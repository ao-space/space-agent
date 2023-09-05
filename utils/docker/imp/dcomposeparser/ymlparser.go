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

package dcomposeparser

import (
	"github.com/spf13/viper"
)

type ServiceStruct struct {
	ContainerName string            `mapstructure:"container_name"`
	Image         string            `mapstructure:"image"`
	Restart       string            `mapstructure:"restart"`
	Ports         []string          `mapstructure:"ports"`
	DependsOn     []string          `mapstructure:"depends_on"`
	Environment   map[string]string `mapstructure:"environment"`
	Volumes       []string          `mapstructure:"volumes"`
}

type DockerComposeStruct struct {
	Version string `mapstructure:"version"`

	Services map[string]*ServiceStruct `mapstructure:"services"`

	Networks struct {
		Default struct {
			External struct {
				Name string `mapstructure:"name"`
			} `mapstructure:"external"`
		} `mapstructure:"default"`
	} `mapstructure:"networks"`
}

func ParseYml(filePath string) (*DockerComposeStruct, error) {
	// 设置配置文件信息
	viper.SetConfigType("yml")
	viper.SetConfigFile(filePath)

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// 将文件内容解析后封装到cfg对象中
	var c DockerComposeStruct
	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
