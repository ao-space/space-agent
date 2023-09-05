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

// NetworkConfig godoc
// @Summary network config[for client LAN/Call]
// @Description
// @ID POSTNetworkConfig
// @Tags network
// @Accept  json
// @Produce  json
// @Param   networkConfigReq body network.NetworkConfigReq true  "网络配置参数".
// @Success 200 {object} dto.BaseRspStr "code=AG-200 成功."
// @Router /agent/v1/api/network/config [POST]
func PostNetworkConfig(c *gin.Context) {
	logger.AppLogger().Debugf("NetworkConfig POST:%+v", c.Request)

	var reqObject network.NetworkConfigReq
	svc := networkservice.NewPostNetworkConfigService()
	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, &reqObject))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
	}

}

// NetworkConfig godoc
// @Summary get network info。[for client LAN/call]
// @Description
// @ID GETNetworkStatus
// @Tags network
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=network.NetworkStatusRsp} "code=AG-200 success."
// @Router /agent/v1/api/network/config [GET]
func GetNetworkConfig(c *gin.Context) {
	logger.AppLogger().Debugf("GetNetworkConfig GET:%+v", c.Request)

	svc := networkservice.NewGetNetworkConfigService()
	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, nil))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, nil))
	}
}
