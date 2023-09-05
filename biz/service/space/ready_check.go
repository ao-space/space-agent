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

package space

import (
	"agent/biz/model/clientinfo"
	"agent/biz/model/disk_initial/model"
	"agent/biz/model/dto"
	"agent/biz/model/dto/space"
	"agent/biz/service/base"
	"agent/utils/logger"
)

type ReadyCheckService struct {
	base.BaseService
}

func (svc *ReadyCheckService) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("ReadyCheckService Process")

	paired := clientinfo.GetAdminPairedStatus()
	//diskInitialInfo := model.ReadDiskInitialInfo()
	//if diskInitialInfo.DiskInitialCode == model.DiskInitialCode_Nomal {
	//
	//	disks, err := info.GetDiskInfos(device_ability.GetAbilityModel().SupportUSBDisk)
	//	if err != nil {
	//		logger.AppLogger().Warnf("GetDiskInfos err:%+v", err)
	//		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
	//			Message: err.Error()}
	//	}
	//
	//	// 检测每一块磁盘是否还存并且接口号正确
	//	missingMainStorage, _ := manager.GetMissingDisk(diskInitialInfo, disks)
	//	if missingMainStorage {
	//		err := fmt.Errorf("miss main storage")
	//		logger.AppLogger().Warnf("%+v", err)
	//		return dto.BaseRspStr{Code: dto.AgentCodeMissingMainStorageFailedStr,
	//			Message: err.Error()}
	//	}
	//}
	rsp := &space.ReadyCheckRsp{Paired: paired,
		DiskInitialCode: model.DiskInitialCode_Nomal}
	svc.Rsp = rsp
	return svc.BaseService.Process()
}
