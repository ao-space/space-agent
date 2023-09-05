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

package method

import (
	"agent/biz/model/did"
	"agent/biz/model/did/leveldb"
	"agent/biz/model/dto"
	"agent/biz/model/dto/did/document/method"
	"agent/biz/service/base"
	"agent/utils/logger"
	"encoding/base64"
	"fmt"
)

type UpdateDocumentMethod struct {
	base.BaseService
}

func NewUpdateDocumentMethod() *UpdateDocumentMethod {
	svc := new(UpdateDocumentMethod)
	return svc
}

func (svc *UpdateDocumentMethod) Process() dto.BaseRspStr {
	req := svc.Req.(*method.UpdateDocumentMethodReq)
	logger.AppLogger().Debugf("UpdateDocumentMethod Process, svc.RequestId:%v, req:%+v", svc.RequestId, req)
	if req == nil {
		err1 := fmt.Errorf("request error")
		logger.AppLogger().Debugf(err1.Error())
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, RequestId: svc.RequestId, Message: err1.Error()}
	}

	levelDBTrans, err := leveldb.BeginTransaction() // 开启事务
	defer levelDBTrans.Rollback()                   // 退出时回滚事务. 如果成功, 函数返回之前主动 commit.

	didDocBytes, newDidStr, err := did.ResetPasswordVerficationMethod(levelDBTrans, req.DID, req.AOID, req.NewPassword)
	if err != nil {
		err1 := fmt.Errorf("ResetPasswordVerficationMethod err:%v", err)
		logger.AppLogger().Debugf(err1.Error())
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, RequestId: svc.RequestId, Message: err1.Error()}
	}

	svc.Rsp = &method.UpdateDocumentMethodRsp{DIDDoc: base64.StdEncoding.EncodeToString(didDocBytes),
		DID: newDidStr}
	if levelDBTrans != nil {
		levelDBTrans.Commit() // commit 提交事务
	}
	return svc.BaseService.Process()
}
