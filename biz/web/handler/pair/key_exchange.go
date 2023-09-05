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
 * @Date: 2021-12-14 17:07:20
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-15 09:04:38
 * @Description:
 */

package pair

import (
	"agent/biz/model/dto"
	dtopair "agent/biz/model/dto/pair"
	servicepair "agent/biz/service/pair"
	"fmt"
	"net/http"

	"agent/utils/logger"
	"github.com/gin-gonic/gin"
)

// KeyExchange godoc
// @Summary exchange Symmetric key [for client]
// @Description client request server generate Symmetric key
// @ID KeyExchange
// @Tags Pair
// @Accept  json
// @Produce  json
// @Param  keyExchangeReq body dtopair.KeyExchangeReq true  "{clientPreSecret: client generate random string，32 byte,encBtid:String obtained by base64 after encrypting the btid with the public key of the server}"
// @Success 200 {object} dto.BaseRspStr{results=dtopair.KeyExchangeRsp} "code=AG-200 success;"
// @Router /agent/v1/api/keyexchange [POST]
func KeyExchange(c *gin.Context) {
	logger.AppLogger().Debugf("KeyExchange POST, req=%+v", c.Request)

	var reqObj dtopair.KeyExchangeReq
	if err := c.ShouldBindJSON(&reqObj); err != nil {
		err1 := fmt.Errorf("failed ShouldBindJSON, %+v", err)
		logger.AppLogger().Debugf("KeyExchange POST, %+v", err1)
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, Message: err1.Error()})
		return
	}
	rsp, _ := servicepair.ServiceKeyExchange(&reqObj)
	c.JSON(http.StatusOK, rsp)
}
