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
 * @Date: 2021-12-09 10:25:41
 * @LastEditors: jeffery
 * @LastEditTime: 2022-04-13 17:13:07
 * @Description:
 */
package docker

import (
	"agent/biz/model/device"
	"agent/config"
	"agent/utils/version"
	"fmt"
	"path"

	"agent/biz/model/device_ability"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

func processTmpEnv(envs, tmpEnv map[string]map[string]string) map[string]map[string]string {
	if len(tmpEnv) > 0 {
		for tmpFile, tmpFileEnv := range tmpEnv {
			if fileEnv, ok := envs[tmpFile]; ok && len(fileEnv) > 0 {
				for k, v := range tmpFileEnv {
					fileEnv[k] = v
				}
				envs[tmpFile] = fileEnv
			}
		}
	}
	return envs
}

func ProcessEnv(composeFile string, tmpEnv map[string]map[string]string) error {
	logger.AppLogger().Debugf("ProcessEnv")
	envs := getEnvMap()
	envs = processTmpEnv(envs, tmpEnv)

	for k, v := range envs {
		content := ""
		for key, value := range v {
			content += fmt.Sprintf("%v=%v\n", key, value)
		}
		p := fileutil.GetFilePath(composeFile)
		f := path.Join(p, k)
		err := fileutil.WriteToFile(f, []byte(content), true)
		if err != nil {
			return err
		}
	}

	return nil
}

func getEnvMap() map[string]map[string]string {
	ret := make(map[string]map[string]string)

	putEnvIntoGateway(ret)
	putEnvIntoClient(ret)
	putEnvIntoNginx(ret)

	logger.AppLogger().Debugf("GetEnvMap, ret:%+v", ret)
	return ret
}

func putEnvIntoGateway(ret map[string]map[string]string) {
	envs := map[string]string{}
	envs["APP_BOX_UUID"] = device.GetDeviceInfo().BoxUuid

	envs["APP_BOX_BTID"] = device.GetDeviceInfo().Btid

	envs["APP_SSPLATFORM_URL"] = device.GetApiBaseUrl()
	envs["APP_PSPLATFORM_URL"] = config.Config.PSPlatform.APIBase.Url

	envs["APP_APPSTORE_APPSIGN_URL"] = config.Config.AppStore.Sign.Url
	envs["APP_APPSTORE_APPAPI_URL"] = config.Config.AppStore.Api.Url

	boxVersion, err := version.GetInstalledAgentVersion()
	if err != nil {
		logger.AppLogger().Errorf("putEnvIntoGateway : %s", err)
	}
	if len(boxVersion) < 1 {
		boxVersion = config.VersionNumber
	}
	envs["APP_BOX_VERSION"] = boxVersion

	envs["APP_SECURITY_API_URL"] = config.Config.Box.SecurityChipAgentHttpAddr

	fingerPrint, err := device.GetDevicePubKeyFingerprint()
	if err != nil {
		logger.AppLogger().Errorf("putEnvIntoGateway GetBoxPubKeyFingerprint failed, err:%v", err)
	} else {
		logger.AppLogger().Debugf("putEnvIntoGateway fingerPrint: [%s]", fingerPrint)
	}
	envs["APP_ACCOUNT_SYSTEM_AGENT_URL_DEVICE_INFO"] = config.Config.EnvDefaultVal.SYSTEM_AGENT_URL_DEVICE_INFO
	envs["APP_SYSTEM_AGENT_URL_BASE"] = config.Config.EnvDefaultVal.SYSTEM_AGENT_URL_BASE
	envs["APP_BOX_KEYFINGERPRINT"] = fingerPrint

	envs["APP_BOX_SUPPORT_SECURITY_CHIP"] = fmt.Sprintf("%v", device_ability.GetAbilityModel().SecurityChipSupport) // 是否支持加密芯片, APP_BOX_SUPPORT_SECURITY_CHIP="true"
	envs["APP_BOX_SUPPORT_EXTERNAL_DISK"] = fmt.Sprintf("%v", device_ability.GetAbilityModel().InnerDiskSupport)    // 是否支持外接磁盘

	envs["APP_BOX_DEPLOY_METHOD"] = "PhysicalBox"
	if device_ability.GetAbilityModel().RunInDocker {
		envs["APP_BOX_DEPLOY_METHOD"] = "DockerBox"
	}
	envs["APP_BOX_DEVICE_MODEL_NUMBER"] = fmt.Sprintf("%v", device_ability.GetAbilityModel().DeviceModelNumber)
	envs["APP_VERSION_TYPE"] = "open_source"
	ret["aospace-gateway.env"] = envs
}

func putEnvIntoClient(ret map[string]map[string]string) {
	envs := map[string]string{}

	// envs["NETWORK_REMOTEAPI"] = config.Config.Platform.APIBase.Url + config.Config.Platform.NetworkRemoteApi.Path
	// logger.AppLogger().Debugf("V1, NETWORK_REMOTEAPI=%v", envs["NETWORK_REMOTEAPI"])

	envs["NETWORK_REMOTEAPI"] = device.GetApiBaseUrl() + config.Config.Platform.NetworkRemoteApi.Path
	logger.AppLogger().Debugf("V2, NETWORK_REMOTEAPI=%v", envs["NETWORK_REMOTEAPI"])

	logger.AppLogger().Debugf("boxInfo=%+v", device.GetDeviceInfo())
	logger.AppLogger().Debugf("boxInfo.NetworkClient=%+v", device.GetDeviceInfo().NetworkClient)
	if device.GetDeviceInfo().NetworkClient != nil {
		logger.AppLogger().Debugf("boxInfo.NetworkClient.SecretKey=%+v", device.GetDeviceInfo().NetworkClient.SecretKey)
		envs["NETWORK_SECRET"] = device.GetDeviceInfo().NetworkClient.SecretKey
		envs["SPACE_NAME_DOMAIN"] = device.GetDeviceInfo().NetworkClient.ClientID
	}

	ret["aonetwork-client.env"] = envs
}

func putEnvIntoNginx(ret map[string]map[string]string) {
	envs := map[string]string{}

	envs["CONFIG_WEBURL"] = config.Config.Platform.WebBase.Url
	ret["aospace-nginx.env"] = envs
}
