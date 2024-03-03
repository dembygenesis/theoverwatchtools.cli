package docker_env

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utils/slice_util"
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
	HostPort     int // Specifying the host port in port bindings.
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
	m.Lock()
	defer m.Unlock()

	allContainers, err := dm.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", err
	}

	targetHostPort := strconv.Itoa(dm.cfg.HostPort)

	for _, ctn := range allContainers {
		hasContainerNameCollision := slice_util.Compare(ctn.Names, func(name string) bool {
			return name == "/"+dm.cfg.Name
		})

		if hasContainerNameCollision {
			if !recreate {
				dm.ContainerID = ctn.ID
				return ctn.ID, nil
			}

			if err := dm.client.ContainerRemove(ctx, ctn.ID, container.RemoveOptions{Force: true}); err != nil {
				return "", fmt.Errorf("failed to remove existing ctn (ID: %s) due to name collision: %w", ctn.ID, err)
			}
		}

		hasPortBindingCollision := slice_util.Compare(ctn.Ports, func(port types.Port) bool {
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

	if err := dm.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	if err := dm.waitForPort(ctx, "0.0.0.0", dm.cfg.HostPort, 5*time.Minute); err != nil {
		return "", fmt.Errorf("error waiting for port %d to become active: %w", dm.cfg.HostPort, err)
	}

	return resp.ID, nil
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
