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

package model

import (
	"agent/utils/logger"
)

type NetworkTestResult struct {
	PingCloudHost                 bool `json:"pingCloudHost"` // Ping CloudHost 成功
	PingThirdPartyHost            bool `json:"pingThirdPartyHost"`
	PingCloudIpv4                 bool `json:"pingCloudIpv4"`
	CurlCloudStatusHost           bool `json:"curlCloudStatusHost"`
	CurlHttpHeaderCloudStatusIpv4 bool `json:"curlHttpHeaderCloudStatusIpv4"`
}

var networkTestResult *NetworkTestResult

func init() {
	networkTestResult = &NetworkTestResult{}
	networkTestResult.PingCloudHost = false

	networkTestResult.PingThirdPartyHost = false
	networkTestResult.PingCloudIpv4 = false
	networkTestResult.CurlCloudStatusHost = false
	networkTestResult.CurlHttpHeaderCloudStatusIpv4 = false
}

func Refresh(result *NetworkTestResult) {
	networkTestResult.PingCloudHost = result.PingCloudHost

	networkTestResult.PingThirdPartyHost = result.PingThirdPartyHost
	networkTestResult.PingCloudIpv4 = result.PingCloudIpv4
	networkTestResult.CurlCloudStatusHost = result.CurlCloudStatusHost
	networkTestResult.CurlHttpHeaderCloudStatusIpv4 = result.CurlHttpHeaderCloudStatusIpv4
}

func RefreshPingCloudHost(pingCloudHost bool) {
	networkTestResult.PingCloudHost = pingCloudHost

	logger.CheckLogger().Debugf("RefreshPingCloudHost, networkTestResult.PingCloudHost:%v", networkTestResult.PingCloudHost)
}

func Get() *NetworkTestResult {
	logger.CheckLogger().Debugf("Get, NetworkTestResult.PingCloudHost:%v", networkTestResult.PingCloudHost)
	return networkTestResult
}
