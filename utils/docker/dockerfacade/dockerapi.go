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

/*
 * @Author: jeffery
 * @Date: 2022-02-21 10:01:15
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-28 15:30:23
 * @Description:
 */
package dockerfacade

import (
	"agent/utils/docker/dockermodel"
	"agent/utils/docker/imp/dcomposeapi"
	"agent/utils/docker/imp/dengineapi"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/spf13/viper"

	"agent/utils/logger"
)

type DockerFacade struct {
}

func NewDockerFacade() *DockerFacade {
	return &DockerFacade{}
}

func (dock *DockerFacade) SetClientVersion(v string) {
	dengineapi.SetClientVersion(v)
}

func (dock *DockerFacade) GetEngineInfo() (string, error) {
	return dengineapi.Info()
}

func (dock *DockerFacade) ChansReader() (chan string, chan string) {
	stdOutput := make(chan string, 128)
	errOutput := make(chan string, 128)
	go func(stdOutput chan string) {
		for line := range stdOutput {
			logger.AppLogger().Debugf("##[" + line + "]")
		}
	}(stdOutput)
	go func(errOutput chan string) {
		for line := range errOutput {
			logger.AppLogger().Debugf("##[" + line + "]")
		}
	}(errOutput)

	return stdOutput, errOutput
}

func (dock *DockerFacade) ListImages() ([]*dockermodel.DockerImage, error) {
	return dengineapi.ListImages(nil)
}

func (dock *DockerFacade) ListContainers() ([]*dockermodel.DockerContainer, error) {
	return dengineapi.ListContainers(nil)
}

func (dock *DockerFacade) RemoveImage(imageId string) error {
	return dengineapi.RemoveImage(nil, imageId, types.ImageRemoveOptions{})
}

func (dock *DockerFacade) RemoveContainer(containerId string) error {
	return dengineapi.RemoveContainer(nil, containerId, types.ContainerRemoveOptions{})
}

func (dock *DockerFacade) Exec(containerId string, cmd []string) error {
	return dengineapi.Exec(nil, containerId, cmd)
}

func (dock *DockerFacade) FindContainer(containerName string) (string, error) {
	return dengineapi.FindContainer(containerName)
}

func (dock *DockerFacade) RestartContainer(containerName string) error {
	return dengineapi.RestartContainer(containerName)
}

func (dock *DockerFacade) CreateNetwork(networkName string) error {
	stdOutput, errOutput := dock.ChansReader()
	return dcomposeapi.CreateNetwork(networkName, stdOutput, errOutput)
}

func (dock *DockerFacade) RemoveNetwork(networkName string) error {
	stdOutput, errOutput := dock.ChansReader()
	return dcomposeapi.RemoveNetwork(networkName, stdOutput, errOutput)
}

func (dock *DockerFacade) Start(composeFile string) error {
	stdOutput, errOutput := dock.ChansReader()
	return dcomposeapi.StartContainers(composeFile, stdOutput, errOutput)
}

func (dock *DockerFacade) Pull(composeFile string) error {
	stdOutput, errOutput := dock.ChansReader()
	return dcomposeapi.PullImages(composeFile, stdOutput, errOutput)
}

func (dock *DockerFacade) Create(composeFile string) error {
	stdOutput, errOutput := dock.ChansReader()
	return dcomposeapi.CreateContainers(composeFile, stdOutput, errOutput)
}

func (dock *DockerFacade) DownContainers(composeFile string) error {
	stdOutput, errOutput := dock.ChansReader()
	return dcomposeapi.DownContainers(composeFile, stdOutput, errOutput)
}

func (dock *DockerFacade) StopSpecifiedContainers(composeFile string, containers ...string) error {
	stdOutput, errOutput := dock.ChansReader()
	return dcomposeapi.StopSpecifiedContainers(composeFile, stdOutput, errOutput, containers...)
}

func (dock *DockerFacade) UpContainers(composeFile string, excludeServices []string) (chan string, chan string, error) {
	if excludeServices == nil || len(excludeServices) < 1 {
		stdOutput, errOutput := dock.ChansReader()
		err := dcomposeapi.UpContainers(composeFile, stdOutput, errOutput)
		return stdOutput, errOutput, err
	}

	services, err := getComposeFileServiceNames(composeFile)
	if err != nil {
		return nil, nil, err
	}
	includeServices := arrayComplement(services, excludeServices)
	if len(services) < 1 {
		return nil, nil, fmt.Errorf("excludeServices:%v, includeServices:%v", excludeServices, includeServices)
	}

	return dock.UpSpecifiedContainers(composeFile, includeServices...)
}

func (dock *DockerFacade) UpContainersWithNoRecreate(composeFile string, excludeServices []string) (chan string, chan string, error) {
	if excludeServices == nil || len(excludeServices) < 1 {
		stdOutput, errOutput := dock.ChansReader()
		err := dcomposeapi.UpContainersWithNoRecreate(composeFile, stdOutput, errOutput)
		return stdOutput, errOutput, err
	}

	services, err := getComposeFileServiceNames(composeFile)
	if err != nil {
		return nil, nil, err
	}
	includeServices := arrayComplement(services, excludeServices)
	if len(services) < 1 {
		return nil, nil, fmt.Errorf("excludeServices:%v, includeServices:%v", excludeServices, includeServices)
	}

	return dock.UpSpecifiedContainers(composeFile, includeServices...)
}

func (dock *DockerFacade) UpSpecifiedContainers(composeFile string, containers ...string) (chan string, chan string, error) {
	stdOutput, errOutput := dock.ChansReader()
	err := dcomposeapi.UpSpecifiedContainers(composeFile, stdOutput, errOutput, containers...)
	return stdOutput, errOutput, err
}

// 全集是a, 子集b. 返回补集.
func arrayComplement(a, b []string) []string {
	setExclude := make(map[string]string, 0)
	for _, v := range b {
		setExclude[v] = ""
	}

	includes := make([]string, 0, len(a))
	for _, cur := range a {
		if _, ok := setExclude[cur]; !ok {
			includes = append(includes, cur)
		}
	}
	return includes
}

func getComposeFileServiceNames(dockerComposeFile string) ([]string, error) {

	dir, filename := filepath.Split(dockerComposeFile)
	viper.SetConfigName(strings.TrimSuffix(filename, path.Ext(filename)))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	v := viper.Get("services")
	m := v.(map[string]interface{})

	services := make([]string, 0)
	for k := range m {
		services = append(services, k)
	}
	return services, nil
}
