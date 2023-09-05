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
 * @LastEditTime: 2021-12-22 17:36:12
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

// Initial godoc
// @Summary init aospace bind progress [For Client]
// @Description query init progress ,return after done init progress
// @ID initial
// @Tags Pair
// @Accept  json
// @Produce  json
// @Param   passwordInfo      body dtopair.PasswordInfo true  "admin pass"
// @Success 200 {object} dto.BaseRspStr{results=call.MicroServerRsp} "code=AG-200 success."
// @Router /agent/v1/api/initial [POST]
func Initial(c *gin.Context) {
	logger.AppLogger().Debugf("initial POST:%+v", c.Request)

	var reqObj dtopair.PasswordInfo
	if err := c.ShouldBindJSON(&reqObj); err != nil {
		err1 := fmt.Errorf("failed ShouldBindJSON, %+v", err)
		logger.AppLogger().Debugf("revoke POST, %+v", err1)
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, Message: err1.Error()})
		return
	}

	rsp, _ := servicepair.ServiceInitial(&reqObj)
	c.JSON(http.StatusOK, rsp)
}
