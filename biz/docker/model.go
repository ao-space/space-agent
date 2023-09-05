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

package docker

import (
	"agent/config"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

type AOComposeFile struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks Network            `yaml:"networks"`
}

type Service struct {
	ContainerName string                   `yaml:"container_name"`
	Image         string                   `yaml:"image"`
	Restart       string                   `yaml:"restart"`
	Healthcheck   HealthcheckConfig        `yaml:"healthcheck,omitempty"`
	Ports         []string                 `yaml:"ports,omitempty"`
	DependOn      map[string]DependsConfig `yaml:"depend_on,omitempty"`
	Environment   map[string]string        `yaml:"environment,omitempty"`
	Volumes       []string                 `yaml:"volumes,omitempty"`
}

type HealthcheckConfig struct {
	Test        interface{} `yaml:"test"`
	Interval    string      `yaml:"interval"`
	Timeout     string      `yaml:"timeout"`
	Retries     int         `yaml:"retries"`
	StartPeriod string      `yaml:"start_period"`
}

type DependsConfig struct {
	Condition string `yaml:"condition,omitempty"`
}

type Network struct {
	External map[string]ExternalConfig `yaml:"external"`
}

type ExternalConfig struct {
	Name string `yaml:"name"`
}

func (a *AOComposeFile) FixVolume(homeDir string) error {
	for _, service := range a.Services {
		for _, volume := range service.Volumes {
			volume = homeDir + volume
		}
	}
	return a.SaveComposeFile()
}

func (a *AOComposeFile) SaveComposeFile() error {
	return writeYAML(config.Config.Docker.CustomComposeFile, *a)
}

func LoadComposeFile() (AOComposeFile, error) {
	var (
		composeFile AOComposeFile
		err         error
	)
	err = readYAML(config.Config.Docker.CustomComposeFile, &composeFile)
	return composeFile, err
}

func readYAML(filePath string, config interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	return nil
}

func writeYAML(filePath string, yamlConfig interface{}) error {
	newData, err := yaml.Marshal(&yamlConfig)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Dir(filePath), os.ModeDir)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, newData, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
