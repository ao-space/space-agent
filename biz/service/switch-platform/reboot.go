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
 * @Date: 2022-11-24 16:52:10
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-24 17:50:30
 */

package switchplatform

import (
	"agent/config"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

func RetryUnfinishedStatus() {
	//加载历史状态文件
	mtx.Lock()
	defer mtx.Unlock()

	var siLast StatusInfo
	err := fileutil.ReadFileJsonToObject(config.Config.Box.SwithStatusFile, &siLast)
	if err != nil {
		logger.AppLogger().Errorf("Read Switch-Platform file failed, file:%v, err:%v", config.Config.Box.BoxInfoFile, err)
		return
	}

	si = &siLast
	logger.AppLogger().Debugf("transid:%v, read Switch-Platform, status-info:%+v", si.TransId, siLast)

	switch si.Status {
	case StatusAbort, StatusOK, StatusInit:
		return
	}

	go doLastStatus()
}

func doLastStatus() {
	err := doStatusFlow(true)
	logger.AppLogger().Infof("Retry-switch-platform, transid:%v, status-info=%+v, err:%v", si.TransId, si, err)
}
