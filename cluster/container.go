package cluster

import (
	"github.com/samalba/dockerclient"
)

// Container is exported
type Container struct {
	dockerclient.Container

	Info   dockerclient.ContainerInfo
	Engine *Engine
}

// Start the container.
func (c *Container) Start() error {
	return c.Engine.client.StartContainer(c.Id, nil)
}

// Stop the container.
func (c *Container) Stop() error {
	return c.Engine.client.StopContainer(c.Id, 8)
}

// Restart the container.
func (c *Container) Restart(timeout int) error {
	return c.Engine.client.RestartContainer(c.Id, timeout)
}

// Kill the container.
func (c *Container) Kill(signal string) error {
	return c.Engine.client.KillContainer(c.Id, signal)
}

// Pause the container.
func (c *Container) Pause() error {
	return c.Engine.client.PauseContainer(c.Id)
}

// Unpause the container.
func (c *Container) Unpause() error {
	return c.Engine.client.UnpauseContainer(c.Id)
}
