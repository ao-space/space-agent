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
 * @Date: 2021-12-01 09:33:36
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-22 15:23:54
 * @Description:
 */
package pair

import (
	"fmt"
	"net/http"

	"agent/biz/model/dto"
	dtopair "agent/biz/model/dto/pair"
	_ "agent/biz/service/call"
	servicepair "agent/biz/service/pair"

	"agent/utils/logger"
	"github.com/gin-gonic/gin"
)

// Pairing godoc
// @Summary aospace client bind to aospace server [for client]
// @Description
// @ID Pairing
// @Tags Pair
// @Accept  json
// @Produce  json
// @Param   pairingBoxInfo body dtopair.PairingReq true  "pair params"
// @Success 200 {object} dto.BaseRspStr{results=call.MicroServerRsp} "code=200 success; 410 clientUuid have been registered as member,it could not register as admin"
// @Router /agent/v1/api/pairing [POST]
func Pairing(c *gin.Context) {
	logger.AppLogger().Debugf("Pairing POST:%+v", c.Request)

	var reqObj dtopair.PairingReq
	if err := c.ShouldBindJSON(&reqObj); err != nil {
		err1 := fmt.Errorf("failed ShouldBindJSON, %+v", err)
		logger.AppLogger().Debugf("Pairing POST, %+v", err1)
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, Message: err1.Error()})
		return
	}

	rsp, _ := servicepair.ServicePairing(&reqObj)
	c.JSON(http.StatusOK, rsp)
}
