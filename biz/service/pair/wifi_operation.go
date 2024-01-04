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

package pair

import (
	dtopair "agent/biz/model/dto/pair"
	"agent/config"
	"agent/utils/rpi/network"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

var lastWifiList []*dtopair.WifiListRsp

func sortWifiList(ret []*network.ListWifiInfo) []*network.ListWifiInfo {

	// slice 按照信号强度升序排序.
	sort.SliceStable(ret, func(i, j int) bool {

		iSig, err := strconv.Atoi(ret[i].SIGNAL)
		if err != nil {
			logger.AppLogger().Warnf("failed strconv.Atoi, ret[i].SIGNAL=%v, err:%v", ret[i].SIGNAL, err)
			return false
		}
		jSig, err := strconv.Atoi(ret[j].SIGNAL)
		if err != nil {
			logger.AppLogger().Warnf("failed strconv.Atoi, ret[j].SIGNAL=%v, err:%v", ret[j].SIGNAL, err)
			return false
		}

		return iSig > jSig
	})
	logger.AppLogger().Debugf("sortWifiList, SliceStable len(ret):%+v", len(ret))

	return ret
}

// https://pm.eulix.xyz/bug-view-164.html
func removeDuplicatedWifi(wifi []*network.ListWifiInfo) []*network.ListWifiInfo {
	logger.AppLogger().Debugf("removeDuplicatedWifi, len(wifi):%+v", len(wifi))

	// 把重复的 wifiname 去掉, 同一个名称的wifi 只保留信号强度最高的.
	mapName := make(map[string]*network.ListWifiInfo)
	for _, info := range wifi {
		if val, ok := mapName[info.SSID]; ok {

			infoSig, err := strconv.Atoi(info.SIGNAL)
			if err != nil {
				logger.AppLogger().Warnf("failed strconv.Atoi, info.SIGNAL=%v, err:%v", info.SIGNAL, err)
				continue
			}
			valSig, err := strconv.Atoi(val.SIGNAL)
			if err != nil {
				logger.AppLogger().Warnf("failed strconv.Atoi, val.SIGNAL=%v, err:%v", val.SIGNAL, err)
				continue
			}

			if infoSig > valSig {
				mapName[info.SSID] = info
			}
		} else {
			mapName[info.SSID] = info
		}
	}
	logger.AppLogger().Debugf("removeDuplicatedWifi, len(mapName):%+v", len(mapName))

	// map 放入 slice
	ret := []*network.ListWifiInfo{}
	for _, v := range mapName {
		ret = append(ret, v)
	}
	logger.AppLogger().Debugf("removeDuplicatedWifi, len(ret):%+v", len(ret))

	return sortWifiList(ret)
}

func GetWifiList() []*dtopair.WifiListRsp {
	logger.AppLogger().Debugf("GetWifiList")
	// var err error
	lst, err := network.ListWifi()
	if err != nil {
		logger.AppLogger().Warnf("failed ListWifi, err:%v", err)
		return []*dtopair.WifiListRsp{}
	} else {
		// for i, v := range lst {
		// logger.AppLogger().Debugf("getWList, lst[%v] = %+v", i, v)
		// }
	}
	// logger.AppLogger().Debugf("getWList, lst:%+v", lst)

	// 按照信号强度来排序，并且去掉重新的wifi 名称.
	// https://pm.eulix.xyz/bug-view-164.html
	wifi := removeDuplicatedWifi(lst)
	// logger.AppLogger().Debugf("getWList, wifi:%+v", wifi)

	// for i, v := range wifi {
	// 	logger.AppLogger().Debugf("getWList, after removeDuplicatedWifi, wifi[%v] = %+v", i, v)
	// }

	wifiList := make([]*dtopair.WifiListRsp, 0)
	for _, v := range wifi {
		sig, err := strconv.Atoi(v.SIGNAL)
		if err != nil {
			logger.AppLogger().Warnf("failed strconv.Atoi, v.SIGNAL=%v, err:%v", v.SIGNAL, err)
		}
		wifiList = append(wifiList, &dtopair.WifiListRsp{Name: v.SSID, Addr: v.BSSID, Signal: int8(sig)})
	}

	lastWifiList = wifiList

	return wifiList
	// logger.AppLogger().Debugf("getWList, wifiList:%+v", wifiList)
}

func GetLocalIpBySSID(ssid string) (string, string) {
	logger.AppLogger().Debugf("GetLocalIpBySSID, ssid=%v, lastWifiList=%+v", ssid, lastWifiList)

	name := ""
	ip := ""
	if lastWifiList == nil {
		logger.AppLogger().Debugf("GetLocalIpBySSID, ssid=%v, lastWifiList is nil", ssid)
		return name, ip
	}

	wifiName := ""
	for _, v := range lastWifiList {
		if v.Addr == ssid {
			wifiName = v.Name
			break
		}
	}
	if len(wifiName) < 1 {
		logger.AppLogger().Debugf("GetLocalIpBySSID, ssid=%v, wifiName=%+v, return", ssid, wifiName)
		return name, ip
	}
	name = wifiName
	logger.AppLogger().Debugf("GetLocalIpBySSID, ssid=%v, wifiName=%+v", ssid, wifiName)

	deviceInfos, err := network.GetIpAddress()
	if err != nil {
		logger.AppLogger().Warnf("GetLocalIpBySSID, failed GetIpAddress, err:%v", err)
		return name, ip
	}

	// 去掉ip地址最后的 "/"
	for i, v := range deviceInfos {
		logger.AppLogger().Debugf("GetLocalIpBySSID, deviceInfos[%d]: %+v", i, v)

		arr := strings.Split(v.Ip4Address, "/")
		if len(arr) > 1 {
			logger.AppLogger().Debugf("GetLocalIpBySSID, v.GENERALCONNECTION=%v, wifiName=%+v",
				v.GeneralConnection, wifiName)
			if v.GeneralConnection == wifiName {
				ip = arr[0]
				logger.AppLogger().Debugf("GetLocalIpBySSID, arr[0]=%v", arr[0])
				break
			}
		}
	}
	return name, ip
}

func GetLocalIp() []string {
	deviceInfos, err := network.GetIpAddress()
	if err != nil {
		logger.AppLogger().Warnf("GetLocalIp, failed GetIpAddress, err:%v", err)
		return []string{}
	}
	for i, v := range deviceInfos {
		logger.AppLogger().Debugf("GetLocalIp, deviceInfos[%d]: %+v", i, v)
	}

	// 去掉ip地址最后的 "/"
	rt := make([]string, 0)
	for _, v := range deviceInfos {
		arr := strings.Split(v.Ip4Address, "/")
		if len(arr) > 1 {
			rt = append(rt, arr[0])
		}
	}
	return rt
}

func IsNetworkConnected() (bool, error) {

	ok, err := network.Ping(config.Config.Box.PingHost)
	if err != nil {
		logger.AppLogger().Warnf("failed Ping, err:%v", err)
		return false, err
	}
	logger.AppLogger().Debugf("Ping ok:%v", ok)
	return ok, nil
}

func ConnectToWifi(BSSID, PWD string) error {
	logger.AppLogger().Debugf("connectToWifi, BSSID:%+v", BSSID)

	succ, err := network.ConnectWifi(BSSID, PWD)
	if err != nil {
		logger.AppLogger().Warnf("failed ConnectWifi, err:%v", err)
		return err
	}
	logger.AppLogger().Debugf("ConnectWifi result:%v\n", succ)
	if !succ {
		return fmt.Errorf("ConnectWifi return false")
	}

	return nil
}

func GetConnectedNetwork() []*dtopair.Network {
	logger.AppLogger().Debugf("GetConnectedNetwork")

	deviceInfos, err := network.GetIpAddress()
	if err != nil {
		logger.AppLogger().Warnf("GetConnectedNetwork, failed GetIpAddress, err:%v", err)

		if fileutil.IsFileExist(config.Config.Box.HostIpFile) {
			logger.AppLogger().Debugf("GetConnectedNetwork, IsFileExist, HostIpFile:%v", config.Config.Box.HostIpFile)
			ipaddr, err := fileutil.ReadFromFile(config.Config.Box.HostIpFile)
			if err != nil {
				return []*dtopair.Network{}
			}
			logger.AppLogger().Debugf("GetConnectedNetwork, IsFileExist, ipaddr:%v", string(ipaddr))

			arr := strings.Split(string(ipaddr), ":")
			if len(arr) >= 2 {
				logger.AppLogger().Debugf("GetConnectedNetwork, IsFileExist, Ip:%v", string(arr[0]))
				ret := make([]*dtopair.Network, 0)
				ret = append(ret, &dtopair.Network{Ip: arr[0], Wire: true, WifiName: "",
					Port:    config.Config.GateWay.LanPort,
					TlsPort: config.Config.GateWay.TlsLanPort})
				return ret
			}

			return []*dtopair.Network{}
		}
		logger.AppLogger().Debugf("GetConnectedNetwork, File NOT Exist, HostIpFile:%v", config.Config.Box.HostIpFile)

		return []*dtopair.Network{}
	}
	for i, v := range deviceInfos {
		logger.AppLogger().Debugf("GetConnectedNetwork, deviceInfos[%d]: %+v", i, v)
	}

	rt := make([]*dtopair.Network, 0)
	for _, v := range deviceInfos {

		// 去掉ip地址最后的 "/"
		ip := v.Ip4Address
		arr := strings.Split(v.Ip4Address, "/")
		if len(arr) > 1 {
			ip = arr[0]
		}

		n := &dtopair.Network{Ip: ip, Wire: true,
			WifiName: v.GeneralDevice,
			Port:     config.Config.GateWay.LanPort,
			TlsPort:  config.Config.GateWay.TlsLanPort}

		if v.GeneralType == "wifi" {
			n = &dtopair.Network{Ip: ip, Wire: false,
				WifiName: v.GeneralConnection,
				Port:     config.Config.GateWay.LanPort,
				TlsPort:  config.Config.GateWay.TlsLanPort}
		}
		rt = append(rt, n)

	}

	return rt
}
