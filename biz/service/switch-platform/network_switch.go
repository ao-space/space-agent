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
 * @Author: zhongguang xuyang
 * @Date: 2022-11-15 16:06:11
 * @Last Modified by: xuyang
 * @Last Modified time: 2023-08-14 20:24:31
 */

package switchplatform

import (
	"agent/biz/docker"
	"agent/biz/model/device"
	"agent/biz/model/gt"
	"agent/biz/service/network"
	"agent/config"
	"agent/utils"

	"agent/utils/logger"
)

func networkSwitch(newSSP bool) error {
	// 修改 network client 环境变量
	envFiles := make(map[string]map[string]string)

	if newSSP {
		clientEnvFile := make(map[string]string)

		clientEnvFile["NETWORK_REMOTEAPI"], _ = utils.JoinUrl(si.NewApiBaseUrl, config.Config.Platform.NetworkRemoteApi.Path)
		clientEnvFile["SPACE_NAME_DOMAIN"] = si.ImigrateResult.NetworkClient.ClientID
		clientEnvFile["NETWORK_SECRET"] = si.ImigrateResult.NetworkClient.SecretKey
		envFiles["aonetwork-client.env"] = clientEnvFile

		gwEnvFile := make(map[string]string)
		gwEnvFile["APP_SSPLATFORM_URL"] = si.NewApiBaseUrl

		envFiles["aospace-gateway.env"] = gwEnvFile
	}

	if err := docker.ProcessEnv(config.Config.Docker.ComposeFile, envFiles); err != nil {
		return err
	}

	logger.AppLogger().Debugf("networkClientTurn, transId=%v, envFiles:%+v", si.TransId, envFiles)

	if err := docker.DockerUpImmediately(envFiles); err != nil {
		return err
	}

	logger.AppLogger().Debugf("networkClientTurn, transId=%v, docker.DockerUpImmediately() end ",
		si.TransId)

	return nil
}

func networkSwitchV2(newSSP bool) error {
	envFiles := make(map[string]map[string]string)
	if newSSP {
		// env 文件
		clientEnvFile := make(map[string]string)
		clientEnvFile["NETWORK_REMOTEAPI"], _ = utils.JoinUrl(si.NewApiBaseUrl, config.Config.Platform.NetworkRemoteApi.Path)
		clientEnvFile["SPACE_NAME_DOMAIN"] = si.ImigrateResult.NetworkClient.ClientID
		clientEnvFile["NETWORK_SECRET"] = si.ImigrateResult.NetworkClient.SecretKey
		envFiles["aonetwork-client.env"] = clientEnvFile

		// gt config aonetwork-client.yml
		gtConfig := new(gt.Config)
		remoteAPI, _ := utils.JoinUrl(si.NewApiBaseUrl, config.Config.Platform.NetworkRemoteApi.Path)
		secret := si.ImigrateResult.NetworkClient.SecretKey
		clientId := si.ImigrateResult.NetworkClient.ClientID

		gwEnvFile := make(map[string]string)
		gwEnvFile["APP_SSPLATFORM_URL"] = si.NewApiBaseUrl
		envFiles["aospace-gateway.env"] = gwEnvFile
		err := gtConfig.Switch(remoteAPI, clientId, secret)
		if err != nil {
			return err
		}
	} else {
		gtConfig := new(gt.Config)
		remoteAPI, _ := utils.JoinUrl(device.GetApiBaseUrl(), config.Config.Platform.NetworkRemoteApi.Path)
		secret := device.GetDeviceInfo().NetworkClient.SecretKey
		clientId := device.GetDeviceInfo().NetworkClient.ClientID
		err := gtConfig.Switch(remoteAPI, clientId, secret)
		if err != nil {
			logger.AppLogger().Errorf("rollback gt client config err:%v", err)
			return err
		}
	}
	err := network.RestartGTClient()
	if err != nil {
		return err
	}
	logger.AppLogger().Debugf("networkClientTurn, transId=%v, restart gt client end ", si.TransId)
	logger.AppLogger().Debugf("networkClientTurn, transId=%v, envFiles:%+v", si.TransId, envFiles)
	if err := docker.ProcessEnv(config.Config.Docker.ComposeFile, envFiles); err != nil {
		return err
	}
	if err := docker.DockerUpImmediately(envFiles); err != nil {
		return err
	}
	logger.AppLogger().Debugf("networkClientTurn, transId=%v, docker.DockerUpImmediately() end ",
		si.TransId)

	return nil
}
