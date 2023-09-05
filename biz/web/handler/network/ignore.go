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

package network

import (
	"agent/biz/model/dto/network"
	networkservice "agent/biz/service/network"
	"agent/config"
	"net/http"

	"agent/utils/logger"
	"github.com/gin-gonic/gin"
)

// NetworkIgnore godoc
// @Summary ignore this network(SSID) ,delete net link。[for client LAN/call]
// @Description
// @ID NetworkIgnore
// @Tags network
// @Produce  json
// @Param   networkIgnoreReq body network.NetworkIgnoreReq true  "params"
// @Success 200 {object} dto.BaseRspStr "code=AG-200 success."
// @Router /agent/v1/api/network/ignore [POST]
func NetworkIgnore(c *gin.Context) {
	logger.AppLogger().Debugf("NetworkIgnore POST:%+v", c.Request)

	var reqObject network.NetworkIgnoreReq
	svc := new(networkservice.NetworkIgnoreService)
	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, &reqObject))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
	}
}
