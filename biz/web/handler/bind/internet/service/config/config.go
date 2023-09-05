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

package config

import (
	dtoconfig "agent/biz/model/dto/bind/internet/service/config"
	serviceConfig "agent/biz/service/bind/internet/service/config"
	"agent/config"
	"agent/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InternetServicePostConfig godoc
// @Summary config internet tunnel [for client/gateway bluetooth/LAN]
// @Description When the Internet tunnel is configured,
// @Description system-agent starts to register the server with the platform, and when it succeeds,
// @Description it calls the gateway interface, and the gateway registers the member information with the platform. And return the result to client.
// @ID InternetServicePostConfig
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Param   configReq      body dtoconfig.ConfigReq true  "config params"
// @Success 200 {object} dto.BaseRspStr{results=dtoconfig.ConfigRsp} "code=AG-200 success;"
// @Router /agent/v1/api/bind/internet/service/config [POST]
func PostConfig(c *gin.Context) {
	logger.AppLogger().Debugf("%+v", c.Request)

	var reqObject dtoconfig.ConfigReq
	svc := new(serviceConfig.InternetServiceConfig)

	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, &reqObject))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
	}

}

// InternetServiceGetConfig godoc
// @Summary get internet tunnel config [for client/gateway bluetooth/LAN]
// @Description
// @ID InternetServiceGetConfig
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Param   body query string true "clientUuid and aoid"
// @Success 200 {object} dto.BaseRspStr{results=dtoconfig.GetConfigRsp} "code=AG-200 success;"
// @Router /agent/v1/api/bind/internet/service/config [GET]
func GetConfig(c *gin.Context) {
	logger.AppLogger().Debugf("%+v", c.Request)

	var reqObject dtoconfig.GetConfigReq
	svc := new(serviceConfig.InternetServiceGetConfig)
	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, &reqObject))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
	}
}
