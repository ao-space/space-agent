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
 * @Date: 2021-10-30 20:30:23
 * @LastEditors: wenchao
 * @LastEditTime: 2021-11-09 08:55:01
 * @Description:数据模型定义
 */

package network

/*
nmcli device status
nmcli device show eth0
nmcli device show wlan0

[root@EulixOS ~]# nmcli  dev  show eth1
GENERAL.DEVICE:                         eth1
GENERAL.TYPE:                           ethernet
GENERAL.HWADDR:                         B2:A3:3B:4E:11:C4
GENERAL.MTU:                            1500
GENERAL.STATE:                          100 (connected)
GENERAL.CONNECTION:                     Wired connection 1
GENERAL.CON-PATH:                       /org/freedesktop/NetworkManager/ActiveConnection/1
WIRED-PROPERTIES.CARRIER:               on
IP4.ADDRESS[1]:                         192.168.124.112/24
IP4.GATEWAY:                            192.168.124.1
IP4.ROUTE[1]:                           dst = 0.0.0.0/0, nh = 192.168.124.1, mt = 100
IP4.ROUTE[2]:                           dst = 192.168.124.0/24, nh = 0.0.0.0, mt = 100
IP4.DNS[1]:                             192.168.124.1
IP6.ADDRESS[1]:                         fe80::358:9e52:1489:a770/64
IP6.GATEWAY:                            --
IP6.ROUTE[1]:                           dst = fe80::/64, nh = ::, mt = 100
*/

type NetDevice struct {
	GeneralDevice     string `json:"generalDevice"`
	GeneralHWAddr     string `json:"generalHWAddress"`
	GeneralType       string `json:"generalType"`
	GeneralConnection string `json:"generalConnection"`
	Ip4Address        string `json:"ip4Address"`
}

type ListWifiInfo struct {
	INUSE    bool   `json:"INUSE"`
	BSSID    string `json:"BSSID"`
	SSID     string `json:"SSID"`
	CHAN     string `json:"CHAN"`
	RATE     string `json:"RATE"`
	SIGNAL   string `json:"SIGNAL"`
	SECURITY string `json:"SECURITY"`
}

type DevStatus struct {
	DEVICE     string `json:"DEVICE"`
	TYPE       string `json:"TYPE"`
	STATE      string `json:"STATE"`
	CONNECTION string `json:"CONNECTION"`
}
