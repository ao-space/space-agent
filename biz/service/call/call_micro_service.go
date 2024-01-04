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
 * @Date: 2021-12-10 16:06:55
 * @LastEditors: wenchao
 * @LastEditTime: 2022-01-06 00:51:09
 * @Description:
 */

package call

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"agent/utils/logger"

	utilshttp "agent/utils/network/http"

	"github.com/dungeonsnd/gocom/encrypt/random"
)

func CallServiceByPost(url string, headers map[string]string, req, rsp interface{}) error {

	if headers == nil {
		headers = map[string]string{"Request-Id": random.GenUUID()}
	}
	httpReq, httpRsp, body, err1 := utilshttp.PostJsonWithHeaders(url, req, headers, rsp)
	if err1 != nil {
		// logger.AppLogger().Warnf("Failed CallServiceByPost, err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@body:%v", err1, httpReq, httpRsp, string(body))
		return err1
	}
	// logger.AppLogger().Debugf("CallServiceByPost, req:%+v", req)
	logger.AppLogger().Debugf("CallServiceByPost, rsp:%+v", rsp)
	logger.AppLogger().Infof("CallServiceByPost, httpReq:%+v", httpReq)
	logger.AppLogger().Infof("CallServiceByPost, httpRsp:%+v", httpRsp)
	logger.AppLogger().Debugf("CallServiceByPost, body:%v", string(body))
	return nil
}

func CallServiceByGet(url string, headers map[string]string, req, rsp interface{}) error {
	if headers == nil {
		headers = map[string]string{"Request-Id": random.GenUUID()}
	}
	_, _, _, err1 := utilshttp.SendJsonWithHeaders("GET", url, req, headers, rsp)
	if err1 != nil {
		return err1
	}
	return nil
}

func CallServiceByForm(method, url string, reqMap map[string]string, rsp interface{}) (*http.Response, error) {
	s := ""
	i := 0
	for k, v := range reqMap {
		s += fmt.Sprintf("%v=%v&", k, v)
		i++
	}
	if len(s) > 0 && s[len(s)-1:] == "&" {
		s = s[:len(s)-1]
	}
	return CallServiceByFormStr(method, url, s, rsp)
}

func CallServiceByFormStr(method, url string, reqStr string, rsp interface{}) (*http.Response, error) {
	if len(method) < 1 {
		method = "POST"
	}
	payload := strings.NewReader(reqStr)
	client := &http.Client{Timeout: time.Second * 20}

	// logger.AppLogger().Infof("CallServiceByForm, method:%+v, url:%+v, reqStr:%+v, payload:%+v", method, url, reqStr, payload)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("NewRequest err:%+v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Request-Id", random.GenUUID())
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do err:%+v", err)
	}
	defer resp.Body.Close()

	// logger.AppLogger().Infof("CallServiceByForm, httpReq:%+v", req)
	logger.AppLogger().Infof("CallServiceByForm, httpRsp:%+v", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("ReadAll err:%+v", err)
	}

	logger.AppLogger().Infof("CallServiceByForm, string(body):%+v", string(body))

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return resp, fmt.Errorf("failed Unmarshal, err:%v, body:%v", err, string(body))
	}

	logger.AppLogger().Debugf("CallServiceByPost, body:%v", string(body))

	return resp, nil
}

func DoHTTPRequest(method, url string, headers map[string]string, reqBody []byte, params map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.AppLogger().Errorf("do req error:%v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.AppLogger().Debugf("response body:%s", string(rspBody))
		return nil, err
	}
	return rspBody, nil
}
