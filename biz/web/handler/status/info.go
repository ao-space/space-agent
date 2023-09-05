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
 * @Date: 2021-10-16 10:23:56
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-13 13:34:02
 * @Description:
 */
package status

import (
	"agent/biz/docker"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/dto/status"
	"agent/config"
	"agent/utils/deviceid"
	"fmt"
	"net/http"

	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/encrypt/random"
	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/gin-gonic/gin"
)

// Info godoc
// @Summary get server info
// @Description get server info
// @ID info
// @Tags Info
// @Accept  plain
// @Produce  json
// @Success 200 {object}  dto.BaseRspStr{results=status.Info} "code=AG-200 success"
// @Router /agent/info [GET]
func Info(c *gin.Context) {

	var result *status.Info
	if config.Config.DebugMode {

		result = &status.Info{Status: "OK",
			Version:    config.Version,
			TheBoxInfo: device.GetDeviceInfo(),

			// IsClientPaired: device.IsClientPaired(),
			// IsBoxInit:      device.IsBoxInit(),
			DockerStatus: docker.GetDockerStatus(),

			// TheClientInfo:     device.GetClientInfo(),
			TheBoxPriKeyBytes: string(device.GetDevicePriKey()),
			TheBoxPublicKey:   string(device.GetDevicePubKey()),
		}
	} else {
		result = &status.Info{Status: "OK",
			Version: config.Version,
			// IsClientPaired: device.IsClientPaired(),
			// IsBoxInit:      device.IsBoxInit(),
			DockerStatus: docker.GetDockerStatus(),
		}
	}

	abilityModel := device_ability.GetAbilityModel()
	if abilityModel.RunInDocker {
		err := fileutil.WriteToFile(config.Config.Box.HostIpFile, []byte(c.Request.Host), true)
		if err != nil {
			err1 := fmt.Errorf("failed write HostIpFile, %+v", err)
			logger.AppLogger().Debugf("info POST, %+v", err1)
			c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
				RequestId: random.GenUUID(),
				Message:   err1.Error()})
			return
		}
		result.QrCode = device.GetQrCode()

		snNumber, err := deviceid.GetSnNumber(config.Config.Box.SnNumberStoreFile)
		result.TryoutCodeVerified = false
		if err == nil && len(snNumber) > 0 {
			result.TryoutCodeVerified = true
		}
	}

	c.IndentedJSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		Message:   "OK",
		RequestId: random.GenUUID(),
		Results:   result})
}
