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

package document

import (
	"agent/biz/model/did"
	"agent/biz/model/dto"
	"agent/biz/model/dto/did/document"
	"agent/biz/service/base"
	"agent/utils/logger"
	"encoding/base64"
	"fmt"
)

type GetDocument struct {
	base.BaseService
}

func NewGetDocument() *GetDocument {
	svc := new(GetDocument)
	return svc
}

func (svc *GetDocument) Process() dto.BaseRspStr {
	req := svc.Req.(*document.GetDocumentReq)
	// logger.AppLogger().Debugf("GetDocument Process, svc.RequestId:%v, req:%+v", svc.RequestId, req)

	didDocBytes, err := did.GetDocumentFromFile(nil, req.AOID, req.DID)
	if err != nil {
		err1 := fmt.Errorf("GetDocument err:%v", err)
		logger.AppLogger().Debugf(err1.Error())
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, RequestId: svc.RequestId, Message: err1.Error()}
	}

	encryptedPriKeyBytes, found, err := did.GetEncryptedPriKeyBytes(nil, req.AOID)
	if err != nil {
		err1 := fmt.Errorf("GetDocument err:%v", err)
		logger.AppLogger().Debugf(err1.Error())
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, RequestId: svc.RequestId, Message: err1.Error()}
	}

	rsp := &document.GetDocumentRsp{}
	rsp.DIDDoc = base64.StdEncoding.EncodeToString(didDocBytes)
	if found {
		rsp.EncryptedPriKeyBytes = base64.StdEncoding.EncodeToString(encryptedPriKeyBytes)
	}
	svc.Rsp = rsp
	return svc.BaseService.Process()
}
