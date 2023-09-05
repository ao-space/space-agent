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
	"strconv"
	"strings"
)

const (
	DevStatus_Type_Wire     = "ethernet"
	DevStatus_Type_Wireless = "wifi"
)

type DevInfo struct {
	GENERALDEVICE     string `json:"GENERALDEVICE"`     // GENERAL.DEVICE:wlan0, GENERAL.DEVICE:eth0
	GENERALTYPE       string `json:"GENERALTYPE"`       // GENERAL.TYPE:wifi, GENERAL.TYPE:ethernet
	GENERALHWADDR     string `json:"GENERALHWADDR"`     // GENERAL.HWADDR:AE:FA:BF:E4:C7:26, GENERAL.HWADDR:10:2C:6B:7F:B1:78
	GENERALCONUUID    string `json:"GENERALCONUUID"`    // GENERAL.CON-UUID:3962b85f-3d69-496e-8974-6072229edfa2, GENERAL.CON-UUID:
	CAPABILITIESSPEED string `json:"CAPABILITIESSPEED"` // CAPABILITIES.SPEED:1000 Mb/s, CAPABILITIES.SPEED:unknown
}

type DevStatus struct {
	DEVICE     string `json:"DEVICE"`     // DEVICE:eth1, wlan0
	TYPE       string `json:"TYPE"`       // TYPE:ethernet, wifi, bridge
	STATE      string `json:"STATE"`      // STATE:connected, disconnected, unavailable
	CONNECTION string `json:"CONNECTION"` // CONNECTION:Wired connection 1
	CONUUID    string `json:"CONUUID"`    // CON-UUID:ddf82224-fcb2-3f61-88f8-4896ddd8f8f1
}

func (devStatus *DevStatus) IsEthernetAndWifi() bool {
	return devStatus.IsWire() || devStatus.IsWireless()
}

func (devStatus *DevStatus) IsWire() bool {
	return strings.Contains(devStatus.TYPE, DevStatus_Type_Wire)
}

func (devStatus *DevStatus) IsWireless() bool {
	return strings.Contains(devStatus.TYPE, DevStatus_Type_Wireless)
}

func (devStatus *DevStatus) IsConnected() bool {
	return strings.EqualFold(devStatus.STATE, "connected")
}

type ConInfo struct {
	ConnectionId            string `json:"connectionId"`            // connection.id:Wired connection 1
	ConnectionUuid          string `json:"connectionUuid"`          // connection.uuid:ddf82224-fcb2-3f61-88f8-4896ddd8f8f1
	ConnectionInterfaceName string `json:"connectionInterfaceName"` // connection.interface-name:eth1
	ConnectionAutoconnect   string `json:"connectionAutoconnect"`   // connection.autoconnect:yes
	Ipv4Method              string `json:"ipv4Method"`              // ipv4.method:auto, ipv4.method:manual (none)
	Ipv6Method              string `json:"ipv6Method"`              // ipv6.method:auto

	IP4ADDRESS string `json:"iP4ADDRESS"` // IP4.ADDRESS[1]:192.168.124.112/24
	IP4GATEWAY string `json:"iP4GATEWAY"` // IP4.GATEWAY:192.168.124.1
	IP4DNS1    string `json:"iP4DNS1"`    // IP4.DNS[1]:192.168.124.1
	IP4DNS2    string `json:"iP4DNS2"`    // IP4.DNS[2]:8.8.4.4
	IP6ADDRESS string `json:"iP6ADDRESS"` // IP6.ADDRESS[1]:fe80::9a44:b614:b787:83a0/64
	IP6GATEWAY string `json:"iP6GATEWAY"` // IP6.GATEWAY:

	SeenBssids  string `json:"seenBssids"`  // 802-11-wireless.seen-bssids:DC:FE:18:19:AA:8B
	GENERALNAME string `json:"generalName"` // GENERAL.NAME:JefEnterprise_5G
}

func (conInfo *ConInfo) UseDhcp() bool {
	return strings.EqualFold(conInfo.Ipv4Method, "auto")
}

func (conInfo *ConInfo) UseDhcpIpv6() bool {
	return strings.EqualFold(conInfo.Ipv6Method, "auto")
}

func (conInfo *ConInfo) Ipv4Address() string {
	arr := strings.Split(conInfo.IP4ADDRESS, "/")
	if len(arr) < 2 || len(arr[0]) < 1 || len(arr[1]) < 1 {
		return ""
	}
	return arr[0]
}

func (conInfo *ConInfo) Ipv4NetMask() string {
	arr := strings.Split(conInfo.IP4ADDRESS, "/")
	if len(arr) < 2 || len(arr[0]) < 1 || len(arr[1]) < 1 {
		return ""
	}
	l, err := strconv.Atoi(arr[1])
	if err != nil {
		return ""
	}

	return LenToSubNetMask(l)
}
