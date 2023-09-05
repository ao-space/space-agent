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
 * @Date: 2022-11-16 10:32:03
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-21 20:52:31
 */

package switchplatform

import (
	"agent/biz/service/call"
	"agent/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"agent/utils/logger"
)

type ClientInfo struct {
	ClientUUID string `json:"clientUUID"`
	ClientType string `json:"clientType"`
}
type AccountInfo struct {
	UserId      string       `json:"userId"`
	UserDomain  string       `json:"userDomain"`
	UserType    string       `json:"userType"`
	ClientInfos []ClientInfo `json:"clientInfos"`
}

type MigrateInfo struct {
	NetworkClinetId string        `json:"networkClientId,omitempty"`
	UserInfos       []AccountInfo `json:"userInfos"`
}

func httpGet(transId string, url string, rsp interface{}) error {

	if resp, err := http.Get(url); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("url=%v, code=%v", url, resp.StatusCode)
	} else {
		defer resp.Body.Close()

		if data, err := ioutil.ReadAll(resp.Body); err != nil {
			return err
		} else {

			logger.AppLogger().Debugf("transid=%v, test_data:%+v, ", transId, string(data))

			if err = json.Unmarshal(data, rsp); err != nil {
				return err
			}

			logger.AppLogger().Debugf("transid=%v, rsp:%+v, ", transId, rsp)
		}
	}
	return nil
}

func httpGetWithHeaders(url string, headers map[string]string, ret interface{}) error {
	client := &http.Client{Timeout: time.Second * 60}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("errcode:%v, body:%v", resp.StatusCode, string(body))
	}

	logger.AppLogger().Debugf("url:%v, status:%v, body:%v", url, resp.StatusCode, string(body))

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return fmt.Errorf("failed Unmarshal, err:%v, body:%v", err, string(body))
	}
	return nil
}

func getAccountInfos(transId string) ([]AccountInfo, error) {
	migration := &MigrateInfo{}
	rsp := &call.MicroServerRsp{}
	rsp.Results = migration
	if err := httpGetWithHeaders(config.Config.Account.Migrate.Url, map[string]string{"Request-Id": transId}, rsp); err != nil {
		return nil, err
	}

	return migration.UserInfos, nil
}
