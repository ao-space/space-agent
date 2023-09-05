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
 * @Date: 2021-11-16 17:13:15
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-17 17:22:56
 * @Description:
 */

package net

import (
	dtopair "agent/biz/model/dto/pair"
	servicepair "agent/biz/service/pair"
	"agent/config"
	"net/http"

	"agent/utils/logger"

	"github.com/gin-gonic/gin"
)

// NetConfig godoc
// @Summary get Wi-Fi list  [for client]
// @Description get Wi-Fi list
// @ID PairNetNetConfig
// @Tags net
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=[]dtopair.WifiListRsp} "code=AG-200 success;"
// @Router /agent/v1/api/pair/net/netconfig [GET]
func NetConfig(c *gin.Context) {
	logger.AppLogger().Debugf("NetConfig GET, c.Request.Host:%+v, DockerLocalListenAddr:%+v",
		c.Request.Host, config.Config.Web.DockerLocalListenAddr)

	var reqObj dtopair.WifiListReq
	rsp, _ := servicepair.ServiceWifiList(&reqObj, c.Request.Host != config.Config.Web.DockerLocalListenAddr)
	c.JSON(http.StatusOK, rsp)
}

// NetConfigDevice godoc
// @Summary get Wi-Fi list [for gateway]
// @Description get Wi-Fi list
// @ID PairNetNetConfigDevice
// @Tags device
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=[]dtopair.WifiListRsp} "code=AG-200 success;"
// @Router /agent/v1/api/device/netconfig [GET]
func NetConfigDevice(c *gin.Context) {
	logger.AppLogger().Debugf("NetConfig GET, c.Request.Host:%+v, DockerLocalListenAddr:%+v",
		c.Request.Host, config.Config.Web.DockerLocalListenAddr)

	var reqObj dtopair.WifiListReq
	rsp, _ := servicepair.ServiceWifiList(&reqObj, c.Request.Host != config.Config.Web.DockerLocalListenAddr)
	c.JSON(http.StatusOK, rsp)
}
