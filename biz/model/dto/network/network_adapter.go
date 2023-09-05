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

type NetworkAdapter struct {
	AdapterName  string `json:"adapterName"`  // 网卡名称 eth0, eth1, wlan0 (获取网络信息: 返回;  其他: 不传;)
	Wired        bool   `json:"wired"`        // 有线还是无线网卡。 true: 有线; false: 无线。 (获取网络信息: 返回;  其他: 必传;)
	WIFIAddress  string `json:"wIFIAddress"`  // 路由器无线网络地址(不是盒子网卡地址)。有线连接时为空串。 (获取网络信息: 不返回;  无线连上修改配置: 不传;  无线未连上修改配置:  必传;  有线时修改配置:  不传;)
	WIFIName     string `json:"wIFIName"`     // WIFI名称。有线连接时为空串。 (获取网络信息: 返回;  无线连上修改配置: 必传;  无线未连上修改配置:  必传;  有线时修改配置:  不传;)
	WIFIPassword string `json:"wIFIPassword"` // WIFI密码。有线连接时为空串。 (获取网络信息: 不返回;  无线连上修改配置: 不传;  无线未连上修改配置:  必传;  有线时修改配置:  不传;)
	Connected    bool   `json:"connected"`    // 是否已连接 (获取网络信息: 返回;  其他: 不传;)
	MACAddress   string `json:"mACAddress"`   // 盒子网卡地址(不是路由器的网络地址)  (获取网络信息: 返回;  其他: 不传;)

	Ipv4UseDhcp    bool   `json:"ipv4UseDhcp"`    // ipv4 使用 dhcp 自动获取。 (获取网络信息: 返回;  其他: 必传;)
	Ipv4           string `json:"ipv4"`           // ipv4 地址 (获取网络信息: 返回;  无线连上修改配置: 必传;  无线未连上修改配置:  可选;  有线时修改配置:  必传;)
	SubNetMask     string `json:"subNetMask"`     // 子网掩码 (获取网络信息: 返回;  无线连上修改配置: 必传;  无线未连上修改配置:  可选;  有线时修改配置:  必传;)
	DefaultGateway string `json:"defaultGateway"` // 默认网关 (获取网络信息: 返回;  无线连上修改配置: 必传;  无线未连上修改配置:  可选;  有线时修改配置:  必传;)
	Ipv4DNS1       string `json:"ipv4DNS1"`       // ipv4 dNS1 地址 (获取网络信息: 返回; 其他: 选传)
	Ipv4DNS2       string `json:"ipv4DNS2"`       // ipv4 dNS2 地址 (获取网络信息: 返回; 其他: 选传)

	Ipv6UseDhcp        bool   `json:"ipv6UseDhcp"`        // ipv6 使用 dhcp 自动获取 (获取网络信息: 返回;  其他: 必传;)
	Ipv6               string `json:"ipv6"`               // ipv6  (获取网络信息: 返回;  无线连上修改配置: 必传;  无线未连上修改配置:  可选;  有线时修改配置:  必传;)
	SubNetPreLen       string `json:"subNetPreLen"`       // 子网前缀长度 (获取网络信息: 返回;  无线连上修改配置: 必传;  无线未连上修改配置:  可选;  有线时修改配置:  必传;)
	Ipv6DefaultGateway string `json:"ipv6DefaultGateway"` // ipv6 默认网关 (获取网络信息: 返回;  无线连上修改配置: 必传;  无线未连上修改配置:  可选;  有线时修改配置:  必传;)
}
