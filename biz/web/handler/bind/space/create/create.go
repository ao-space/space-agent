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

package create

import (
	"agent/biz/model/dto/bind/space/create"
	serviceCreate "agent/biz/service/bind/space/create"
	"agent/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SpaceCreate godoc
// @Summary create aospace account(finished bind progress) [for client bluetooth/LAN]
// @Description This interface will submit the form data for space configuration on client,
// @Description such as space identification information, administrator space password and server Internet service configuration, etc.
// @Description When client calls this interface, it will synchronize the binding success status with server and store the status data of both sides. server will return the access_token to client.
// @Description 解绑后重新绑定情况下，创建空间时不传 SpaceName 和 EnableInternetAccess，传 Password 参数。
// @ID SpaceCreate
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Param   createReq      body create.CreateReq true  "admin password"
// @Success 200 {object} dto.BaseRspStr{results=create.CreateRsp} "code=AG-200 success; AG-577 register failed on self-host platform; AG-560 failed to create member"
// @Router /agent/v1/api/bind/space/create [POST]
func Create(c *gin.Context) {
	logger.AppLogger().Debugf("%+v", c.Request)

	var reqObject create.CreateReq
	svc := new(serviceCreate.SpaceCreateService)
	c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &reqObject))
}
