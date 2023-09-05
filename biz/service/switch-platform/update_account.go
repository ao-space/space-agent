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
 * @Date: 2022-11-21 15:07:01
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-21 18:00:34
 */
package switchplatform

import (
	"agent/biz/service/call"
	"agent/config"
	"fmt"
	"net/http"

	"agent/utils/logger"
	utilshttp "agent/utils/network/http"
)

func updateAccount() error {
	var mi MigrateInfo
	mi.UserInfos = si.ImigrateResult.UserInfos

	url := config.Config.Account.Migrate.Url
	logger.AppLogger().Debugf("updateAccount, transid=%v, url:%+v, parms:%+v", si.TransId, url, mi)

	var headers = map[string]string{"Request-Id": si.TransId}
	var rsp call.MicroServerRsp
	httpReq, httpRsp, body, err := utilshttp.SendJsonWithHeaders("PUT", url, mi, headers, &rsp)

	if err != nil {
		logger.AppLogger().Warnf("Failed PostJson, transid=%v, err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@body:%v", si.TransId, err, httpReq, httpRsp, string(body))
		return err
	}
	logger.AppLogger().Infof("updateAccount transid=%v,, parms:%+v", si.TransId, mi)
	logger.AppLogger().Infof("updateAccount transid=%v, rsp:%+v", si.TransId, rsp)
	logger.AppLogger().Infof("updateAccount transid=%v, httpReq:%+v", si.TransId, httpReq)
	logger.AppLogger().Infof("updateAccount transid=%v, httpRsp:%+v", si.TransId, httpRsp)
	logger.AppLogger().Infof("updateAccount transid=%v, body:%v", si.TransId, string(body))

	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf("http code:%v", httpRsp.StatusCode)
	}

	return nil
}
