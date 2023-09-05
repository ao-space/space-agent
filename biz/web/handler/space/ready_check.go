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

package space

import (
	servicespace "agent/biz/service/space"
	"net/http"

	"agent/utils/logger"
	"github.com/gin-gonic/gin"
)

// ReadyCheck godoc
// @Summary check if space is ready [for client bluetooth|LAN]
// @Description
// @ID ReadyCheck
// @Tags space
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=space.ReadyCheckRsp} "code=AG-200 success."
// @Router /agent/v1/api/space/ready/check [GET]
func ReadyCheck(c *gin.Context) {
	logger.AppLogger().Debugf("ReadyCheck GET:%+v", c.Request)
	svc := new(servicespace.ReadyCheckService)
	c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, nil))
}
