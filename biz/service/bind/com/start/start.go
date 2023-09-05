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

package start

import (
	"agent/biz/docker"
	"agent/biz/model/clientinfo"
	"agent/biz/model/dto"
	"agent/biz/service/base"
	"agent/utils/logger"
	"fmt"
)

type ComStartService struct {
	base.BaseService
}

func (svc *ComStartService) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("ComStartService Process")

	pairedStatus := clientinfo.GetAdminPairedStatus()
	dockerStatus := docker.GetDockerStatus()

	if pairedStatus == clientinfo.DeviceAlreadyBound {
		err := fmt.Errorf("pairedStatus:%+v", pairedStatus)
		return dto.BaseRspStr{Code: dto.AgentCodeAlreadyPairedStr,
			Message: err.Error()}
	}
	if dockerStatus == docker.ContainersStarting || dockerStatus == docker.ContainersDownloading {
		err := fmt.Errorf("dockerStatus:%+v", dockerStatus)
		return dto.BaseRspStr{Code: dto.AgentCodeDockerStarting,
			Message: err.Error()}
	}
	if dockerStatus == docker.ContainersStarted {
		err := fmt.Errorf("dockerStatus:%+v", dockerStatus)
		return dto.BaseRspStr{Code: dto.AgentCodeDockerStarted,
			Message: err.Error()}
	}

	go docker.PostEvent(docker.EventPairing)

	return svc.BaseService.Process()
}
