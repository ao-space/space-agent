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

package password

import (
	passwordservice "agent/biz/service/did/document/password"
	"agent/config"
	"net/http"

	_ "agent/biz/model/dto/did/document"
	"agent/biz/model/dto/did/document/password"
	"agent/utils/logger"

	"github.com/gin-gonic/gin"
)

// UpdateDocumentPassword godoc
// @Summary get UpdateDocumentPassword [client call through gateway ,for gateway]
// @Description
// @ID UpdateDocumentPassword
// @Tags did
// @Produce  json
// @Param   updateDocumentPasswordReq      body password.UpdateDocumentPasswordReq true  "params"
// @Success 200 {object} dto.BaseRspStr "code=AG-200 success."
// @Router /agent/v1/api/did/document/password [PUT]
func UpdateDocumentPassword(c *gin.Context) {
	logger.AppLogger().Debugf("GetDIDDoc GET:%+v", c.Request)

	var reqObject password.UpdateDocumentPasswordReq

	svc := passwordservice.NewUpdateDocumentPassword()
	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, &reqObject))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
	}
}
