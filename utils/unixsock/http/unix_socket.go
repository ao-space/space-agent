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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"agent/utils/logger"
)

func GetJsonWithHeadersByUnixSock(sockAddr, path string, parms interface{},
	headers map[string]string, ret interface{}) (*http.Response, error) {
	logger.AppLogger().Debugf("GetJsonWithHeadersByUnixSock, path:%v", path)

	// fmt.Println("Unix HTTP client")
	httpc := http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockAddr)
			},
		},
	}
	response, err := httpc.Get(fmt.Sprintf("http://unix.sock%v", path))
	if err != nil {
		logger.AppLogger().Debugf("GetJsonWithHeadersByUnixSock, path:%v, err:%v", path, err)
		return response, err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return response, err
		}
		if ret != nil {
			err = json.Unmarshal(body, &ret)
			if err != nil {
				return response, err
			}
		}
	} else {
		return response, fmt.Errorf("http StatusCode : %v", response.StatusCode)
	}
	return response, nil
}

func PostJsonWithHeadersByUnixSock(sockAddr, path string, parms interface{},
	headers map[string]string, ret interface{}) (*http.Request, *http.Response, []byte, error) {
	logger.AppLogger().Debugf("PostJsonWithHeadersByUnixSock, path:%v", path)

	client := &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockAddr)
			},
		}}

	data, err := json.Marshal(parms)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed Marshal, err:%v, parms:%+v", err, parms)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://unix.sock%v", path), bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed NewRequest, err:%v, data:%+v", err, string(data))
	}

	req.Header.Set("Content-Type", "application/json")
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

func SendMessageToSocket(sockAddr, message string) error {
	var conn net.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("unix", sockAddr)
		if err != nil {
			logger.UpgradeLogger().Errorf("failed to connect to socket,%s", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	defer conn.Close()

	// Send a command to the aospace-upgrade

	conn.Write([]byte(message))

	// Read the response from the aospace-upgrade server
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		logger.UpgradeLogger().Errorf("read response message error:%v", err)
		return err
	}
	response := string(buf[:n])
	logger.UpgradeLogger().Infof("get response message :%v", response)
	return nil
}
