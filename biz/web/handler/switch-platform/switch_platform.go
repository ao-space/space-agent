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

// SwitchPlatform godoc
// @Summary switch platform
// @Description switch platform to self-host platform
// @ID SwitchPlatform
// @Tags Switch
// @Accept  json
// @Produce  json
// @Param   SwitchPlatformReq body modelsp.SwitchPlatformReq true  "new platform params"
// @Success 200 {object} dto.BaseRspStr{results=modelsp.SwitchPlatformResp} "code=200 success;
// @Router /agent/v1/api/switch [POST]
func SwitchPlatform(c *gin.Context) {

	logger.AppLogger().Debugf("SwitchPlatform POST:%+v", c.Request)

	var reqObj modelsp.SwitchPlatformReq
	if err := c.ShouldBindJSON(&reqObj); err != nil {
		err1 := fmt.Errorf("failed ShouldBindJSON, %+v", err)
		logger.AppLogger().Debugf("SwitchPlatform POST, %+v", err1)
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, Message: err1.Error()})
		return
	}

	rsp, _ := serviceswithplatform.ServiceSwitchPlatform(&reqObj)
	c.JSON(http.StatusOK, rsp)

}
