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
 * @Date: 2022-11-15 16:07:59
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-15 16:08:22
 */

package switchplatform

import (
	"fmt"
	"time"

	"agent/utils/logger"
)

func networkDetect(transId string, domain string, info *ImigrateRsp) error {
	type StatusRsp struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}
	for i := 0; i < 60; i++ {
		time.Sleep(time.Second)
		rsp := StatusRsp{}

		url := "https://" + info.UserInfos[0].UserDomain + "/space/status"
		if err := httpGet(transId, url, &rsp); err != nil {
			logger.AppLogger().Debugf("network test, transId=%v, url=%v,  i=%v, err=%v",
				transId, url, i, err)
			continue
		} else if rsp.Status != "ok" {
			return fmt.Errorf("status:%v", rsp.Status)
		} else {
			return nil
		}
	}

	return fmt.Errorf("try too much times")
}
