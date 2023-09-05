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
 * @Date: 2021-10-29 13:55:30
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-13 15:05:24
 * @Description:
 */
package pair

type TempKeyInfo struct {
	Key string `json:"key"`
	Iv  string `json:"iv"`
}

type WifiListReq struct {
	Count int `json:"count"`
}

type WifiListRsp struct {
	Name   string `json:"name"`
	Addr   string `json:"addr"`
	Signal int8   `json:"signal"`
}

type WifiPwdReq struct {
	Addr string `json:"addr"`
	Pwd  string `json:"pwd"`
}

type WifiStatusRsp struct {
	Name    string   `json:"name"`
	Addr    string   `json:"addr"`
	Status  int      `json:"status"`
	LocalIp []string `json:"ipAddress"`
}

type Network struct {
	Ip       string `json:"ip"`
	Wire     bool   `json:"wire"` // 有线
	WifiName string `json:"wifiName"`
	Port     uint16 `json:"port"`
	TlsPort  uint16 `json:"tlsPort"`
}
