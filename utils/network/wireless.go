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
// ## 手动设置 WIFI 网卡的ip、默认网关、DNS 等
// 之前已经连接 WIFI
// nmcli con mod b0764f9d-0dbc-4f66-a2af-e3ff2465027c ipv4.method manual ipv4.addr 192.168.0.203/24 ipv4.gateway 192.168.0.1
// nmcli con mod b0764f9d-0dbc-4f66-a2af-e3ff2465027c ipv4.dns "8.8.8.8 8.8.4.4"
// nmcli con up b0764f9d-0dbc-4f66-a2af-e3ff2465027c
///////////////////////////////////////////////////////////////////////////////////////////////

// 之前已经连接 WIFI
// nmcli con mod b0764f9d-0dbc-4f66-a2af-e3ff2465027c ipv4.method manual ipv4.addr 192.168.0.203/24 ipv4.gateway 192.168.0.1
func SetWirelessIpManual(connectionUuid string, ipv4 string, subnet string, gatewayv4 string) error {
	logger.AppLogger().Debugf("SetWirelessIpManual, connectionUuid:%v, ipv4:%v, subnet:%v, gatewayv4:%v",
		connectionUuid, ipv4, subnet, gatewayv4)
	prelen, err := SubNetMaskToLen(subnet)
	if err != nil {
		return err
	}
	ip4 := fmt.Sprintf("%v/%v", ipv4, prelen)
	params := []string{"con", "mod", connectionUuid, "ipv4.method", "manual", "ipv4.addr", ip4, "ipv4.gateway", gatewayv4}
	if err := runCmd(params); err != nil {
		return err
	}

	return nil
}

// nmcli con mod b0764f9d-0dbc-4f66-a2af-e3ff2465027c ipv4.dns "8.8.8.8 8.8.4.4"
func SetWirelessDnsManual(connectionUuid string, dns ...string) error {
	logger.AppLogger().Debugf("SetWirelessDnsManual, connectionUuid:%v, dns:%+v", connectionUuid, dns)
	dnsValid := []string{}
	for _, d := range dns {
		if len(d) > 0 {
			dnsValid = append(dnsValid, d)
		}
	}
	if len(dnsValid) < 1 {
		return nil
	}

	params := []string{"con", "mod", connectionUuid, "ipv4.dns", strings.Join(dnsValid, " ")}
	return runCmd(params)
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// ## 设置网卡 dhcp
// 无线时不能删除网卡配置，要不然 wifi 账号/密码也被删除了。
// nmcli con mod 247088ec-6cc6-44f9-85ae-cae29cce5d1a ipv4.method auto
// nmcli con mod 247088ec-6cc6-44f9-85ae-cae29cce5d1a ipv4.addresses "" ipv4.gateway "" ipv4.dns ""
// nmcli con down 247088ec-6cc6-44f9-85ae-cae29cce5d1a
// nmcli con up 247088ec-6cc6-44f9-85ae-cae29cce5d1a
// /////////////////////////////////////////////////////////////////////////////////////////////
func SetWirelessIpAuto(connectionUuid string) error {
	logger.AppLogger().Debugf("SetWirelessIpAuto, connectionUuid:%v", connectionUuid)
	params := []string{"con", "mod", connectionUuid, "ipv4.method", "auto"}
	if err := runCmd(params); err != nil {
		return err
	}

	// bash -c "nmcli con mod b0764f9d-0dbc-4f66-a2af-e3ff2465027c ipv4.addresses '' ipv4.gateway '' ipv4.dns ''"
	script := fmt.Sprintf("nmcli con mod %v ipv4.addresses '' ipv4.gateway '' ipv4.dns ''",
		connectionUuid)
	params = []string{"-c", script}
	if err := runCmd2("bash", params); err != nil {
		logger.AppLogger().Warnf("SetWirelessIpAuto, clear original ipv4.addresses failed. err:%v", err)
	}

	return nil
}
