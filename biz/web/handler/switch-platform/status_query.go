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
 * @Author: zhongguang
 * @Date: 2022-11-23 09:36:53
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-24 17:58:36
 */

package switchplatform

import (
	"agent/biz/model/dto"
	"fmt"
	"net/http"

	modelsp "agent/biz/model/switch-platform"
	serviceswithplatform "agent/biz/service/switch-platform"

	"agent/utils/logger"

	"github.com/gin-gonic/gin"
)

// SwitchStatusQuery godoc
// @Summary switch platform
// @Description query switch platform status
// @ID SwitchStatusQuery
// @Tags Switch
// @Accept  json
// @Param   transId query string true  "switch platform params"
// @Success 200 {object} dto.BaseRspStr{results=modelsp.SwitchStatusQueryResp} "code=200 success;
// @Router /agent/v1/api/switch/status [GET]
// app 空间平台切换接口 v1
func SwitchStatusQuery(c *gin.Context) {

	logger.AppLogger().Debugf("SwitchPlatform GET:%+v", c.Request)

	var reqObj modelsp.SwitchStatusQueryReq
	if err := c.ShouldBind(&reqObj); err != nil {
		err1 := fmt.Errorf("failed ShouldBindQuery, %+v", err)
		logger.AppLogger().Debugf("SwitchStatusQuery GET, %+v", err1)
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, Message: err1.Error()})
		return
	}

	rsp, _ := serviceswithplatform.ServiceSwitchStatusQuery(&reqObj)
	c.JSON(http.StatusOK, rsp)

}
