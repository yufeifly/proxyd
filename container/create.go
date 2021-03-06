package container

import (
	"encoding/json"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/yufeifly/migrator/api/types"
)

// CreateContainer create a container
func CreateContainer(opts types.CreateOpts) (container.ContainerCreateCreatedBody, error) {
	header := "container.CreateContainer"

	config := &container.Config{
		Image: opts.ImageName,
	}
	// unmarshal cmd
	if opts.Cmd != "" {
		var cmd []string
		err := json.Unmarshal([]byte(opts.Cmd), &cmd)
		if err != nil {
			logrus.Errorf("%s, unmarshal cmd err: %v", header, err)
			return container.ContainerCreateCreatedBody{}, err
		}
		config.Cmd = cmd
	}

	openPort, _ := nat.NewPort("tcp", opts.ContainerPort)
	if opts.ExposedPorts != "" {
		exposedPorts := nat.PortSet{}
		err := json.Unmarshal([]byte(opts.ExposedPorts), &exposedPorts)
		if err != nil {
			logrus.Errorf("%s, unmarshal ExposedPorts err: %v", header, err)
			return container.ContainerCreateCreatedBody{}, err
		}
		config.ExposedPorts = exposedPorts
	} else if opts.ContainerPort != "" {

		config.ExposedPorts = nat.PortSet{
			openPort: struct{}{}, //docker容器对外开放的端口
		}
	}

	hostConfig := &container.HostConfig{}
	if opts.PortBindings != "" {
		portBindings := nat.PortMap{}
		err := json.Unmarshal([]byte(opts.PortBindings), &portBindings)
		if err != nil {
			logrus.Errorf("%s, unmarshal PortBindings err: %v", header, err)
			return container.ContainerCreateCreatedBody{}, err
		}
		hostConfig.PortBindings = portBindings
	} else if opts.HostPort != "" {
		hostConfig.PortBindings = nat.PortMap{
			openPort: []nat.PortBinding{nat.PortBinding{
				HostIP:   "0.0.0.0",     //docker容器映射的宿主机的ip
				HostPort: opts.HostPort, //docker 容器映射到宿主机的端口
			}},
		}
	}

	body, err := dockerCli.ContainerCreate(ctx, config, hostConfig, nil, opts.ContainerName)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	logrus.WithFields(logrus.Fields{
		"ContainerID": body.ID,
	}).Info("container created")

	return body, err
}
