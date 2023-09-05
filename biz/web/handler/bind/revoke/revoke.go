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

package revoke

import (
	"agent/biz/model/dto/bind/revoke"
	serviceRevoke "agent/biz/service/bind/revoke"
	"agent/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindRevoke godoc
// @Summary admin unbind with server [for client bluetooth/LAN]
// @Description
// @ID BindRevoke
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Param  verifyReq body revoke.RevokeReq true  "params"
// @Success 200 {object} dto.BaseRspStr{results=revoke.RevokeRsp} "code=AG-200 success;"
// @Router /agent/v1/api/bind/revoke [POST]
func Revoke(c *gin.Context) {
	logger.AppLogger().Debugf("%+v", c.Request)

	var reqObject revoke.RevokeReq
	svc := new(serviceRevoke.RevokeService)
	c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
}
