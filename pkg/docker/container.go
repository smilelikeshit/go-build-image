package docker

import (
	"context"

	docker "github.com/fsouza/go-dockerclient"
)

type Container interface {
	StopAndRemove() error
	GetIP() (string, error)
}

type dockerContainer struct {
	cli *docker.Client
	ctx context.Context
	id  string
}

func newContainer(cli *docker.Client, ctx context.Context, id string) Container {
	return &dockerContainer{cli: cli, ctx: ctx, id: id}
}

func (c *dockerContainer) StopAndRemove() error {
	err := c.cli.StopContainer(c.id, 100)
	if err != nil {
		return err
	}

	if err := c.cli.RemoveContainer(docker.RemoveContainerOptions{ID: c.id, RemoveVolumes: true, Context: c.ctx}); err != nil {
		return err
	}
	return nil
}

func (c *dockerContainer) GetIP() (string, error) {
	container, err := c.cli.InspectContainerWithContext(c.id, c.ctx)
	if err != nil {
		return "", err
	}

	if container.NetworkSettings != nil {
		return container.NetworkSettings.IPAddress, nil
	}

	return "", nil
}
