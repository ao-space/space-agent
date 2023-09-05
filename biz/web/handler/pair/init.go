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
 * @Date: 2021-11-24 09:30:32
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-15 09:02:26
 * @Description:
 */

package pair

import (
	_ "agent/biz/model/dto"
	_ "agent/biz/model/dto/pair"
	servicepair "agent/biz/service/pair"
	"net/http"

	"agent/utils/logger"
	"github.com/gin-gonic/gin"
)

// PairInit godoc
// @Summary wired network binding [for client]
// @Description wired network binding
// @ID PairInit
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=pair.InitResult} "code=AG-200 success;"
// @Router /agent/v1/api/pair/init [GET]
func Init(c *gin.Context) {
	logger.AppLogger().Debugf("Init")
	rsp, _ := servicepair.ServiceInit()
	c.JSON(http.StatusOK, rsp)
}
