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
	"agent/utils/logger"
	"regexp"
	"strings"
)

// device: eth0
// nmcli -t -m multiline --fields all dev show eth0
// nmcli -t -m multiline --fields all dev show wlan0
func GetNetworkDevInfo(device string) (*DevInfo, error) {
	logger.AppLogger().Debugf("GetNetworkDevInfo")
	params := []string{"-t", "-m", "multiline", "--fields", "all", "dev", "show", device}
	stdOutput, err := runCmdOutput(params)
	if err != nil {
		return nil, err
	}

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)
	devInfo := &DevInfo{}
	for _, line := range lines {
		// logger.AppLogger().Debugf("GetNetworkDevInfo, line(%v)=%+v", i, line)
		arr := strings.SplitN(line, ":", 2)
		if len(arr) != 2 {
			continue
		}
		k := arr[0]
		v := arr[1]

		if strings.EqualFold(k, "GENERAL.DEVICE") {
			devInfo.GENERALDEVICE = v
		} else if strings.EqualFold(k, "GENERAL.TYPE") {
			devInfo.GENERALTYPE = v
		} else if strings.EqualFold(k, "GENERAL.HWADDR") {
			devInfo.GENERALHWADDR = v
		} else if strings.EqualFold(k, "GENERAL.CON-UUID") {
			devInfo.GENERALCONUUID = v
		} else if strings.EqualFold(k, "CAPABILITIES.SPEED") {
			devInfo.CAPABILITIESSPEED = v
		}
	}

	logger.AppLogger().Debugf("GetNetworkDevInfo, return devInfo: %+v", devInfo)
	return devInfo, nil
}

// device: eth0
func GetNetworkDevStatusByDevice(device string, devStatuss []*DevStatus) *DevStatus {
	logger.AppLogger().Debugf("GetNetworkDevStatusByDevice")
	for i, devStatus := range devStatuss {
		logger.AppLogger().Debugf("GetNetworkDevStatusByDevice, devStatus(%v/%v): %+v", i, len(devStatuss), devStatus)
		if strings.EqualFold(devStatus.DEVICE, device) {
			logger.AppLogger().Debugf("GetNetworkDevStatusByDevice, device: %+v, devStatus: %+v", device, devStatus)
			return devStatus
		}
	}
	return nil
}

// nmcli -t -m multiline --fields all dev status
func GetNetworkDevStatus() ([]*DevStatus, error) {
	logger.AppLogger().Debugf("GetNetworkDevStatus")
	params := []string{"-t", "-m", "multiline", "--fields", "all", "dev", "status"}
	stdOutput, err := runCmdOutput(params)
	if err != nil {
		return nil, err
	}

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)
	devStatuss := []*DevStatus{}
	devStatus := &DevStatus{}
	for i, line := range lines {
		// logger.AppLogger().Debugf("GetNetworkDevStatus, line(%v/%v)=%+v", i, len(lines), line)
		arr := strings.SplitN(line, ":", 2)
		// for j, a := range arr {
		// 	logger.AppLogger().Debugf("GetNetworkDevStatus, arr[%v/%v]=%v", j, len(arr), a)
		// }
		if len(arr) != 2 {
			continue
		}
		k := arr[0]
		v := arr[1]
		if strings.EqualFold(k, "DEVICE") {
			if i > 0 {
				devStatuss = append(devStatuss, devStatus)
			}
			devStatus = &DevStatus{}
			devStatus.DEVICE = v
			// logger.AppLogger().Debugf("GetNetworkDevStatus, i:%v, devStatus.DEVICE= %+v, len(devStatuss):%v",
			// 	i, devStatus.DEVICE, len(devStatuss))
		} else if strings.EqualFold(k, "TYPE") {
			devStatus.TYPE = v
		} else if strings.EqualFold(k, "STATE") {
			devStatus.STATE = v
		} else if strings.EqualFold(k, "CONNECTION") {
			devStatus.CONNECTION = v
		} else if strings.EqualFold(k, "CON-UUID") {
			devStatus.CONUUID = v
		}

		// logger.AppLogger().Debugf("GetNetworkDevStatus, i:%+v, len(lines):%v", i, len(lines))
		if i == len(lines)-1 {
			// logger.AppLogger().Debugf("GetNetworkDevStatus, devStatus= %+v", devStatus)
			devStatuss = append(devStatuss, devStatus)
		}
	}

	for i, devStatus := range devStatuss {
		logger.AppLogger().Debugf("GetNetworkDevStatus, return devStatus(%v/%v): %+v", i, len(devStatuss), devStatus)
	}
	return devStatuss, nil
}

// nmcli -t connection show ddf82224-fcb2-3f61-88f8-4896ddd8f8f1
func GetNetworkConInfo(conUuid string) (*ConInfo, error) {
	logger.AppLogger().Debugf("GetNetworkConInfo")
	params := []string{"-t", "connection", "show", conUuid}
	stdOutput, err := runCmdOutput(params)
	if err != nil {
		return nil, err
	}

	IP4DNSArr := []string{}

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)
	conInfo := &ConInfo{}
	for _, line := range lines {
		// logger.AppLogger().Debugf("GetNetworkConInfo, line(%v)=%+v", i, line)
		arr := strings.SplitN(line, ":", 2)
		if len(arr) != 2 {
			continue
		}
		k := arr[0]
		v := arr[1]

		if strings.EqualFold(k, "connection.id") {
			conInfo.ConnectionId = v
		} else if strings.EqualFold(k, "connection.uuid") {
			conInfo.ConnectionUuid = v
		} else if strings.EqualFold(k, "connection.interface-name") {
			conInfo.ConnectionInterfaceName = v
		} else if strings.EqualFold(k, "connection.autoconnect") {
			conInfo.ConnectionAutoconnect = v
		} else if strings.EqualFold(k, "ipv4.method") {
			conInfo.Ipv4Method = v
		} else if strings.EqualFold(k, "ipv6.method") {
			conInfo.Ipv6Method = v
		} else if strings.Index(k, "IP4.ADDRESS") == 0 {
			conInfo.IP4ADDRESS = v
			logger.AppLogger().Debugf("[IP4.ADDRESS], k:%v, v:%v, line:%v", k, v, line)
		} else if strings.EqualFold(k, "IP4.GATEWAY") {
			conInfo.IP4GATEWAY = v
		} else if strings.Index(k, "IP6.ADDRESS") == 0 {
			conInfo.IP6ADDRESS = v
		} else if strings.EqualFold(k, "IP6.GATEWAY") {
			conInfo.IP6GATEWAY = v
		} else if strings.EqualFold(k, "802-11-wireless.seen-bssids") {
			conInfo.SeenBssids = v
		} else if strings.EqualFold(k, "GENERAL.NAME") {
			conInfo.GENERALNAME = v
		} else if strings.Index(k, "IP4.DNS") == 0 {
			IP4DNSArr = append(IP4DNSArr, v)
			logger.AppLogger().Debugf("[IP4DNSArr], k:%v, v:%v, line:%v", k, v, line)
		}

		if len(IP4DNSArr) >= 2 {
			conInfo.IP4DNS2 = IP4DNSArr[len(IP4DNSArr)-1]
			conInfo.IP4DNS1 = IP4DNSArr[len(IP4DNSArr)-2]
		} else if len(IP4DNSArr) >= 1 {
			conInfo.IP4DNS1 = IP4DNSArr[len(IP4DNSArr)-1]
		}

	}

	logger.AppLogger().Debugf("GetNetworkConInfo, return conInfo: %+v", conInfo)
	return conInfo, nil
}
