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

/*
 * @Author: wenchao
 * @Date: 2021-11-10 14:41:01
 * @LastEditors: wenchao
 * @LastEditTime: 2021-11-30 15:15:38
 * @Description:
 */
package status

import (
	"agent/biz/model/dto"
	"agent/biz/model/dto/status"
	"agent/config"
	"net/http"

	"github.com/dungeonsnd/gocom/encrypt/random"
	"github.com/gin-gonic/gin"
)

// Status godoc
// @Summary get server status
// @Description get server status by this api.
// @ID status
// @Tags Info
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=status.Status} "code=AG-200 success"
// @Router /agent/status [GET]
func Status(c *gin.Context) {
	c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		RequestId: random.GenUUID(),
		Message:   "OK",
		Results:   &status.Status{Status: "OK", Version: config.Version}})
}
