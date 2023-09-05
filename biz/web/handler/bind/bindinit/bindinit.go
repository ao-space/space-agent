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

package bindinit

import (
	"agent/biz/model/dto/bind/bindinit"
	servicesinit "agent/biz/service/bind/init"

	"net/http"

	"agent/utils/logger"

	"github.com/gin-gonic/gin"
)

// BindInit godoc
// @Summary APP get aospace server base info [for client bluetooth/LAN]
// @Description get aospace server base info.
// @ID BindInit
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Param   initReq      body bindinit.InitReq true  "query params"
// @Success 200 {object} dto.BaseRspStr{results=pair.InitResult} "code=AG-200 success;"
// @Router /agent/v1/api/bind/init [GET]
func Init(c *gin.Context) {
	logger.AppLogger().Debugf("%+v", c.Request)
	var reqObject bindinit.InitReq
	svc := new(servicesinit.InitService)
	c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
}
