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
 * @Author: zhongguang
 * @Date: 2022-11-22 08:44:00
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-24 17:29:16
 */
package switchplatform

import (
	"agent/config"
	"fmt"
	"sync"

	"github.com/dungeonsnd/gocom/file/fileutil"
	"agent/utils/logger"
)

const (
	StatusInit          int = 0   //初始默认值
	StatusStart         int = 1   //状态管理开始
	StatusUpdateGateway int = 2   //更新网关信息完成
	StatusUpdateBoxInfo int = 3   //替换盒子本地信息完成
	StatusAbort         int = 99  //终止，可再次发起空间平台切换
	StatusOK            int = 100 //切换完成，域名重定向设置完成后
)

type StatusInfo struct {
	TransId        string
	Domain         string
	Status         int
	StatusMsg      string
	OldAccount     []AccountInfo
	OldApiBaseUrl  string
	ImigrateResult ImigrateRsp
	NewApiBaseUrl  string
}

var mtx sync.Mutex
var si *StatusInfo

func createStatus() (*StatusInfo, error) {
	mtx.Lock()
	defer mtx.Unlock()

	if si != nil && (si.Status != StatusAbort && si.Status != StatusOK) {
		return si, fmt.Errorf("task is doing")
	}

	si = &StatusInfo{}

	return si, nil
}

func UpdateStatus(status int, msg string) error {
	mtx.Lock()
	defer mtx.Unlock()

	si.Status = status
	si.StatusMsg = msg

	err := fileutil.WriteToFileAsJson(config.Config.Box.SwithStatusFile, si, "  ", true)
	if err != nil {
		logger.AppLogger().Errorf("Write StatusInfo file failed, file:%v, err:%v", config.Config.Box.SwithStatusFile, err)
	} else {
		logger.AppLogger().Debugf("Write StatusInfo file succ, file:%v, status info:%+v", config.Config.Box.SwithStatusFile, si)
	}
	logger.AppLogger().Debugf("Write StatusInfo file succ")

	return err
}
