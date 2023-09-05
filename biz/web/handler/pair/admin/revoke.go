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
 * @Date: 2021-11-01 13:31:24
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-22 14:23:07
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

// Revoke godoc
// @Summary admin unbinding from the server [for client]
// @Description
// @ID PairAdminRevoke
// @Tags Pair
// @Accept  json
// @Produce  json
// @Param   revokeReq body pair.RevokeReq true  "server admin pass"
// @Success 200 {object} dto.BaseRspStr{results=call.MicroServerRsp}
// @Router /agent/v1/api/admin/revoke [POST]
func Revoke(c *gin.Context) {
	logger.AppLogger().Debugf("revoke POST, req=%+v", c.Request)
	var reqObj dtopair.RevokeReq
	if err := c.ShouldBindJSON(&reqObj); err != nil {
		err1 := fmt.Errorf("failed ShouldBindJSON, %+v", err)
		logger.AppLogger().Debugf("revoke POST, %+v", err1)
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, Message: err1.Error()})
		return
	}

	rsp, _ := servicepair.ServiceRevoke(&reqObj)
	c.JSON(http.StatusOK, rsp)
}
