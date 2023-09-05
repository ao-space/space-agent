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

//  获取网络信息: ;无线连上修改配置: ;  无线未连上修改配置:  ;  有线时修改配置:  ;
type NetworkStatusRsp struct {
	InternetAccess  bool              `json:"internetAccess"`  // 是否可以访问互联网 (获取网络信息: 返回;  其他: 不传;)
	DNS1            string            `json:"dNS1"`            // ipv4 dNS1 地址 (获取网络信息: 返回; 其他: 选传)
	DNS2            string            `json:"dNS2"`            // ipv4 dNS2 地址 (获取网络信息: 返回; 其他: 选传)
	LanAccessPort   uint16            `json:"lanAccessPort"`   // 局域网访问网关服务的端口号 (获取网络信息: 返回;  其他: 不传;)
	Ipv6DNS1        string            `json:"ipv6DNS1"`        // ipv6 dNS1 地址 (获取网络信息: 返回; 其他: 选传)
	Ipv6DNS2        string            `json:"ipv6DNS2"`        // ipv6 dNS2 地址 (获取网络信息: 返回; 其他: 选传)
	NetworkAdapters []*NetworkAdapter `json:"networkAdapters"` // 网络适配器列表 (获取网络信息: 返回; 其他: 必传)
}
