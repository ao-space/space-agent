package identify

import (
	"agent/biz/model/device"
	"agent/biz/model/dto"
	"agent/biz/model/dto/bind/identify"
	"agent/biz/service/base"
	"agent/biz/service/pair"
	"agent/config"
	"agent/utils/logger"
	utilshttp "agent/utils/network/http"
	"github.com/dungeonsnd/gocom/encrypt/random"
	"net/http"
	"time"
)

type TicketService struct {
	base.BaseService
}

func (svc *TicketService) Process() dto.BaseRspStr {
	req := svc.Req.(*identify.TicketReq)
	ticket, err := GetTempTicket(req.PlatformUrl)
	if err != nil {
		return dto.BaseRspStr{Code: dto.AgentGetTicketFailed, Message: err.Error()}
	}

	svc.Rsp = ticket

	return svc.BaseService.Process()
}

func GetTempTicket(platformUrl string) (*identify.TicketRsp, error) {
	//先获取 box-reg-key
	if boxRegKeyInfo, err := pair.GetDeviceRegKey("https://" + platformUrl); err != nil {
		logger.AppLogger().Warnf("ServiceRegisterBox, failed GetBoxRegKey, err:%+v", err)
		return nil, err
	} else {
		logger.AppLogger().Debugf("ServiceRegisterBox, succ GetBoxRegKey, boxRegKeyInfo:%+v", boxRegKeyInfo)
		device.SetDeviceRegKey(boxRegKeyInfo.BoxRegKey, boxRegKeyInfo.ExpiresAt)
	}
	// 平台请求结构
	type registryStruct struct {
		BoxUUID string `json:"boxUUID"`
	}
	// 平台响应结构
	type ticketRspStruct struct {
		Ticket    string `json:"ticket"`
		ExpiresAt string `json:"expiresAt"`
	}

	// 请求平台
	parms := &registryStruct{BoxUUID: device.GetDeviceInfo().BoxUuid}
	// url := config.Config.Platform.APIBase.Url + config.Config.Platform.RegistryBox.Path
	url := "https://" + platformUrl + config.Config.Platform.GetChannelTicket.Path
	logger.AppLogger().Debugf("ServiceGetChannelTempTicket, v2, url:%+v, parms:%+v", url, parms)

	var headers = map[string]string{"Request-Id": random.GenUUID(), "Box-Reg-Key": device.GetDeviceInfo().BoxRegKey}
	var rsp identify.TicketRsp

	tryTotal := 3
	// var httpReq *http.Request
	var httpRsp *http.Response
	var body []byte
	var err1 error
	for i := 0; i < tryTotal; i++ {
		_, httpRsp, body, err1 = utilshttp.PostJsonWithHeaders(url, parms, headers, &rsp)
		if err1 != nil {
			// logger.AppLogger().Warnf("Failed PostJson, err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@body:%v", err1, httpReq, httpRsp, string(body))
			if i == tryTotal-1 {
				return nil, err1
			}
			time.Sleep(time.Second * 2)
			continue
		} else {
			break
		}
	}

	logger.AppLogger().Infof("ServiceRegisterBox, parms:%+v", parms)
	logger.AppLogger().Infof("ServiceRegisterBox, rsp:%+v", rsp)
	// logger.AppLogger().Infof("ServiceRegisterBox, httpReq:%+v", httpReq)
	logger.AppLogger().Infof("ServiceRegisterBox, httpRsp:%+v", httpRsp)
	logger.AppLogger().Infof("ServiceRegisterBox, body:%v", string(body))

	return &rsp, nil
}
