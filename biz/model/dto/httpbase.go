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
 * @LastEditors: wenchao
 * @LastEditTime: 2021-11-30 15:23:38
 * @Description:
 */
package dto

import (
	"fmt"
	"net/http"
)

const (
	AgentCodeOk = http.StatusOK

	AgentCodeOkStr          = "AG-200"
	AgentCodeInitialSuccStr = "AG-201"

	AgentCodeBadReqStr            = "AG-400"
	AgentCodeParamErr             = "AG-401" //参数错误
	AgentCodeResBusyErr           = "AG-402" //资源忙，使用中
	AgentCodeUnsupportedFunction  = "AG-403" //当前环境不支持此项功能
	AgentCodeAlreadyPairedStr     = "AG-460" // 已经绑定
	AgentCodePwdErrorOverLimitStr = "AG-461"
	AgentCodeUnpairedBeforeStr    = "AG-462"
	AgentCodeAdminPwdError        = "AG-463"
	AgentCodeRepeatedRequest      = "AG-464"
	AgentCodeTryOutCodeError      = "AG-465" // 试用码错误
	AgentCodeTryOutCodeExpired    = "AG-466" // 试用码过期
	AgentCodeTryOutCodeHasUsed    = "AG-467" // 试用码已经使用过了
	AgentCodeTryOutCodeDisabled   = "AG-468" // 试用码禁用
	AgentCodeDockerPulling        = "AG-469" // 容器下载中
	AgentCodeDockerStarting       = "AG-470" // 容器启动中
	AgentCodeDockerStarted        = "AG-471" // 容器已经启动

	AgentCodeServerErrorStr              = "AG-500"
	AgentCodeCallServiceFailedStr        = "AG-560"
	AgentCodeConnectWifiFailedStr        = "AG-561"
	AgentCodeRevokeDirectFailedStr       = "AG-562"
	AgentCodeCallDiskUsageFailedStr      = "AG-563"
	AgentCodeMissingMainStorageFailedStr = "AG-590"

	AgentCodeSwitchDomainErr         = "AG-570" //域名错误
	AgentCodeSwitchToNewSSPErr       = "AG-571" //临时指向到新的空间平台错误
	AgentCodeSwitchNetworkTestErr    = "AG-572" //网络测试失败
	AgentCodeSwitchRecallGatewayErr  = "AG-573" //回写网关失败
	AgentCodeSwitchDoingErr          = "AG-574" //切换任务正在执行
	AgentCodeSwitchGetAccountErr     = "AG-575" //从网关获取账号失败
	AgentCodeSwitchImigrateErr       = "AG-576" //新空间平台迁入失败
	AgentCodePrivateSSPRegBoxErr     = "AG-577" //私有空间平台注册盒子错误
	AgentCodeCopyDataToStatusInfoErr = "AG-578" //将迁入请求数据转移到状态信息失败

	AgentCodeSwitchTaskNotFoundErr = "AG-580" //切换任务未找到
	AgentCodeConnectErr            = "AG-581" //连接错误
	//AlreadyLatestVersion           = "AG-591"

	GatewayCodeOkStr = "GW-200"
	AccountCodeOkStr = "ACC-200"
)

type BaseRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message, omitempty"`
	Results interface{} `json:"results, omitempty"`
}

// c.JSON(http.StatusOK, gin.H{"code": AgentCodeOk, "message": "OK"})
func NewBaseRsp(code int, message string, results interface{}) BaseRsp {
	return BaseRsp{Code: code, Message: message, Results: results}
}

// 后来开发规范要求 code 是 string 了
type BaseRspStr struct {
	Code      string      `json:"code"`
	RequestId string      `json:"requestId, omitempty"`
	Message   string      `json:"message, omitempty"`
	Results   interface{} `json:"results, omitempty"`
}

// 客户端的附加参数，如鉴权字段。由于蓝牙没有类似于 HTTP HEADER 字段，所以暂时把这类信息放在请求体中。
type ReqExtraInfo struct {
	AgentToken string `json:"agentToken"` // 客户端通过局域网/蓝牙调用 agent 接口时的鉴权参数 access_token，其中包含公开字段: boxUuid、clientUuid 和 tokenType。
}

// 局域网、蓝牙调用的请求结构
type LanInvokeReq struct {
	Body string `json:"body" form:"body"` // 请求体
}

// c.JSON(http.StatusOK, gin.H{"code": AgentCodeOkStr, "message": "OK"})
func NewBaseRspStr(code string, message string, results interface{}) *BaseRspStr {
	return &BaseRspStr{Code: code, Message: message, Results: results}
}

// c.JSON(http.StatusOK, gin.H{"code": AgentCodeOkStr, "message": "OK"})
func NewBaseResponse(code, message, requestId string, results interface{}) *BaseRspStr {
	return &BaseRspStr{Code: code, Message: message, RequestId: requestId, Results: results}
}

type UpdateError error

var (
	AlreadyLatestVersion UpdateError = fmt.Errorf("already the latest version")
)
