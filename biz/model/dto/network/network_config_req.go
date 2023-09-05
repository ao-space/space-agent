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

package network

type NetworkConfigReq struct {
	DNS1            string            `json:"dNS1"`            // ipv4 dNS1 地址
	DNS2            string            `json:"dNS2"`            // ipv4 dNS2 地址
	Ipv6DNS1        string            `json:"ipv6DNS1"`        // ipv6 dNS1 地址
	Ipv6DNS2        string            `json:"ipv6DNS2"`        // ipv6 dNS2 地址
	NetworkAdapters []*NetworkAdapter `json:"networkAdapters"` // 网络适配器列表
}
