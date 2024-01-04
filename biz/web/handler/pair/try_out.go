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

package pair

import (
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/dto/pair/tryout"
	servicepair "agent/biz/service/pair"
	"agent/config"
	"fmt"
	"net/http"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/random"
	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/gin-gonic/gin"
)

// TryOutCode godoc
// @Summary verify trial code [for web client LAN]
// @Description
// @ID TryOutCode
// @Tags TryOut
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr "code=AG-200 success;"
// @Router /agent/v1/api/pair/tryout/code [POST]
func TryOutCode(c *gin.Context) {
	logger.AppLogger().Debugf("TryOutCode")

	tryoutCode := c.PostForm("tryoutCode")
	email := c.PostForm("email")

	var reqObj tryout.TryoutCodeReq
	if len(tryoutCode) < 1 || len(email) < 1 {
		err1 := fmt.Errorf("request params error, tryoutCode:%+v, email:%+v", tryoutCode, email)
		// logger.AppLogger().Debugf("TryOutCode POST, %+v", err1)
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, Message: err1.Error()})
		return
	}
	reqObj.TryoutCode = tryoutCode
	reqObj.Email = email

	abilityModel := device_ability.GetAbilityModel()
	if abilityModel.RunInDocker {
		err := fileutil.WriteToFile(config.Config.Box.HostIpFile, []byte(c.Request.Host), true)
		if err != nil {
			err1 := fmt.Errorf("failed write HostIpFile, %+v", err)
			logger.AppLogger().Debugf("info POST, %+v", err1)
			c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
				RequestId: random.GenUUID(),
				Message:   err1.Error()})
			return
		}
	}

	rsp, _ := servicepair.ServiceTryout(&reqObj)
	c.JSON(http.StatusOK, rsp)
}
