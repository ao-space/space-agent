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

package switchplatform

import (
	"agent/biz/model/device"
	"agent/biz/model/dto"
	modelsp "agent/biz/model/switch-platform"
	"agent/biz/service/encwrapper"
	"agent/utils"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"net"
	"strings"

	"agent/utils/logger"
)

// 切换盒子对接的空间平台
func ServiceSwitchPlatform(req *modelsp.SwitchPlatformReq) (dto.BaseRspStr, error) {
	logger.AppLogger().Debugf("ServiceSwitchPlatform, req:%+v", req)
	logger.AccessLogger().Debugf("[ServiceSwitchPlatform], req:%+v", req)

	err := encwrapper.Check()
	if err != nil {
		logger.AppLogger().Warnf("check, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}

	rt, err := encwrapper.Dec(req.TransId, req.NewDomain)
	if err != nil {
		logger.AppLogger().Warnf("dec, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}
	TransId := rt[0]
	NewDomain := rt[1]

	logger.AppLogger().Debugf("transId:%+v, domain:%+v",
		TransId, NewDomain)
	rsp, err := doSwitch(TransId, NewDomain)

	logger.AppLogger().Debugf("ServiceSwitchPlatform, transId:%v,  rsp:%+v, err:%v", TransId, rsp, err)
	return rsp, err
}

func init() {
	// go func() {
	// 	time.Sleep(30 * time.Second)
	// 	rsp, err := doSwitch("test_migrate", "dev-services.eulix.xyz")
	// 	//rsp, err := doSwitch("test_migrate", "qa.eulix.xyz")

	// 	logger.AppLogger().Debugf("transId=%v, *** rsp=%v , err:=%v",
	// 		"test_migrate", rsp, err)
	// }()
}

func doSwitch(transId string, domain string) (dto.BaseRspStr, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	logger.AppLogger().Debugf("transId=%v, domain=%v ",
		transId, domain)

	host, _, _ := utils.ParseUrl(device.GetApiBaseUrl())

	//检测参数
	if len(domain) == 0 {
		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeParamErr
		basersp.Message = fmt.Sprintf("domain[%v] is invalid", domain)
		return basersp, errors.New(basersp.Message)
	}

	if domain == host {
		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeResBusyErr
		basersp.Message = fmt.Sprintf("domain[%v] is using", domain)
		return basersp, errors.New(basersp.Message)
	}

	if conn, err := net.Dial("tcp", domain+":443"); err != nil {
		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeConnectErr
		basersp.Message = fmt.Sprintf("failed to connect to %v:443, err:%v", domain, err.Error())
		return basersp, errors.New(basersp.Message)
	} else {
		conn.Close()
	}

	if _, err := createStatus(); err != nil {

		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeSwitchDoingErr
		basersp.Message = "task is doing"
		return basersp, errors.New(basersp.Message)
	}

	//---------- 此功能预计在第二阶段真正实现
	//获取账号
	accountInfos, err := getAccountInfos(transId)
	if err != nil {
		logger.AppLogger().Debugf("transId=%v, accountInfos err:%+v", transId, err)

		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeSwitchGetAccountErr
		basersp.Message = "failed to get account info from gateway"

		UpdateStatus(StatusAbort, basersp.Message)
		return basersp, err
	}
	si.OldAccount = accountInfos
	si.OldApiBaseUrl = device.GetApiBaseUrl()
	si.NewApiBaseUrl = "https://" + domain
	si.Domain = domain
	si.TransId = transId

	logger.AppLogger().Debugf("transId=%v, accountInfos:%+v, ",
		transId, accountInfos)

	//执行迁入
	imigrateRsp, err := imigrate()
	logger.AppLogger().Debugf("imigrateRsp:%+v, ", imigrateRsp)
	if err != nil {
		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeSwitchImigrateErr
		basersp.Message = "failed to imigrate to new SSP."
		UpdateStatus(StatusAbort, basersp.Message)
		return basersp, err
	}

	err = copier.Copy(&si.ImigrateResult, imigrateRsp)
	if err != nil {
		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeCopyDataToStatusInfoErr
		basersp.Message = "failed to transport data to StatusInfo."
		UpdateStatus(StatusAbort, basersp.Message)
		return basersp, err
	}

	logger.AppLogger().Debugf("transId=%v, imigrateRsp:%+v, ",
		transId, imigrateRsp)

	// network-client 对接新的空间平台
	if err := networkSwitchV2(true); err != nil {
		networkSwitchV2(false)

		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeSwitchToNewSSPErr
		basersp.Message = "failed to network to SSP."

		UpdateStatus(StatusAbort, basersp.Message)
		return basersp, err
	}

	logger.AppLogger().Debugf("transId=%v, networkSwitch succ. ",
		transId)

	//进行测试
	if err := networkDetect(transId, domain, &si.ImigrateResult); err != nil {
		logger.AppLogger().Debugf("transId=%v, testNetwork:%+v, ",
			transId, err)
		networkSwitchV2(false)

		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeSwitchNetworkTestErr
		basersp.Message = "failed to test new network."
		UpdateStatus(StatusAbort, basersp.Message)
		return basersp, err
	}
	logger.AppLogger().Debugf("transId=%v, test network succ ", transId)

	UpdateStatus(StatusStart, "StatusStart")

	if err = doStatusFlow(false); err != nil {

		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeSwitchRecallGatewayErr
		basersp.Message = "failed to test new network."

		UpdateStatus(StatusAbort, fmt.Sprintf("failed to doStatusFlow.err:%v", err))
		return basersp, err
	} else {
		var results modelsp.SwitchPlatformResp
		results.UserDomain, _ = si.ImigrateResult.GetAdminDomain()
		results.TransId = si.TransId
		logger.AppLogger().Debugf("transId=%v, results: %+v", transId, results)

		return encwrapper.Enc(results)
	}

}

func doStatusFlow(reboot bool) error {

	if reboot && si.Status != StatusUpdateBoxInfo {
		// network-client，gateway 对接新的空间平台
		networkSwitchV2(true)
	}

	if si.Status == StatusStart {

		//网关更新域名
		if err := updateAccount(); err != nil {
			// 恢复对接老的空间平台
			networkSwitchV2(false)
			UpdateStatus(StatusAbort, "failed to updateAccount")
			logger.AppLogger().Debugf("transId=%v, updateAccount err:%+v, ",
				si.TransId, err)
			return err
		}
		UpdateStatus(StatusUpdateGateway, "StatusUpdateGateway")
	}

	if si.Status == StatusUpdateGateway {
		logger.AppLogger().Debugf("transId=%v, updateAccount succ ", si.TransId)
		// 切换环境
		device.SetApiBaseUrl("https://" + si.Domain)
		device.SetNetworkClient(&si.ImigrateResult.NetworkClient)

		UpdateStatus(StatusUpdateBoxInfo, "StatusUpdateBoxInfo")
	}

	// 异步执行迁出
	go emigrate()
	return nil
}
