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

package passthrough

import (
	"net/http"

	"agent/biz/model/passthrough"
	servicepassthrough "agent/biz/service/passthrough"

	"agent/utils/logger"

	"github.com/gin-gonic/gin"
)

// Passthrough godoc
// @Summary pass through gateway interface [for client bluetooth/LAN]
// @Description client pass through call gateway api
// @ID Passthrough
// @Tags Passthrough
// @Accept  json
// @Produce  json
// @Param   Request-Id      header     string     true  "request id "
// @Param   lanInvokeReq      body dto.LanInvokeReq true  "pass through params"
// @Success 200 {object} dto.BaseRspStr{results=string} "code=AG-200 success;"
// @Router /agent/v1/api/passthrough [POST]
func Passthrough(c *gin.Context) {
	logger.AppLogger().Debugf("Passthrough POST, req=%+v", c.Request)
	var reqObject passthrough.PassthroughReq
	svc := new(servicepassthrough.PassthroughService)
	c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
}
