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

package checkerimp

import (
	"agent/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/random"
)

const (
	ContainerNameGateway = "Gateway"
)

type GatewayAliveChecker struct {
}

func (checker *GatewayAliveChecker) Name() string {
	return ContainerNameGateway
}

func (checker *GatewayAliveChecker) Enable() bool {
	return config.Config.AliveChecker.GateWay.Enable
}

func (checker *GatewayAliveChecker) Check() bool {

	type Rsp struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}
	url := config.Config.AliveChecker.GateWay.UrlGateway
	var rsp Rsp
	_, err := GetJsonWithHeaders(url, nil, nil, &rsp)
	if err != nil {
		logger.CheckLogger().Warnf("failed check %v, err:%v", url, err)
		return false
	}
	logger.CheckLogger().Debugf("check %v return %+v", url, rsp)
	return strings.EqualFold(rsp.Status, "OK")
}

func (checker *GatewayAliveChecker) Restart() bool {
	// TODO: 调用 docker 模块来重启.
	// 目前处于开发阶段, 为了问题能够暴露并被发现来定位, 所以暂时没加重启逻辑.
	return true
}

func GetJsonWithHeaders(remoteUrl string, queryValues url.Values,
	headers map[string]string, ret interface{}) (*http.Request, error) {
	client := &http.Client{Timeout: time.Second * 3}

	if headers == nil {
		headers = map[string]string{"Request-Id": random.GenUUID()}
	}
	uri, err := url.Parse(remoteUrl)
	if err != nil {
		return nil, err
	}
	if queryValues != nil {
		values := uri.Query()
		for k, v := range values {
			queryValues[k] = v
		}
		uri.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed NewRequest, err:%v, uri.String():%+v", err, uri.String())
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cookie", "name=anny")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return req, fmt.Errorf("failed client.Do, err:%v, uri.String():%+v", err, uri.String())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return req, fmt.Errorf("failed ReadAll, err:%v", err)
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return req, fmt.Errorf("failed Unmarshal, err:%v, body:%v", err, string(body))
	}
	return req, nil
}
