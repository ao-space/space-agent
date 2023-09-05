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
 * @Author: wenchao
 * @Date: 2021-12-25 11:36:22
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-24 11:24:24
 * @Description:
 */

package docker

import (
	"agent/config"
	"time"

	"agent/utils/logger"
)

func waitEngineReady() {
	logger.AppLogger().Debugf("waitEngineReady")

	setDockerStatus(ContainersWaitOSReady)
	for {
		// 如果 Docker Engine 还没有启动，则等待
		info, err := docker.GetEngineInfo()
		if err == nil {
			logger.AppLogger().Debugf("get docker engine info:%v ", info)
			return
		}
		logger.AppLogger().Warnf("Will Sleep after GetEngineInfo return err:%v ", err)
		// fmt.Printf("Will Sleep after GetEngineInfo return err:%v \n", err)
		time.Sleep(time.Duration(config.Config.Docker.DockerEngineReadyWaitingCheckInterval) * time.Second)
	}
}

func createDockerNetwork() {
	logger.AppLogger().Debugf("createDockerNetwork")

	// err := docker.RemoveNetwork(config.Config.Docker.NetworkName)
	// if err != nil {
	// 	logger.AppLogger().Warnf("failed docker.RemoveNetwork, err:%v ", err)
	// }

	err := docker.CreateNetwork(config.Config.Docker.NetworkName)
	if err != nil {
		logger.AppLogger().Warnf("failed docker.CreateNetwork, err:%v ", err)
		PublishDockerNetworkStatus() // TODO: 这里应该判断是否已经存在了. 但是使用 sdk api 方式才更好判断.
		// return
	} else {
		logger.AppLogger().Debugf("SUCC docker.CreateNetwork, Docker.NetworkName:%v ", config.Config.Docker.NetworkName)
		PublishDockerNetworkStatus()
	}
}
