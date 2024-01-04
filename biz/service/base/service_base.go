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

package base

import (
	"agent/biz/model/dto"
	"agent/biz/service/encwrapper"
	"agent/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"agent/utils/jwt"
	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/gin-gonic/gin"
)

type IService interface {
	Process() dto.BaseRspStr
}

const (
	CalledType_Lan       = 1 // 局域网
	CalledType_Bluetooth = 2 // 蓝牙
	CalledType_Gateway   = 3 // 网关
)

type BaseService struct {
	// 以下是调用者设置
	RequestId                   string      `json:"requestId"`
	Header                      http.Header `json:"header"`                      // http 请求头
	EncryptedLanRequestBodyData []byte      `json:"encryptedLanRequestBodyData"` // 局域网/蓝牙加密的请求数据
	CalledType                  int         `json:"calledType"`                  // 1 局域网, 2 蓝牙, 3 网关
	ginContext                  *gin.Context

	// 以下是内部使用
	Req      interface{} `json:"req"`
	Rsp      interface{} `json:"rsp"`
	RspBytes []byte      `json:"rspBytes"`
}

func init() {
	jwt.SetSignHost(config.Config.Box.SecurityChipAgentHttpLocalAddr)
}

func (svc *BaseService) InitLanService(RequestId string, header http.Header, ginContext *gin.Context) *BaseService {
	// logger.AppLogger().Debugf("InitLanService, ginContext:%+v", ginContext)
	if len(header.Get("RequestId")) > 0 && len(RequestId) < 1 {
		RequestId = header.Get("RequestId")
	}
	logger.AccessLogger().Debugf("[LanService] RequestId:%+v, Request:%+v, header:%+v", RequestId, ginContext.Request, header)
	svc.RequestId = RequestId
	svc.Header = header
	svc.ginContext = ginContext
	svc.CalledType = CalledType_Lan
	return svc
}

func (svc *BaseService) InitBluetoothService(RequestId string, cmd int, encryptedLanRequestBodyData []byte) *BaseService {
	// logger.AppLogger().Debugf("InitBluetoothService")
	logger.AccessLogger().Debugf("[BluetoothService] RequestId:%+v, cmd:%+v, len(encryptedLanRequestBodyData):%+v", RequestId, cmd, len(encryptedLanRequestBodyData))
	svc.RequestId = RequestId
	svc.EncryptedLanRequestBodyData = encryptedLanRequestBodyData
	svc.CalledType = CalledType_Bluetooth
	return svc
}

func (svc *BaseService) InitGatewayService(RequestId string, header http.Header, ginContext *gin.Context) *BaseService {
	// logger.AppLogger().Debugf("InitGatewayService, header:%+v", header)
	logger.AccessLogger().Debugf("[GatewayService] RequestId:%+v, Request:%+v, header:%+v", RequestId, ginContext.Request, header)
	svc.RequestId = RequestId
	svc.Header = header
	svc.ginContext = ginContext
	svc.CalledType = CalledType_Gateway
	return svc
}

// 进入函数。
// reqObj 是客户端请求对象, 无参接口 reqObj 传 nil.
// 如果是局域网/蓝牙则传入LanInvokeReq结构中 body 对应的对象, 不需要调用者去解密和反序列化。如果是网关则传入反序列化好的请求对象。
func (svc *BaseService) Enter(iService IService, reqObj interface{}) dto.BaseRspStr {
	// logger.AppLogger().Debugf("BaseService Enter, reqObj:%+v", reqObj)

	if (svc.CalledType == CalledType_Lan || svc.CalledType == CalledType_Bluetooth) && reqObj != nil { // 局域网或蓝牙调用

		if svc.CalledType == CalledType_Lan {
			logger.AppLogger().Debugf("BaseService Enter, CalledType_Lan")

			var lanInvokeReq dto.LanInvokeReq
			if err := svc.ginContext.ShouldBind(&lanInvokeReq); err != nil {
				err1 := fmt.Errorf("failed ShouldBind, %+v", err)
				logger.AppLogger().Warnf("BaseService Enter, %+v", err1)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err1.Error()}
			}

			// logger.AppLogger().Debugf("BaseService Enter, lanInvokeReq:%+v", lanInvokeReq)
			svc.EncryptedLanRequestBodyData = []byte(lanInvokeReq.Body)
		}

		var bodyDataDec []byte
		if config.Config.EncryptLanSessionData {
			logger.AppLogger().Debugf("BaseService Enter, EncryptLanSessionData")

			if svc.EncryptedLanRequestBodyData == nil {
				err := fmt.Errorf("EncryptedLanRequestBodyData is nil")
				logger.AppLogger().Warnf("BaseService.Enter, err:%+v", err)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err.Error()}
			}

			err := encwrapper.Check()
			if err != nil {
				logger.AppLogger().Warnf("check, err:%+v", err)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err.Error()}
			}

			bodyDataDec, err = encwrapper.DecParam(string(svc.EncryptedLanRequestBodyData))
			if err != nil {
				logger.AppLogger().Warnf("BaseService.Enter, DecParam err:%+v", err)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err.Error()}
			}
		} else {
			bodyDataDec = svc.EncryptedLanRequestBodyData
		}

		// logger.AppLogger().Debugf("BaseService Enter, bodyDataDec:%+v", string(bodyDataDec))
		err := encoding.JsonDecode(bodyDataDec, reqObj)
		if err != nil {
			logger.AppLogger().Warnf("BaseService.Enter, JsonDecode err:%+v", err)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err.Error()}
		}

	} else if svc.CalledType == CalledType_Gateway && reqObj != nil {
		logger.AppLogger().Debugf("BaseService Enter, CalledType_Gateway")

		if err := svc.ginContext.ShouldBind(reqObj); err != nil {
			err1 := fmt.Errorf("failed ShouldBind, %+v", err)
			logger.AppLogger().Warnf("BaseService Enter, %+v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err.Error()}
		}
	}

	svc.Req = reqObj

	// logger.AppLogger().Debugf("svc.Req:%+v", svc.Req) // 可能包含隐私字段，暂时不输出请求。
	if svc.Req != nil {
		valueOfReq := reflect.ValueOf(svc.Req)
		// logger.AppLogger().Debugf("valueOfReq:%+v", valueOfReq)
		reqExtraInfoValue := reflect.Indirect(valueOfReq).FieldByName("ReqExtraInfo")
		logger.AppLogger().Debugf("reqExtraInfoValue:%+v", reqExtraInfoValue)
		if reqExtraInfoValue.IsValid() { // 如果请求结构定义中没有 ReqExtraInfo 字段，则不会继续执行.
			agentTokenValue := reflect.Indirect(reqExtraInfoValue).FieldByName("AgentToken")
			logger.AppLogger().Debugf("agentTokenValue:%+v", agentTokenValue)
			if agentTokenValue.IsValid() {
				// 获取interface{}类型的值, 通过类型断言转换
				var agentToken string = agentTokenValue.Interface().(string)
				logger.AppLogger().Debugf("agentToken:%+v", agentToken)

				if err := VerifyAgentToken(agentToken); err != nil {
					err1 := fmt.Errorf("failed VerifyAgentToken, %+v", err)
					logger.AppLogger().Warnf("BaseService Enter, %+v", err1)
					return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
						Message: err.Error()}
				}
			}
		}
	}

	// logger.AppLogger().Debugf("BaseService, req: %+v", reqObj)
	// logger.AccessLogger().Debugf("[BaseService.Enter] svc.RequestId:%v, svc.Req:%+v", svc.RequestId, svc.Req)
	return iService.Process()
}

// 离开函数。主要是加密返回。
func (svc *BaseService) Process() dto.BaseRspStr {
	// logger.AppLogger().Debugf("BaseService Process")
	return svc.Leave()
}

// 去除json中的转义字符
func JsonEncodeWithDisableEscapeHtml(data interface{}) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(data); err != nil {
		return "", err
	}
	return bf.String(), nil
}

// 处理函数。主要是业务处理。
func (svc *BaseService) Leave() dto.BaseRspStr {
	if svc.Rsp != nil {
		// buf, err := encoding.JsonEncode(svc.Rsp)
		buf, err := JsonEncodeWithDisableEscapeHtml(svc.Rsp)
		if err != nil {
			logger.AppLogger().Debugf("BaseService Leave, JsonEncode %+v failed, err:%v", svc.Rsp, err)
		} else {
			logger.AppLogger().Debugf("====== BaseService Leave, rsp:" + buf)
			logger.AccessLogger().Debugf("[BaseService.Leave] svc.RequestId:%v, svc.Rsp:%+v", svc.RequestId, buf)
		}
	} else {
		logger.AppLogger().Debugf("BaseService Leave, rsp empty")
		logger.AccessLogger().Debugf("[BaseService.Leave] svc.RequestId:%v, svc.Rsp:%+v", svc.RequestId, svc.Rsp)
	}

	if config.Config.EncryptLanSessionData && svc.CalledType != CalledType_Gateway {
		if svc.RspBytes != nil {
			rsp, _ := encwrapper.EncBytes(svc.RequestId, svc.RspBytes)
			return rsp
		} else {
			rsp, _ := encwrapper.Enc(svc.Rsp)
			return rsp
		}
	}

	if svc.RspBytes != nil {
		rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr,
			Message: "OK",
			Results: svc.RspBytes}
		return rsp
	} else {
		rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr,
			Message: "OK",
			Results: svc.Rsp}
		// logger.AppLogger().Debugf("BaseService Enter, Rsp:%+v", svc.Rsp)
		return rsp
	}
}
