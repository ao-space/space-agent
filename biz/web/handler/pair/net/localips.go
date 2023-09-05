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
 * @LastEditTime: 2021-12-18 09:04:34
 * @Description:
 */
package net

import (
	"agent/biz/service/pair"
	"agent/biz/model/dto"
	"net/http"

	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/encrypt/random"
	"github.com/gin-gonic/gin"
)

// LocalIps godoc
// @Summary get server's ip in binding progress [for client]
// @Description client can get server's ip by using this api
// @ID PairNetLocalIps
// @Tags net
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=[]pair.Network} "code=AG-200 success;"
// @Router /agent/v1/api/pair/net/localips [GET]
func LocalIps(c *gin.Context) {
	rsp := &dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		RequestId: random.GenUUID(),
		Results:   pair.GetConnectedNetwork()}

	logger.AppLogger().Debugf("LocalIps return rsp:%+v, rsp.Results:%+v", rsp, rsp.Results)
	c.JSON(http.StatusOK, rsp)
}

// LocalIpsDevice godoc
// @Summary get server's ip [for gateway]
// @Description get server's ip
// @ID PairNetLocalIpsDevice
// @Tags device
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=[]pair.Network} "code=AG-200 success;"
// @Router /agent/v1/api/device/localips [GET]
func LocalIpsDevice(c *gin.Context) {
	rsp := &dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		RequestId: random.GenUUID(),
		Results:   pair.GetConnectedNetwork()}

	logger.AppLogger().Debugf("LocalIpsDevice return rsp:%+v, rsp.Results:%+v", rsp, rsp.Results)
	c.JSON(http.StatusOK, rsp)
}
