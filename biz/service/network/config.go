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
	"agent/biz/alivechecker/model"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/dto/network"
	"agent/biz/service/base"
	"agent/config"
	util_network "agent/utils/network"
	rpi_network "agent/utils/rpi/network"
	"fmt"

	"agent/utils/logger"
)

///////////////////////////////////////////////////////////////////////////////
// 用修改配置文件的方式来进行网络配置
// func NewPostNetworkConfigService() *PostNetworkConfigByCfgFileService {
// 	svc := new(PostNetworkConfigByCfgFileService)
// 	return svc
// }

// func NewGetNetworkConfigService() *GetNetworkConfigByCfgFileService {
// 	svc := new(GetNetworkConfigByCfgFileService)
// 	return svc
// }

// 用 执行 nmcli 程序的方式来进行网络配置
func NewPostNetworkConfigService() *PostNetworkConfigService {
	svc := new(PostNetworkConfigService)
	return svc
}

func NewGetNetworkConfigService() *GetNetworkConfigService {
	svc := new(GetNetworkConfigService)
	return svc
}

///////////////////////////////////////////////////////////////////////////////

type GetNetworkConfigService struct {
	base.BaseService
}

func (svc *GetNetworkConfigService) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("GetNetworkConfigService")

	var rsp network.NetworkStatusRsp

	// 读取网络是否连通
	rsp.InternetAccess = model.Get().PingCloudHost
	rsp.LanAccessPort = config.Config.GateWay.LanPort
	rsp.NetworkAdapters = []*network.NetworkAdapter{}

	if device_ability.GetAbilityModel().RunInDocker {
		svc.Rsp = rsp
		return svc.BaseService.Process()
	}

	networkDevInfos, err := util_network.GetNetworkDevStatus()
	if err != nil {
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err.Error()}
	}

	for i, devStatus := range networkDevInfos {
		logger.AppLogger().Debugf("GetNetworkConfigService, devStatus(%v/%v): %+v", i, len(networkDevInfos), devStatus)
		if !devStatus.IsEthernetAndWifi() {
			continue
		}
		if !devStatus.IsConnected() {
			continue
		}

		conInfo, err := util_network.GetNetworkConInfo(devStatus.CONUUID)
		if err != nil {
			logger.AppLogger().Warnf("failed GetNetworkConInfo, devStatus.CONUUID:%v, err:%v", devStatus.CONUUID, err)
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
				Message: err.Error()}
		}
		logger.AppLogger().Debugf("GetNetworkConfigService, conInfo: %+v", conInfo)

		devInfo, err := util_network.GetNetworkDevInfo(devStatus.DEVICE)
		if err != nil {
			logger.AppLogger().Warnf("failed GetNetworkDevInfo, devStatus.DEVICE:%v, err:%v", devStatus.DEVICE, err)
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
				Message: err.Error()}
		}

		networkAdapter := network.NetworkAdapter{}
		networkAdapter.AdapterName = devStatus.DEVICE
		networkAdapter.Wired = devStatus.IsWire()
		if devStatus.IsWireless() {
			networkAdapter.WIFIAddress = conInfo.SeenBssids
			networkAdapter.WIFIName = conInfo.GENERALNAME
		}
		networkAdapter.Connected = devStatus.IsConnected()
		networkAdapter.MACAddress = devInfo.GENERALHWADDR

		networkAdapter.Ipv4UseDhcp = conInfo.UseDhcp()
		networkAdapter.Ipv4 = conInfo.Ipv4Address()
		networkAdapter.SubNetMask = conInfo.Ipv4NetMask()
		networkAdapter.DefaultGateway = conInfo.IP4GATEWAY
		networkAdapter.Ipv4DNS1 = conInfo.IP4DNS1
		networkAdapter.Ipv4DNS2 = conInfo.IP4DNS2
		networkAdapter.Ipv6UseDhcp = conInfo.UseDhcpIpv6()
		networkAdapter.Ipv6 = conInfo.IP6ADDRESS
		networkAdapter.Ipv6DefaultGateway = conInfo.IP6GATEWAY

		rsp.NetworkAdapters = append(rsp.NetworkAdapters, &networkAdapter)

		if len(conInfo.IP4DNS1) > 0 {
			rsp.DNS1 = conInfo.IP4DNS1
		}
		if len(conInfo.IP4DNS2) > 0 {
			rsp.DNS2 = conInfo.IP4DNS2
		}
	}

	// 获取 dns
	dns1, dns2, err := util_network.GetSystemdDns()
	if err != nil {
		logger.AppLogger().Warnf("failed GetSystemdDns, err:%v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
	}
	rsp.DNS1 = dns1
	if len(dns2) > 0 {
		rsp.DNS2 = dns2[len(dns2)-1]
	}

	svc.Rsp = rsp
	return svc.BaseService.Process()
}

type PostNetworkConfigService struct {
	base.BaseService
}

func (svc *PostNetworkConfigService) Process() dto.BaseRspStr {
	req := svc.Req.(*network.NetworkConfigReq)
	// logger.AppLogger().Debugf("PostNetworkConfigService, req:%+v", req)
	abilityModel := device_ability.GetAbilityModel()
	if !abilityModel.InnerDiskSupport {
		err := fmt.Errorf("unsupported function")
		return dto.BaseRspStr{Code: dto.AgentCodeUnsupportedFunction,
			Message: err.Error()}
	}

	// 尝试连接传入的 wifi
	for _, adapter := range req.NetworkAdapters {
		logger.AppLogger().Debugf("PostNetworkConfigService, adapter:%+v", adapter)

		if adapter.Wired {
			continue
		}

		// 查看是否已经连接上了
		alreadyConnectedRequestedWifi := false
		devStatuss, err := rpi_network.GetIpAddress()
		if err != nil {
			logger.AppLogger().Warnf("failed GetIpAddress, err:%v", err)
		} else {
			for _, devStatus := range devStatuss {
				if devStatus.GeneralConnection == adapter.WIFIName {
					alreadyConnectedRequestedWifi = true
					break
				}
			}
		}
		if alreadyConnectedRequestedWifi {
			logger.AppLogger().Debugf("already connected before.")
			break
		}

		// 尝试去连接
		logger.AppLogger().Debugf("try to connect wifi, adapter.WIFIAddress:[%v], adapter.WIFIName:[%v]",
			adapter.WIFIAddress, adapter.WIFIName)
		succ := false
		if len(adapter.WIFIAddress) > 0 {
			succ, err = rpi_network.ConnectWifi(adapter.WIFIAddress, adapter.WIFIPassword)
			logger.AppLogger().Debugf("connect wifi using WIFIAddress result: %v, err:%v", succ, err)
		}
		if !succ && len(adapter.WIFIName) > 0 {
			succ, err = rpi_network.ConnectWifi(adapter.WIFIName, adapter.WIFIPassword)
			logger.AppLogger().Debugf("connect wifi using WIFIName result: %v, err:%v", succ, err)
		}
		if succ {
			logger.AppLogger().Debugf("connect wifi succ, adapter.WIFIName:%+v, adapter.WIFIAddress:%+v", adapter.WIFIName, adapter.WIFIAddress)
			break
		} else {
			err1 := fmt.Errorf("PostNetworkConfigService, connect wifi err:%v", err)
			logger.AppLogger().Warnf(err1.Error())
			if err != nil {
				return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
					Message: err1.Error()}
			} else {
				return dto.BaseRspStr{Code: dto.AgentCodeConnectWifiFailedStr,
					Message: err1.Error()}
			}
		}
	}

	networkDevInfos, err := util_network.GetNetworkDevStatus()
	if err != nil {
		logger.AppLogger().Warnf("failed GetNetworkDevStatus, err:%v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err.Error()}
	}

	for _, adapter := range req.NetworkAdapters {
		networkDevInfo := util_network.GetNetworkDevStatusByDevice(adapter.AdapterName, networkDevInfos)
		logger.AppLogger().Debugf("start configing, adapter: %+v, networkDevInfo=%+v", adapter, networkDevInfo)
		if networkDevInfo == nil {
			continue
		}

		if adapter.Wired { // 有线
			if adapter.Ipv4UseDhcp { // 自动
				if err := util_network.SetWireIpAuto(adapter.AdapterName); err != nil {
					logger.AppLogger().Warnf("failed SetWireIpAuto, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
					return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
				}

				if err = util_network.SetNetworkDeviceDown(adapter.AdapterName); err != nil {
					logger.AppLogger().Warnf("failed SetNetworkDeviceDown, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
					return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
				}
			} else { // 手动
				if err := util_network.SetWireIpManual(adapter.AdapterName, adapter.Ipv4, adapter.SubNetMask, adapter.DefaultGateway); err != nil {
					logger.AppLogger().Warnf("failed SetWireIpManual, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
					return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
				}

				if err := util_network.SetWireDnsManual(adapter.AdapterName, req.DNS1, req.DNS2); err != nil {
					logger.AppLogger().Warnf("failed SetWireDnsManual, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
					return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
				}
			}

			if err = util_network.SetNetworkDeviceUp(adapter.AdapterName); err != nil {
				logger.AppLogger().Warnf("failed SetNetworkDeviceUp, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
				return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
			}

		} else { // 无线
			if adapter.Ipv4UseDhcp { // 自动
				if err := util_network.SetWirelessIpAuto(networkDevInfo.CONUUID); err != nil {
					logger.AppLogger().Warnf("failed SetWirelessIpAuto, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
					return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
				}

				if err = util_network.SetNetworkDeviceDown(networkDevInfo.CONUUID); err != nil {
					logger.AppLogger().Warnf("failed SetNetworkDeviceDown, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
					return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
				}
			} else { // 手动
				if err := util_network.SetWirelessIpManual(networkDevInfo.CONUUID, adapter.Ipv4, adapter.SubNetMask, adapter.DefaultGateway); err != nil {
					logger.AppLogger().Warnf("failed SetWirelessIpManual, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
					return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
				}

				if err := util_network.SetWirelessDnsManual(networkDevInfo.CONUUID, req.DNS1, req.DNS2); err != nil {
					logger.AppLogger().Warnf("failed SetWirelessDnsManual, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
					return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
				}
			}

			if err = util_network.SetNetworkDeviceUp(networkDevInfo.CONUUID); err != nil {
				logger.AppLogger().Warnf("failed SetNetworkDeviceUp, adapter.AdapterName:%v, err:%v", adapter.AdapterName, err)
				return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
			}
		}
	}

	// 修改 dns
	if len(req.DNS1) > 0 && len(req.DNS2) > 0 {
		if err = util_network.SetSystemdDnsManual(req.DNS1, req.DNS2); err != nil {
			logger.AppLogger().Warnf("failed SetSystemdDnsManual, err:%v", err)
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
		}
	} else {
		if err = util_network.SetSystemdDnsDefault(); err != nil {
			logger.AppLogger().Warnf("failed SetSystemdDnsDefault, err:%v", err)
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
		}
	}

	return svc.BaseService.Process()
}
