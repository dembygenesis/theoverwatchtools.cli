package dockerenv

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/sliceutil"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
	m sync.Mutex
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
		return nil, fmt.Errorf("new client: %v", err)
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
	m.Lock()
	defer m.Unlock()

	allContainers, err := dm.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", fmt.Errorf("container list: %v", err)
	}

	targetHostPort := strconv.Itoa(dm.cfg.HostPort)

	for _, ctn := range allContainers {
		// Check if there is a port, and mysql image that is the same
		hasPortBindingCollision := sliceutil.Compare(ctn.Ports, func(port types.Port) bool {
			return strconv.Itoa(int(port.PublicPort)) == targetHostPort && port.Type == "tcp"
		})
		if hasPortBindingCollision {
			if !recreate {
				dm.ContainerID = ctn.ID
				return ctn.ID, nil
			}

			if err := dm.client.ContainerRemove(ctx, ctn.ID, container.RemoveOptions{Force: true}); err != nil {
				return "", fmt.Errorf("failed to remove existing ctn (ID: %s) due to port collision: %w", ctn.ID, err)
			}
		}

		hasContainerNameCollision := sliceutil.Compare(ctn.Names, func(name string) bool {
			return name == "/"+dm.cfg.Name
		})

		hasExistingStoppedContainer := ctn.State == "exited" || ctn.State == "created"
		if hasExistingStoppedContainer && hasContainerNameCollision {
			log.Printf("Existing stopped container found with name %s. Attempting to start it.", dm.cfg.Name)
			if err := dm.client.ContainerStart(ctx, ctn.ID, container.StartOptions{}); err != nil {
				return "", fmt.Errorf("failed to start existing container: %s, with error: %w", ctn.ID, err)
			}
			dm.ContainerID = ctn.ID
			return ctn.ID, nil
		}

		if hasContainerNameCollision {
			if !recreate {
				dm.ContainerID = ctn.ID
				return ctn.ID, nil
			}

			if err := dm.client.ContainerRemove(ctx, ctn.ID, container.RemoveOptions{Force: true}); err != nil {
				return "", fmt.Errorf("failed to remove existing ctn (ID: %s) due to name collision: %w", ctn.ID, err)
			}
		}
	}

	return dm.createContainer(ctx)
}

func (dm *DockerEnv) createContainer(ctx context.Context) (string, error) {
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

	_, _, err := dm.client.ImageInspectWithRaw(ctx, dm.cfg.Image)
	if err != nil {
		if client.IsErrNotFound(err) {
			rc, err := dm.client.ImagePull(ctx, dm.cfg.Image, types.ImagePullOptions{})
			if err != nil {
				return "", fmt.Errorf("pulling image: %v", err)
			}
			defer rc.Close()

			if _, err = new(bytes.Buffer).ReadFrom(rc); err != nil {
				return "", fmt.Errorf("reading image pull response: %v", err)
			}
		} else {
			return "", fmt.Errorf("error checking for image: %v", err)
		}
	}

	resp, err := dm.client.ContainerCreate(ctx, contConfig, hostConfig, nil, nil, dm.cfg.Name)
	if err != nil {
		return "ok", fmt.Errorf("create container: %v", err)
	}

	dm.ContainerID = resp.ID

	if err = dm.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	if err = dm.waitForPort(ctx, "0.0.0.0", dm.cfg.HostPort, 5*time.Minute); err != nil {
		return "", fmt.Errorf("error waiting for port %d to become active: %w", dm.cfg.HostPort, err)
	}

	if dm.cfg.WaitFor != "" {
		timeoutDuration := 60 * time.Second
		if err = dm.waitForLogMessage(ctx, dm.ContainerID, dm.cfg.WaitFor, timeoutDuration); err != nil {
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
