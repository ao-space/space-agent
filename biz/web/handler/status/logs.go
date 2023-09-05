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
 * @LastEditTime: 2021-11-17 09:53:29
 * @Description:
 */
package status

import (
	"agent/config"
	"net/http"
	"path"

	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/gin-gonic/gin"
)

// Logs godoc
// @Summary get server info
// @Description get server info
// @ID logs
// @Tags Info
// @Accept  plain
// @Produce  plain
// @Success 200 {string} string "code=0 success;"
// @Router /agent/logs [GET]
func Logs(c *gin.Context) {

	content, err := fileutil.ReadFromFile(path.Join(config.Config.Log.Path, config.Config.Log.Filename) + ".log")
	if err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, string(content))
	}
}
