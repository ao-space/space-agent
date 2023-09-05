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

package progress

import (
	_ "agent/biz/model/dto/bind/com/progress"
	serviceProgress "agent/biz/service/bind/com/progress"
	"agent/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindComProgress godoc
// @Summary client query aospace microservice starting progress [for client bluetooth/LAN]
// @Description client query aospace microservice starting progress。
// @ID BindComProgress
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=progress.ProgressRsp} "code=AG-200 success; code=AG-460 have been bound; ProgressRsp: {"comStatus":1,"progress":100}, comStatus=0 starting, comStatus=1 started, comStatus=2 failed to start , comStatus<0 docker engine is starting"
// @Router /agent/v1/api/bind/com/progress [GET]
func Progress(c *gin.Context) {
	logger.AppLogger().Debugf("%+v", c.Request)

	svc := new(serviceProgress.ComProgressService)
	c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, nil))
}
