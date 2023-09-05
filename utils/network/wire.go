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

import (
	"fmt"
	"strings"

	"agent/utils/logger"
)

///////////////////////////////////////////////////////////////////////////////////////////////
// ## 手动设置网卡的ip、默认网关、DNS 等
// nmcli c del eth1
// nmcli con add type ethernet con-name eth1 ifname eth1 ip4 192.168.0.213/24 gw4 192.168.0.1
// nmcli con mod eth1 ipv4.dns "8.8.8.8 8.8.4.4"
// nmcli con up eth1
///////////////////////////////////////////////////////////////////////////////////////////////

// nmcli c del eth1
// nmcli con add type ethernet con-name eth1 ifname eth1 ip4 192.168.0.213/24 gw4 192.168.0.1
func SetWireIpManual(device string, ipv4 string, subnet string, gatewayv4 string) error {
	logger.AppLogger().Debugf("SetWireIpManual, device:%v, ipv4:%v, subnet:%v, gatewayv4:%v", device, ipv4, subnet, gatewayv4)
	params := []string{"c", "del", device}

	if err := runCmd(params); err != nil {
		logger.AppLogger().Warnf("SetWireIpManual, err:%v", err)
		// return err // 原来没有手动设置过 eth1, 不需要删除。
	}

	prelen, err := SubNetMaskToLen(subnet)
	if err != nil {
		return err
	}
	ip4 := fmt.Sprintf("%v/%v", ipv4, prelen)
	params = []string{"con", "add", "type", "ethernet", "con-name", device, "ifname", device, "ip4", ip4, "gw4", gatewayv4}
	if err := runCmd(params); err != nil {
		return err
	}

	return nil
}

// nmcli con mod static2 ipv4.dns "8.8.8.8 8.8.4.4"
// or
// nmcli con mod eth1 ipv4.dns "8.8.8.8 8.8.4.4"
func SetWireDnsManual(device string, dns ...string) error {
	logger.AppLogger().Debugf("SetWireDnsManual, device:%v, dns:%+v", device, dns)
	dnsValid := []string{}
	for _, d := range dns {
		if len(d) > 0 {
			dnsValid = append(dnsValid, d)
		}
	}
	if len(dnsValid) < 1 {
		return nil
	}
	params := []string{"con", "mod", device, "ipv4.dns", strings.Join(dnsValid, " ")}
	return runCmd(params)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// // nmcli con mod eth1 ipv4.method auto # 直接改的话原来的 ip 还在，会导致有多个 ip. 所有还是先删除再创建比较好。
// nmcli c del eth1
// nmcli con add type ethernet con-name eth1 ifname eth1 ipv4.method auto
// nmcli con down eth1 # 这2行貌似多余?
// nmcli con up eth1
///////////////////////////////////////////////////////////////////////////////////////////////

// nmcli con mod eth1 ipv4.method auto
func SetWireIpAuto(device string) error {
	// 直接改的话原来的 ip 还在，会导致有多个 ip. 所有还是先删除再创建比较好。
	// logger.AppLogger().Debugf("SetWireIpAuto, device:%v", device)
	// params := []string{"con", "mod", device, "ipv4.method", "auto"}
	// if err := runCmd(params); err != nil {
	// 	return err
	// }

	params := []string{"c", "del", device}

	if err := runCmd(params); err != nil {
		logger.AppLogger().Warnf("SetWireIpAuto, err:%v", err)
	}

	params = []string{"con", "add", "type", "ethernet", "con-name", device, "ifname", device, "ipv4.method", "auto"}
	if err := runCmd(params); err != nil {
		return err
	}

	return nil
}
