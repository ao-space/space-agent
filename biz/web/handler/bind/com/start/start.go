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

package start

import (
	serviceStart "agent/biz/service/bind/com/start"
	"agent/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindComStart godoc
// @Summary run docker containers [for client bluetooth/LAN]
// @Description APP start to pull and run aospace mircoservice。
// @ID BindComStart
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr "code=AG-200 success; code=AG-460 have been bound; code=AG-470 containers is starting; code=AG-471 containers is started; "
// @Router /agent/v1/api/bind/com/start [POST]
func Start(c *gin.Context) {
	logger.AppLogger().Debugf("%+v", c.Request)

	svc := new(serviceStart.ComStartService)
	c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, nil))
}
