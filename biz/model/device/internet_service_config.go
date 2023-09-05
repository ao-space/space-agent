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

package device

import (
	"agent/config"
	"agent/utils/logger"
	"fmt"
	"sync"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

type InternetServiceConfig struct {
	EnableInternetAccess bool `json:"enableInternetAccess"`
}

var c *InternetServiceConfig
var lock sync.Mutex

func init() {
	c = &InternetServiceConfig{EnableInternetAccess: true}
	if !fileutil.IsFileExist(config.Config.Box.InternetServiceConfigFile) {
		fmt.Printf("InternetServiceConfig file not exist, %v \n", config.Config.Box.InternetServiceConfigFile)
		writeToFile()
	} else {
		fmt.Printf("InternetServiceConfig file exist, %v \n", config.Config.Box.InternetServiceConfigFile)
	}
	readFromFile()
}

//func (isc *InternetServiceConfig) Get() *InternetServiceConfig {
//	isc.mu.Lock()
//	defer isc.mu.Unlock()
//	return c
//}
//
//func (isc *InternetServiceConfig) Set() {
//	isc.mu.Lock()
//	defer isc.mu.Unlock()
//	c = isc
//	if c == nil {
//		c = &InternetServiceConfig{EnableInternetAccess: true}
//	}
//	writeToFile()
//}

func GetConfig() *InternetServiceConfig {
	logger.AppLogger().Debugf("InternetServiceConfig GetConfig c:%+v", c)
	lock.Lock()
	defer lock.Unlock()
	return c
}

func SetConfig(config *InternetServiceConfig) {
	logger.AppLogger().Debugf("InternetServiceConfig SetConfig, config:%+v", config)
	lock.Lock()
	defer lock.Unlock()
	c = config
	if c == nil {
		c = &InternetServiceConfig{EnableInternetAccess: true}
	}
	writeToFile()
}

func readFromFile() {
	if fileutil.IsFileExist(config.Config.Box.InternetServiceConfigFile) {
		logger.AppLogger().Debugf("InternetServiceConfig file exist:%v",
			config.Config.Box.InternetServiceConfigFile)
		fileutil.ReadFileJsonToObject(config.Config.Box.InternetServiceConfigFile, c)
	} else {
		logger.AppLogger().Debugf("InternetServiceConfig file not exist:%v",
			config.Config.Box.InternetServiceConfigFile)
	}
}

func writeToFile() {
	err := fileutil.WriteToFileAsJson(config.Config.Box.InternetServiceConfigFile, c, "  ", true)
	logger.AppLogger().Debugf("InternetServiceConfig writeToFile:%v, err:%+v",
		config.Config.Box.InternetServiceConfigFile, err)
}
