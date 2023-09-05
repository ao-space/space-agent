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

package gt

import (
	"agent/biz/model/device"
	"agent/config"
	"agent/utils/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

// Config gt client yaml config
type Config struct {
	Version  string    `yaml:"version" default:"1.0"`
	Services []Service `yaml:"services"`
	Options  Option    `yaml:"options"`
}

// Service gt client service node
type Service struct {
	Local      string `yaml:"local" default:"127.0.0.1:80"`
	HostPrefix string `yaml:"hostPrefix"`
}

// Option gt client Options node
type Option struct {
	ID                string `yaml:"id"`
	Secret            string `yaml:"secret"`
	RemoteConnections int    `yaml:"remoteConnections" default:"5"`
	LogLevel          string `yaml:"logLevel" default:"info"`
	ReconnectDelay    string `yaml:"reconnectDelay" default:"15s"`
	LocalTimeout      string `yaml:"localTimeout" default:"15s"`
	RemoteTimeout     string `yaml:"remoteTimeout" default:"70s"`
	WebrtcMaxPort     int    `yaml:"webrtcMaxPort" default:"62000"`
	WebrtcMinPort     int    `yaml:"webrtcMinPort" default:"61001"`
	RemoteAPI         string `yaml:"remoteAPI"`
}

// Init gt client yaml config init
func (conf *Config) Init() error {
	if _, err := os.Stat(config.Config.GTClient.ConfigPath); err == nil {
		return nil
	}
	var service Service
	conf.Options.ID = device.GetDeviceInfo().NetworkClient.ClientID
	conf.Options.Secret = device.GetDeviceInfo().NetworkClient.SecretKey
	conf.Options.RemoteAPI = device.GetApiBaseUrl() + config.Config.Platform.NetworkRemoteApi.Path
	conf.Options.RemoteConnections = 5
	conf.Options.ReconnectDelay = "15s"
	conf.Options.LocalTimeout = "15s"
	conf.Options.RemoteTimeout = "70s"
	conf.Options.WebrtcMinPort = 61001
	conf.Options.WebrtcMaxPort = 62000
	conf.Options.LogLevel = "info"
	strconv.Itoa(int(config.Config.GateWay.LanPort))
	if config.RunningOnLinux() {
		service = Service{
			Local:      "http://127.0.0.1:" + strconv.Itoa(int(config.Config.GateWay.LanPort)),
			HostPrefix: conf.Options.ID,
		}
	} else {
		service = Service{
			Local:      "http://" + config.Config.Docker.NginxContainerName + ":80",
			HostPrefix: conf.Options.ID,
		}
	}

	//
	//if hardware.RunningInDocker() {
	//	if config.RunningOnLinux() {
	//		service = Service{
	//			Local:      "http://127.0.0.1:9980",
	//			HostPrefix: conf.Options.ID,
	//		}
	//	} else {
	//		service = Service{
	//			Local:      "http://" + config.Config.Docker.NginxContainerName + ":9980",
	//			HostPrefix: conf.Options.ID,
	//		}
	//	}
	//} else {
	//	service = Service{
	//		Local:      "http://127.0.0.1:80",
	//		HostPrefix: conf.Options.ID,
	//	}
	//}
	conf.Services = append(conf.Services, service)
	conf.Version = "1.0"
	// 将配置写入文件
	err := conf.Save()
	if err != nil {
		return err
	}
	return nil
}

// AddNewService when install aospace apps in privacy mode,
// add service node to gt client yaml config
func (conf *Config) AddNewService(appName string, port string, token string) error {
	// 修改 Network Client 的转发配置
	currentConf, err := Load()
	if err != nil {
		return err
	}
	// 新增Service节点
	newService := Service{
		Local:      "http://127.0.0.1:443", // 第三方应用的server and port
		HostPrefix: appName + "-" + token,  // 第三方应用的子域名
	}
	currentConf.Services = append(currentConf.Services, newService)
	err = currentConf.Save()
	if err != nil {
		return err
	}
	logger.AppLogger().Infof("modify network client config successfully.")
	return nil
}

// RemoveService when uninstall aospace apps in privacy mode,
// remove service node in gt client yaml config
func (conf *Config) RemoveService(appName string, port string, appToken string) error {
	// 删除 Network Client 的转发配置
	currentConf, err := Load()
	if err != nil {
		return err
	}
	// 删除Service节点
	prefixToRemove := appName + "-" + appToken // 指定要删除的 hostPrefix 值
	updatedServices := make([]Service, 0)
	for _, service := range currentConf.Services {
		if !strings.Contains(service.HostPrefix, prefixToRemove) {
			updatedServices = append(updatedServices, service)
		}
	}
	err = currentConf.Save()
	if err != nil {
		return err
	}
	logger.AppLogger().Infof("remove gt-client service successfully.")
	return nil
}

// Switch when switch platform,need switch network environment
func (conf *Config) Switch(remoteAPI, clientId, secret string) error {
	// 修改 Network Client 的转发配置
	currentConf, err := Load()
	if err != nil {
		return err
	}
	currentConf.Options.ID = clientId
	currentConf.Options.Secret = secret
	currentConf.Options.RemoteAPI = remoteAPI
	return currentConf.Save()
}

// Save save gt client yaml config
func (conf *Config) Save() error {
	newData, err := yaml.Marshal(*conf)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path.Dir(config.Config.GTClient.ConfigPath)); err != nil {
		err = os.Mkdir(path.Dir(config.Config.GTClient.ConfigPath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(config.Config.GTClient.ConfigPath, newData, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Load load current gt client yaml config
func Load() (*Config, error) {
	var conf Config
	data, err := ioutil.ReadFile(config.Config.GTClient.ConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
