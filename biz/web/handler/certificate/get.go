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

package certificate

import (
	"agent/biz/service/certificate"
	"agent/config"

	"net/http"

	"agent/utils/logger"
	"github.com/gin-gonic/gin"
)

// GetLanCert godoc
// @Summary get certificate(self sign) [for client /gateway]
// @Description get certificate(self sign)
// @ID GetLanCert
// @Tags cert
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=certificate.LanCert} "code=AG-200 success;"
// @Router /agent/v1/api/cert/get [GET]
func GetLanCert(c *gin.Context) {
	logger.AppLogger().Debugf("cert GET:%+v", c.Request)
	svc := new(certificate.LanCertService)
	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, nil))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, nil))
	}
}
