// Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dengineapi

import (
	"agent/utils/docker/dockermodel"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func ListContainers(cli *client.Client) ([]*dockermodel.DockerContainer, error) {

	var err error
	if cli == nil {
		cli, err = NewClient()
		if err != nil {
			return nil, fmt.Errorf("failed NewClient, err:%v", err)
		}
		defer cli.Close()
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	ret := []*dockermodel.DockerContainer{}
	// fmt.Printf("len(containers)=%v \n", len(containers))
	for _, container := range containers {
		// fmt.Printf("container: %+v\n", container)

		ret = append(ret, &dockermodel.DockerContainer{
			ImageID: container.ImageID,
			Image:   container.Image,
			ID:      container.ID,
			State:   container.State,
			Names:   container.Names,
			Created: container.Created})
	}

	return ret, nil
}

// 创建容器
type CreateContainerReq struct {
	ImageName     string
	ContainerName string
	Port          []string
	Env           []string
	Volumes       []string
	Network       string
}

func CreateContainer(cli *client.Client, req *CreateContainerReq) (string, error) {

	exposedPorts, portBindings, _ := nat.ParsePortSpecs(req.Port)

	config := &container.Config{Image: req.ImageName,
		ExposedPorts: exposedPorts,
		Env:          req.Env}

	m := []mount.Mount{}
	for _, v := range req.Volumes {
		s := strings.Split(v, ":")
		if len(s) < 2 {
			continue
		}
		m = append(m, mount.Mount{
			Type:     mount.TypeBind,
			ReadOnly: false,
			Source:   s[0],
			Target:   s[1],
		})
	}
	hostConfig := &container.HostConfig{PortBindings: portBindings,
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Mounts:        m,
		// NetworkMode:   container.NetworkMode(req.Network),
	}

	networkConfig := &network.NetworkingConfig{}

	body, err := cli.ContainerCreate(context.Background(), config, hostConfig, networkConfig,
		nil, req.ContainerName)
	if err != nil {
		return "", err
	}
	fmt.Printf("ID: %s\n", body.ID)
	return body.ID, nil
}

func Exec(cli *client.Client, containerId string, cmd []string) error {
	var err error
	if cli == nil {
		cli, err = NewClient()
		if err != nil {
			return fmt.Errorf("failed NewClient, err:%v", err)
		}
		defer cli.Close()
	}
	execId, err := cli.ContainerExecCreate(context.Background(), containerId, types.ExecConfig{
		Cmd: cmd,
	})
	if err != nil {
		return err
	}
	err = cli.ContainerExecStart(context.Background(), execId.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	})
	if err != nil {
		return err
	}
	return nil
}

func FindContainer(containerName string) (string, error) {

	cli, err := NewClient()
	if err != nil {
		return "", fmt.Errorf("failed NewClient, err:%v", err)
	}
	defer cli.Close()

	var args = filters.NewArgs()
	args.Add("name", containerName)
	containerInfos, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: args})

	if err != nil {
		return "", err
	}
	if len(containerInfos) > 0 {
		for i, ci := range containerInfos {
			if strings.TrimLeft(ci.Names[0], "/") == containerName {
				return containerInfos[i].ID, nil
			}
		}
	}
	return "", nil

}

// 删除
func RemoveContainer(cli *client.Client, containerID string, removeOpt types.ContainerRemoveOptions) error {
	var err error
	if cli == nil {
		cli, err = NewClient()
		if err != nil {
			return fmt.Errorf("failed NewClient, err:%v", err)
		}
		defer cli.Close()
	}
	err = cli.ContainerRemove(context.Background(), containerID, removeOpt)
	return err
}

func RestartContainer(containerId string) error {

	var duration = 10 * time.Second
	cli, err := NewClient()
	if err != nil {
		return fmt.Errorf("failed NewClient, err:%v", err)
	}
	defer cli.Close()

	err = cli.ContainerRestart(context.Background(), containerId, &duration)
	if err != nil {
		return err
	}

	return nil
}
