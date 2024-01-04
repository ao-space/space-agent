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
 * @Date: 2021-10-30 17:54:56
 * @LastEditors: jeffery
 * @LastEditTime: 2022-04-13 14:37:04
 * @Description: 管理树莓派的 WIFI 连接
 */

package network

import (
	"agent/utils/tools"
	"fmt"
	"regexp"
	"strings"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/sys/run"
)

const (
	NMCLI_OUTPUT_CONNECT_EN  = "connected"
	NMCLI_OUTPUT_EXTERNAL_EN = "(externally)"
	NMCLI_OUTPUT_CONNECT_ZH  = "已连接"
	NMCLI_OUTPUT_EXTERNAL_ZH = "（外部）"
)

func GetDefaultGateway(adapterName string) string {
	logger.AppLogger().Debugf("GetDefaultGateway, adapterName:%v", adapterName)
	cmd, parms := "route", []string{"-n"}
	_, output, err := tools.RunCmd(cmd, parms)
	if err != nil {
		logger.AppLogger().Warnf("RunCmd %v %v, err:%v", cmd, strings.Join(parms, " "))
		return ""
	}

	Gateway := ""
	Destination := ""
	Metric := ""
	lines := tools.StringToLines(string(output))
	for _, line := range lines {
		logger.AppLogger().Debugf("line:%v", line)
		fieldsArr := strings.Fields(line)
		if len(fieldsArr) < 8 {
			continue
		}

		// 也许与网卡关系不大，先不判断了
		// Iface := fieldsArr[7]
		// if !strings.EqualFold(Iface, adapterName) {
		// 	continue
		// }

		if strings.EqualFold(fieldsArr[3], "UG") {
			if len(Gateway) > 0 {
				if strings.EqualFold(Destination, "0.0.0.0") && strings.EqualFold(fieldsArr[0], "0.0.0.0") {
					if strings.Compare(Metric, fieldsArr[4]) > 0 {
						Destination = fieldsArr[0]
						Gateway = fieldsArr[1]
						Metric = fieldsArr[4]
					}

				} else if strings.EqualFold(fieldsArr[0], "0.0.0.0") {
					Destination = fieldsArr[0]
					Gateway = fieldsArr[1]
					Metric = fieldsArr[4]
				}
			} else {
				Destination = fieldsArr[0]
				Gateway = fieldsArr[1]
				Metric = fieldsArr[4]
			}
		}
	}

	return Gateway
}

/**
 * @Title: Ping
 * @Description: 测试能否ping通一个地址,  ping -c 3 192.168.0.112
 * @Author: wenchao
 * @Date: 2021-10-30 18:10:59
 * @param {string} host 目标地址
 * @return {bool} 能否ping通
 * @return {error} 执行失败
 */
func Ping(host string) (bool, error) {
	params := []string{"-c", "3", host}
	logger.AppLogger().Debugf("Ping, run cmd: ping %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("ping", params)
	// fmt.Printf("Ping, stdOutput:\n%v\n\nerrOutput:%v\n\nerr:\n\n%v\n\n",
	// 	string(stdOutput), string(errOutput), err)
	if err != nil {
		return false, fmt.Errorf("failed run Ping %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("Ping, run cmd: ping %v, , stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))

	return true, err
}

/**
 * @Title:
 * @Description: 搜索附近的 wifi, nmcli dev wifi list --rescan yes
 * @Author: wenchao
 * @Date: 2021-10-31 10:28:39
 * @param {*}
 * @return {[]*ListWifiInfo} wifi 列表
 * @return {error} 错误
 */
func ListWifi() ([]*ListWifiInfo, error) {
	params := []string{"dev", "wifi", "list", "--rescan", "yes"}
	logger.AppLogger().Debugf("ListWifi, run cmd: nmcli %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("nmcli", params)
	if err != nil {
		return []*ListWifiInfo{}, fmt.Errorf("failed run ListWifi %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("ListWifi, run cmd: nmcli %v, , stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)
	wifiInfos := []*ListWifiInfo{}
	for row, v := range lines {
		wifiInfos = parseWifiLine(row, v, wifiInfos)
	}

	return wifiInfos, nil
}

func parseWifiListOneRowFields(rowString string) []string {
	fields := strings.Fields(rowString)
	newFields := []string{}
	if len(fields) > 0 {
		modeIdx := 2
		if strings.Contains(fields[0], "*") {
			modeIdx = 3
		}

		ssidHasWhiteSpace := false
		for i, v := range fields {
			if strings.EqualFold(v, "Infra") {
				if i != modeIdx {
					// SSID has white space
					ssidHasWhiteSpace = true
				}
			}
		}

		if ssidHasWhiteSpace {
			for i, v := range fields {
				if i <= 1 && strings.Count(v, ":") > 4 {
					for j, w := range fields {
						if strings.EqualFold(w, "Infra") {
							if j-i > 1 {
								idx1 := strings.Index(rowString, fields[i])
								idx2 := strings.Index(rowString, fields[j])
								if len(rowString) > idx1+len(fields[i]) && len(rowString) > idx2 {
									wifiField := rowString[idx1+len(fields[i]) : idx2]
									wifiField = strings.TrimSpace(wifiField)

									newFields = append(fields[:i+1], wifiField)
									newFields = append(newFields, fields[j:]...)
								}
							}
						}

					}

				}
			}
		}

	}
	if len(newFields) > 0 {
		return newFields
	}
	return fields
}

func parseWifiLine(row int, v string, wifiInfos []*ListWifiInfo) []*ListWifiInfo {

	// fmt.Printf("ListWifi, ####, v:%+v\n\n", v)
	if row == 0 {
		return wifiInfos // 第一行是表头, 忽略
	}

	// filedsRg := regexp.MustCompile(`[\t]`)
	// fileds := filedsRg.Split(v, -1)
	fileds := parseWifiListOneRowFields(v)
	// fmt.Printf("ListWifi, len(fileds):%+v\n\n", len(fileds))
	if len(fileds) < 4 {
		return wifiInfos
	}
	wifiInfo := &ListWifiInfo{INUSE: false}
	for idx, field := range fileds {
		wifiInfo = parseWifiField(idx, field, wifiInfo)
	}
	// fmt.Printf("\n\n")

	if strings.EqualFold(wifiInfo.SSID, "--") {
		return wifiInfos
	}
	wifiInfos = append(wifiInfos, wifiInfo)
	return wifiInfos
}

func parseWifiField(idx int, field string, wifiInfo *ListWifiInfo) *ListWifiInfo {

	// fmt.Printf("[%+v] ", field)
	if idx == 0 && strings.Contains(field, "*") {
		wifiInfo.INUSE = true
	}

	if wifiInfo.INUSE {
		if idx == 1 {
			wifiInfo.BSSID = field
		} else if idx == 2 {
			wifiInfo.SSID = field
		} else if idx == 5 {
			wifiInfo.CHAN = field
		} else if idx == 6 {
			wifiInfo.RATE = field
		} else if idx == 7 {
			wifiInfo.SIGNAL = field
		} else if idx == 9 {
			wifiInfo.SECURITY = field
		}
	} else {
		if idx == 0 {
			wifiInfo.BSSID = field
		} else if idx == 1 {
			wifiInfo.SSID = field
		} else if idx == 4 {
			wifiInfo.CHAN = field
		} else if idx == 5 {
			wifiInfo.RATE = field
		} else if idx == 6 {
			wifiInfo.SIGNAL = field
		} else if idx == 8 {
			wifiInfo.SECURITY = field
		}
	}

	return wifiInfo
}

/**
 * @Title: ConnectWifi
 * @Description: 连接 wifi, nmcli dev wifi connect 803b password pwd803b   (搜索附近的 wifi, nmcli dev wifi list --rescan yes)
 * @Description: 断开wlan0连接, nmcli dev disconnect wlan0   , 忘记一个保存的连接 nmcli c del 803b-b
 * @Author: wenchao
 * @Date: 2021-10-30 19:25:51
 * @param {string} name wifi名称
 * @param {string} pwd wifi密码
 * @return {bool} 连接成功
 * @return {string} 本机ip地址
 * @return {error} 执行失败
 */
func ConnectWifi(BSSID, PWD string) (bool, error) {
	succ := false
	var err error
	// 尝试若干次
	tryTtl := 2
	for i := 0; i < tryTtl; i++ {
		var params []string
		if i < 2 {
			params = []string{"dev", "wifi", "connect", BSSID, "password", PWD}
		} else {
			params = []string{"dev", "wifi", "connect", BSSID, "password", `"` + PWD + `"`}
		}

		logger.AppLogger().Debugf("ConnectWifi (%v/%v), run cmd: nmcli", i+1, tryTtl)
		stdOutput, errOutput, err1 := run.RunExe("nmcli", params)
		// fmt.Printf("ConnectWifi, stdOutput:\n%v\n\nerrOutput:%v\n\nerr:\n\n%v\n\n",
		// 	string(stdOutput), string(errOutput), err)
		if err1 != nil {
			err = fmt.Errorf("failed run ConnectWifi(%v/%v) %v, err is :%v, stdOutput is :%v, errOutput is :%v",
				i+1, tryTtl, params, err, string(stdOutput), string(errOutput))
			time.Sleep(2 * time.Second)
		} else {
			logger.AppLogger().Debugf("run ConnectWifi(%v/%v), run cmd: nmcli, stdOutput is :%v, errOutput is :%v",
				i+1, tryTtl, string(stdOutput), string(errOutput))

			succ = strings.Contains(string(stdOutput), "successfully activated")
		}

		if err == nil && succ {
			break
		}
	}

	return succ, err
}

// 忘记 wifi
// nmcli c del 803b-c
func ForgetWifi(wifiName string) (bool, error) {
	params := []string{"c", "del", wifiName}
	logger.AppLogger().Debugf("ForgetWifi, run cmd: nmcli %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("nmcli", params)
	if err != nil {
		return false, fmt.Errorf("failed run ForgetWifi %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("ForgetWifi, run cmd: nmcli %v, stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))

	return strings.Contains(string(stdOutput), "successfully deleted"), nil
}

// nmcli connection show 803b-b
func GetConnectedWirelessBssids(wifiName string) (string, error) {
	params := []string{"connection", "show", wifiName}
	logger.AppLogger().Debugf("GetConnectedWirelessBssids, run cmd: nmcli %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("nmcli", params)
	if err != nil {
		return "", fmt.Errorf("failed run GetConnectedWirelessBssids %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("GetConnectedWirelessBssids, run cmd: nmcli %v,errOutput is :%v",
		strings.Join(params, " "), string(errOutput))

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)
	key := "seen-bssids:"
	for _, line := range lines { // 802-11-wireless.seen-bssids:            14:DE:39:A6:7C:48
		i1 := strings.Index(line, key)
		if i1 > 0 {
			logger.AppLogger().Debugf("GetConnectedWirelessBssids, line : [%v]", line)
			v := line[i1+len(key):]
			logger.AppLogger().Debugf("GetConnectedWirelessBssids, v : [%v]", v)
			v = strings.TrimSpace(v)
			if strings.Index(v, ",") > 0 { // 802-11-wireless.seen-bssids:            14:DE:39:A6:7C:48,40:FE:95:00:12:66
				arr := strings.Split(v, ",")
				if len(arr) > 0 {
					return arr[0], nil
				} else {
					return "", fmt.Errorf("v:%v", v)
				}
			}
			return strings.TrimSpace(v), nil
		}
	}

	return "", fmt.Errorf("not found %v", key)
}

// GetIpAddress
/*
 * @Description: 获取本机 ip, nmcli dev status, nmcli dev show eth0, nmcli dev show wlan0
 * @Author: wenchao
 * @Date: 2021-10-30 19:27:27
 * @param {*}
 * @return {[]string} ip地址, 有线连接, 无线连接的顺序
 * @return {error} 执行失败
 */
func GetIpAddress() ([]*NetDevice, error) {
	params := []string{"dev", "status"}
	logger.AppLogger().Debugf("GetIpAddress, run cmd: nmcli %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("nmcli", params)
	if err != nil {
		return []*NetDevice{}, fmt.Errorf("failed run GetIpAddress %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("GetIpAddress, run cmd: nmcli %v,errOutput is :%v",
		strings.Join(params, " "), string(errOutput))

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)
	logger.AppLogger().Debugf("lines=%+v", lines)

	devices := []string{}
	for row, v := range lines {
		devices = parseIpAddressLine(row, v, devices)
	}

	ret := []*NetDevice{}
	for idx, v := range devices {
		logger.AppLogger().Debugf("GetIpAddress, devices[%v]=%+v", idx, v)
		deviceInfo, err1 := deviceShow(v)
		if err1 != nil {
			return nil, err1
		}
		logger.AppLogger().Debugf("GetIpAddress, idx=%v, deviceInfo=%+v", idx, deviceInfo)
		if strings.EqualFold(deviceInfo.GeneralType, "ethernet") ||
			strings.EqualFold(deviceInfo.GeneralType, "wifi") {
			ret = append(ret, deviceInfo)
		}
	}

	return ret, nil
}

func parseIpAddressLine(row int, v string, devices []string) []string {
	if row == 0 {
		return devices // 第一行是表头, 忽略
	}

	fields := strings.Fields(v)

	if !arrayHas(fields, NMCLI_OUTPUT_CONNECT_EN) && !arrayHas(fields, NMCLI_OUTPUT_CONNECT_ZH) {
		return devices
	}

	if arrayHas(fields, NMCLI_OUTPUT_EXTERNAL_EN) || arrayHas(fields, NMCLI_OUTPUT_EXTERNAL_ZH) {
		return devices
	}

	for idx, field := range fields {
		if idx == 0 {
			devices = append(devices, field)
		}
	}
	return devices
}

func arrayHas(arr []string, target string) bool {
	for _, v := range arr {
		// logger.AppLogger().Debugf("#### arrayHas, field: %v, target: %v\n", v, target)
		if v == target {
			// logger.AppLogger().Debugf("#### arrayHas, field: %v == target: %v\n", v, target)
			return true
		}
	}
	return false
}

func devShowFileds(s string) []string {
	// logger.AppLogger().Debugf("devShowFileds, s=%v", s)
	ret := make([]string, 0)
	i := strings.Index(s, ":")
	if i > 0 {
		s1 := s[:i]
		s2 := s[i+1:]
		// logger.AppLogger().Debugf("devShowFileds, s1=%v", s1)
		// logger.AppLogger().Debugf("devShowFileds, s2=%v", s2)
		ret = append(ret, strings.TrimSpace(s1))
		// logger.AppLogger().Debugf("devShowFileds, strings.TrimSpace(s1)=[%v]", strings.TrimSpace(s1))
		ret = append(ret, strings.TrimSpace(s2))
		// logger.AppLogger().Debugf("devShowFileds, strings.TrimSpace(s2)=[%v]", strings.TrimSpace(s2))
	}
	return ret
}

func deviceShow(device string) (*NetDevice, error) {
	params := []string{"dev", "show", device}
	logger.AppLogger().Debugf("deviceShow, run cmd: nmcli %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("nmcli", params)
	// fmt.Printf("GetIpAddress show, stdOutput:\n%v\n\nerrOutput:%v\n\nerr:\n\n%v\n\n",
	// 	string(stdOutput), string(errOutput), err)
	if err != nil {
		return nil, fmt.Errorf("failed run deviceShow %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	// logger.AppLogger().Debugf("deviceShow, run cmd: nmcli %v, stdOutput is :%v, errOutput is :%v",
	// 	strings.Join(params, " "), string(stdOutput), string(errOutput))
	logger.AppLogger().Debugf("deviceShow, run cmd: nmcli %v, errOutput is :%v",
		strings.Join(params, " "), string(errOutput))

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)

	deviceInfo := &NetDevice{}
	for _, v := range lines {
		// logger.AppLogger().Debugf("deviceShow, row:%v, line:%v", row, string(v))
		// fields := strings.Fields(v)
		// fmt.Printf("fields: %v\n", fields)
		fields := devShowFileds(v)

		if len(fields) < 1 {
			continue
		}

		k := fields[0]
		k = strings.ReplaceAll(k, ":", "")
		v := fields[1]
		// fmt.Printf("@@  [%v]-->[%v]\n", k, v)

		if k == "GENERAL.DEVICE" {
			deviceInfo.GeneralDevice = v
		}
		if k == "GENERAL.HWADDR" {
			deviceInfo.GeneralHWAddr = v
		}
		if k == "GENERAL.TYPE" {
			deviceInfo.GeneralType = v
		}
		// if k == "GENERAL.STATE" {
		// 	deviceInfo.GENERALSTATE = v
		// }
		if k == "GENERAL.CONNECTION" {
			deviceInfo.GeneralConnection = v
		}
		if k == "IP4.ADDRESS[1]" {
			deviceInfo.Ip4Address = v
		}
		// if k == "IP4.GATEWAY" {
		// 	deviceInfo.IP4GATEWAY = v
		// }
		// if k == "IP6.ADDRESS[1]" {
		// 	deviceInfo.IP6ADDRESS1 = v
		// }
		// if k == "IP6.GATEWAY" {
		// 	deviceInfo.IP6GATEWAY = v
		// }
	}
	return deviceInfo, nil
}
