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

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func PostJsonWithHeaders(url string, parms interface{},
	headers map[string]string, ret interface{}) (*http.Request, *http.Response, []byte, error) {
	return SendJsonWithHeaders("POST", url, parms, headers, ret)
}

func SendJsonWithHeaders(method string, url string, parms interface{},
	headers map[string]string, ret interface{}) (*http.Request, *http.Response, []byte, error) {
	client := &http.Client{Timeout: time.Second * 20}

	// tlsConfig := &tls.Config{InsecureSkipVerify: true}
	// transport := &http.Transport{
	// 	TLSClientConfig:     tlsConfig,
	// 	MaxIdleConnsPerHost: 20,
	// }
	// http2.ConfigureTransport(transport)

	data, err := json.Marshal(parms)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed Marshal, err:%v, parms:%+v", err, parms)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed NewRequest, err:%v, data:%+v", err, string(data))
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cookie", "name=anny")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return req, resp, nil, fmt.Errorf("failed NewRequest, err:%v, data:%+v", err, string(data))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return req, resp, body, fmt.Errorf("failed ReadAll, err:%v", err)
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return req, resp, body, fmt.Errorf("failed Unmarshal, err:%v, body:%v", err, string(body))
	}
	return req, resp, body, nil
}

func PostJson(url string, parms interface{}, ret interface{}) (*http.Request, *http.Response, []byte, error) {
	return PostJsonWithHeaders(url, parms, nil, ret)
}

func PostJsonReturnMap(url string, parms interface{}) (interface{}, error) {
	client := &http.Client{Timeout: time.Second * 60}
	data, err := json.Marshal(parms)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var p interface{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
