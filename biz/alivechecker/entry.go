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
 * @Date: 2021-11-10 14:41:01
 * @LastEditors: jeffery
 * @LastEditTime: 2022-04-13 13:37:47
 * @Description:
 */
package alivechecker

import (
	checkerimp "agent/biz/alivechecker/checkerimplement"
	"agent/biz/alivechecker/model"
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/config"
	"fmt"
	"path"
	"strings"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/sys/run"
)

var tickerDockerAliveChecker *time.Ticker
var tickerNetworkChecker *time.Ticker
var checkers []AliveChecker
var tickCnt int64

func Start() {
	StartTestNetwork()

	// docker.SubscribeAsyncContainerStaus(dockerStatusCallback)
	checkers = append(checkers, new(checkerimp.GatewayAliveChecker))
	// 在这里继续添加其他的 AliveChecker 实例

	StartTimerDockerAliveChecker()
}

func ContainerNameGateway() string {
	return checkerimp.ContainerNameGateway
}

func GetContainerStatus(containerName string) bool {
	for _, c := range checkers {
		if c.Name() == containerName {
			return c.Check()
		}
	}
	return false
}

func Stop() {
	// docker.UnsubscribeContainerStaus(dockerStatusCallback)
}

func dockerStatusCallback(status int) {
	logger.CheckLogger().Debugf("dockerStatusCallback, status=%v", status)

	// if docker.ContainersStarted == status {
	// StartTimerDockerAliveChecker()
	// } else {
	// 	StopTimerDockerAliveChecker()
	// }
}

func timerCallbackDockerAliveChecker(ticker *time.Ticker) {
	for range ticker.C {
		for _, c := range checkers {
			if c.Enable() && !c.Check() && c.Restart() {
				logger.CheckLogger().Debugf("timerCallbackDockerAliveChecker,  c.Name():%v, c.Check():%v", c.Name(), c.Check())
			}
		}
	}
}

func StartTimerDockerAliveChecker() {
	StopTimerDockerAliveChecker()

	tickerDockerAliveChecker = time.NewTicker(time.Second * time.Duration(config.Config.AliveChecker.DockerAliveCheckIntervalSec))
	go timerCallbackDockerAliveChecker(tickerDockerAliveChecker)
}

func StopTimerDockerAliveChecker() {
	if tickerDockerAliveChecker != nil {
		tickerDockerAliveChecker.Stop()
		tickerDockerAliveChecker = nil
	}
}

func StartTimerNetworkChecker() {
	StopTimerNetworkChecker()

	tickerNetworkChecker = time.NewTicker(time.Second)
	go timerCallbackNetworkChecker(tickerNetworkChecker)
}

func StopTimerNetworkChecker() {
	if tickerNetworkChecker != nil {
		tickerNetworkChecker.Stop()
		tickerNetworkChecker = nil
	}
}

func timerCallbackNetworkChecker(ticker *time.Ticker) {
	if config.Config.AliveChecker.LogVersionInfoIntervalSec <= 0 {
		config.Config.AliveChecker.LogVersionInfoIntervalSec = 60
	}
	if config.Config.AliveChecker.TestPlatformNetworkIntervalSec <= 0 {
		config.Config.AliveChecker.TestPlatformNetworkIntervalSec = 60
	}

	logger.CheckLogger().Debugf("TestPlatformNetworkIntervalSec=%v", config.Config.AliveChecker.TestPlatformNetworkIntervalSec)

	for range ticker.C {
		if tickCnt%int64(config.Config.AliveChecker.LogVersionInfoIntervalSec) == 0 {
			logger.CheckLogger().Infof("alive checker: [system-agent version:%v]", config.Version)
			logger.CheckLogger().Infof("alive checker: [GetDeviceInfo:%+v]", device.GetDeviceInfo())
			logger.CheckLogger().Infof("alive checker: [GetDeviceInfo.NetworkClient:%+v]", device.GetDeviceInfo().NetworkClient)
			// logger.CheckLogger().Infof("alive checker: [GetAdminPairedInfo:%+v]", clientinfo.GetAdminPairedInfo())
		}

		if tickCnt%int64(config.Config.AliveChecker.TestPlatformNetworkIntervalSec) == 0 {
			TestNetwork()
		}

		tickCnt++
	}
}

func StartTestNetwork() {
	go func() {
		totalTry := 1000
		if config.Config.DebugMode {
			totalTry = 2
		}
		for i := 0; i < totalTry; i++ {
			TestCloudHost()
			time.Sleep(time.Second * 3)
		}

		StartTimerNetworkChecker()
	}()
}

func TestNetwork() {
	result := &model.NetworkTestResult{}

	ok, _ := Ping(config.Config.NetworkCheck.CloudHost.Url)
	result.PingCloudHost = ok
	Ping(config.Config.NetworkCheck.ThirdPartyHost.Url)
	result.PingThirdPartyHost = ok
	Ping(config.Config.NetworkCheck.CloudIpv4.Url)
	result.PingCloudIpv4 = ok
	Curl(config.Config.NetworkCheck.CloudStatusHost.Url)
	result.CurlCloudStatusHost = ok
	CurlHttpHeader(config.Config.NetworkCheck.CloudStatusIpv4.Url)
	result.CurlHttpHeaderCloudStatusIpv4 = ok
	domain := clientinfo.GetAdminDomain()
	if len(domain) > 0 {
		Curl(path.Join(domain, config.Config.NetworkCheck.BoxStatusPath.Url))
	}

	model.Refresh(result)
}

func TestCloudHost() {
	ok, err := Ping(config.Config.NetworkCheck.CloudHost.Url)
	if err != nil {
		logger.CheckLogger().Warnf("failed Ping %v , err:%v", config.Config.NetworkCheck.CloudHost.Url, err)
		model.RefreshPingCloudHost(false)
	} else {
		model.RefreshPingCloudHost(ok)
	}
}

func Curl(host string) (bool, error) {
	params := []string{"--connect-timeout", "20", "-m", "20", host}
	logger.CheckLogger().Debugf("Curl, run cmd: curl %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("curl", params)
	if err != nil {
		return false, fmt.Errorf("failed run Curl %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.CheckLogger().Debugf("Curl return, run cmd: Curl %v, , stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))
	return true, err
}

func CurlHttpHeader(host string) (bool, error) {
	params := []string{"--connect-timeout", "20", "-m", "20", "-I", host}
	logger.CheckLogger().Debugf("Curl, run cmd: curl %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("curl", params)
	if err != nil {
		return false, fmt.Errorf("failed run Curl %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.CheckLogger().Debugf("Curl return, run cmd: Curl %v, , stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))
	return true, err
}

func Ping(host string) (bool, error) {
	params := []string{"-c", "3", host}
	logger.CheckLogger().Debugf("Ping, running cmd: ping %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("ping", params)
	if err != nil {
		return false, fmt.Errorf("failed run Ping %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.CheckLogger().Debugf("Ping return, run cmd: ping %v, , stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))
	return true, err
}
