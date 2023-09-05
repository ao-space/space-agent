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

package method

import (
	methodservice "agent/biz/service/did/document/method"
	"agent/config"
	"net/http"

	_ "agent/biz/model/dto/did/document"
	"agent/biz/model/dto/did/document/method"
	"agent/utils/logger"

	"github.com/gin-gonic/gin"
)

// UpdateDocumentMethod godoc
// @Summary get UpdateDocumentMethod [for gateway]
// @Description
// @ID UpdateDocumentMethod
// @Tags did
// @Produce  json
// @Param   updateDocumentMethodReq      body method.UpdateDocumentMethodReq true  "params"
// @Success 200 {object} dto.BaseRspStr{results=method.UpdateDocumentMethodRsp} "code=AG-200 成功."
// @Router /agent/v1/api/did/document/method [PUT]
func UpdateDocumentMethod(c *gin.Context) {
	logger.AppLogger().Debugf("GetDIDDoc GET:%+v", c.Request)

	var reqObject method.UpdateDocumentMethodReq

	svc := methodservice.NewUpdateDocumentMethod()
	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, &reqObject))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
	}
}
