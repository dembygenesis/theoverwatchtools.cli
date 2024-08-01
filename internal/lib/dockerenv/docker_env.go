package dockerenv

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	m              sync.Mutex
	waitForTimeout = 30 * time.Second
)

type DockerEnv struct {
	client      *client.Client
	cfg         *ContainerConfig
	ContainerID string
}

type ContainerConfig struct {
	Name         string
	Image        string
	Env          []string
	ExposedPort  int
	ExternalPort int
	HostPort     int    // Specifying the host port in port bindings.
	WaitFor      string // Specifies the string you wait for, before confirming
	Cmd          []string
}

func New(cfg *ContainerConfig) (*DockerEnv, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerEnv{
		client: cli,
		cfg:    cfg,
	}, nil
}

func (dm *DockerEnv) Cleanup(ctx context.Context) error {
	if strings.TrimSpace(dm.ContainerID) == "" {
		return fmt.Errorf("container ID is missing")
	}

	if err := dm.client.ContainerRemove(ctx, dm.ContainerID, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("remove: %v", err)
	}

	return nil
}

// UpsertContainer upserts a new container, be careful with this function
// because it will remove other running instances with colliding port, OR names.
func (dm *DockerEnv) UpsertContainer(ctx context.Context, recreate bool) (string, error) {
	fmt.Println("========= 1")
	m.Lock()
	defer m.Unlock()

	allContainers, err := dm.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", err
	}

	fmt.Println("========= 2")

	targetHostPort := strconv.Itoa(dm.cfg.HostPort)

	for _, ctn := range allContainers {
		fmt.Println("========= 2.5")
		hasPortBindingCollision := checkPortBindingCollision(ctn, targetHostPort)
		hasContainerNameCollision := checkNameCollision(ctn, dm.cfg.Name)
		isStopped := ctn.State == "exited" || ctn.State == "created"

		if isStopped && hasContainerNameCollision {
			log.Printf("Existing stopped container found with name %s. Attempting to start it.", dm.cfg.Name)
			if err := dm.client.ContainerStart(ctx, ctn.ID, container.StartOptions{}); err != nil {
				return "", fmt.Errorf("failed to start existing container: %s, with error: %w", ctn.ID, err)
			}

			dm.ContainerID = ctn.ID
			return ctn.ID, nil
		}

		if hasContainerNameCollision || hasPortBindingCollision {
			if !recreate {
				dm.ContainerID = ctn.ID
				return ctn.ID, nil
			}

			if err := dm.client.ContainerRemove(ctx, ctn.ID, container.RemoveOptions{Force: true}); err != nil {
				return "", fmt.Errorf("failed to remove existing ctn (ID: %s): %w", ctn.ID, err)
			}
		}
	}

	fmt.Println("========= 4")
	return dm.createContainer(ctx)
}

func checkPortBindingCollision(ctn types.Container, targetHostPort string) bool {
	for _, port := range ctn.Ports {
		if strconv.Itoa(int(port.PublicPort)) == targetHostPort && port.Type == "tcp" {
			return true
		}
	}
	return false
}

func checkNameCollision(ctn types.Container, name string) bool {
	for _, n := range ctn.Names {
		if n == "/"+name {
			return true
		}
	}
	return false
}

func (dm *DockerEnv) createContainer(ctx context.Context) (string, error) {
	// Check if the image is already present locally
	imagePresent, err := dm.isImagePresent(ctx, dm.cfg.Image)
	if err != nil {
		return "", fmt.Errorf("failed to check if image is present: %w", err)
	}

	// Pull the image only if it is not present
	if !imagePresent {
		out, err := dm.client.ImagePull(ctx, dm.cfg.Image, image.PullOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to pull image: %w", err)
		}
		defer out.Close()
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading image pull response: %w", err)
		}
	}

	contConfig := &container.Config{
		Image: dm.cfg.Image,
		Env:   dm.cfg.Env,
		ExposedPorts: nat.PortSet{
			nat.Port(fmt.Sprintf("%d/tcp", dm.cfg.ExposedPort)): struct{}{},
		},
		Cmd: dm.cfg.Cmd,
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%d/tcp", dm.cfg.ExposedPort)): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(dm.cfg.HostPort),
				},
			},
		},
	}

	resp, err := dm.client.ContainerCreate(ctx, contConfig, hostConfig, nil, nil, dm.cfg.Name)
	if err != nil {
		return "", err
	}

	dm.ContainerID = resp.ID

	// if err = dm.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
	if err = dm.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	if err = dm.waitForPort(ctx, "0.0.0.0", dm.cfg.HostPort, 5*time.Minute); err != nil {
		return "", fmt.Errorf("error waiting for port %d to become active: %w", dm.cfg.HostPort, err)
	}

	if dm.cfg.WaitFor != "" {
		if err = dm.waitForLogMessage(ctx, dm.ContainerID, dm.cfg.WaitFor, waitForTimeout); err != nil {
			return "", fmt.Errorf("err waiting for text: %s, with err: %v", dm.cfg.WaitFor, err)
		}
	}

	return resp.ID, nil
}

// waitForLogMessage waits for a specific log message from the container.
func (dm *DockerEnv) waitForLogMessage(ctx context.Context, containerID, message string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	options := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true}
	logStream, err := dm.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logStream.Close()

	scanner := bufio.NewScanner(logStream)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, message) {
					log.Printf("Detected log message: %s", message)
					return nil
				}
			} else if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading container logs: %w", err)
			}
		}
	}

	return fmt.Errorf("timeout reached waiting for log message: %s", message)
}

func (dm *DockerEnv) waitForPort(ctx context.Context, host string, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	target := fmt.Sprintf("%s:%d", host, port)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", target, time.Second)
		if err == nil {
			conn.Close()
			log.Printf("Port %d is now active.", port)
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			// Retry after a delay
		}
	}
	return fmt.Errorf("timeout reached waiting for port %d to become active", port)
}

// isImagePresent checks if the specified image is present locally
func (dm *DockerEnv) isImagePresent(ctx context.Context, imageName string) (bool, error) {
	imageFilters := filters.NewArgs()
	imageFilters.Add("reference", imageName)

	images, err := dm.client.ImageList(ctx, image.ListOptions{Filters: imageFilters})
	if err != nil {
		return false, err
	}

	return len(images) > 0, nil
}
