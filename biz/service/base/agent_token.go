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

package base

import (
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/service/encwrapper"
	"agent/utils/jwt"
	"agent/utils/logger"
	"fmt"
	"time"
)

const (
	TokenTypeBind = "BIND_API_TOKEN"
)

func VerifyAgentToken(agentToken string) error {
	if device_ability.GetAbilityModel().SecurityChipSupport {

	} else {

	}
	return nil
}

func CreateAgentToken(clientUuid string) (string, error) {
	logger.AppLogger().Debugf("CreateAgentToken")

	expiredAt := time.Now().Add(24 * time.Hour * 365 * 3)

	if device_ability.GetAbilityModel().SecurityChipSupport {
		logger.AppLogger().Debugf("createAgentToken, using SecurityChipSupport, clientUuid: %v", clientUuid)

		jwtToken, err := jwt.GenerateJWT(device.GetDeviceInfo().BoxUuid, "",
			[]string{clientUuid}, expiredAt, nil,
			map[string]string{"tokenType": TokenTypeBind}, nil)
		if err != nil {
			return "", fmt.Errorf("failed GenerateJWT, err:%v", err)
		}
		return jwtToken, nil

	} else {
		logger.AppLogger().Debugf("createAgentToken, clientUuid: %v", clientUuid)

		pri, err := encwrapper.GetPrivateKey(string(device.GetDevicePriKey()))
		if err != nil {
			return "", err
		}

		jwtToken, err := jwt.GenerateJWT(device.GetDeviceInfo().BoxUuid, "",
			[]string{clientUuid}, expiredAt, nil,
			map[string]string{"tokenType": TokenTypeBind}, pri)
		if err != nil {
			return "", fmt.Errorf("failed GenerateJWT, err:%v", err)
		}
		return jwtToken, nil

	}
}
