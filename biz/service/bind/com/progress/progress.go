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

package progress

import (
	"agent/biz/docker"
	"agent/biz/model/clientinfo"
	"agent/biz/model/dto"
	"agent/biz/model/dto/bind/com/progress"
	"agent/biz/service/base"
	"agent/config"
	"agent/utils/logger"
	"agent/utils/tools"
	"fmt"
	"os"
	"strings"
	"time"
)

type ComProgressService struct {
	base.BaseService
	PairedInfo *clientinfo.AdminPairedInfo
}

func (svc *ComProgressService) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("ComProgressService Process")
	svc.PairedInfo = clientinfo.GetAdminPairedInfo()
	if svc.PairedInfo.AlreadyBound() {
		err := fmt.Errorf("pairedStatus:%+v", svc.PairedInfo.Status())
		return dto.BaseRspStr{Code: dto.AgentCodeAlreadyPairedStr,
			Message: err.Error()}
	}

	singleDockerModeEnv := os.Getenv(config.Config.Box.RunInDocker.AoSpaceSingleDockerModeEnv)
	logger.AppLogger().Debugf("ComProgressService Process, singleDockerModeEnv:%v", singleDockerModeEnv)
	if strings.EqualFold(singleDockerModeEnv, "true") {
		go svc.startSpacePrograms()
		time.Sleep(12 * time.Second)
		rsp := &progress.ProgressRsp{ComStatus: docker.ContainersStarted, Progress: 100}
		svc.Rsp = rsp
		return svc.BaseService.Process()
	}

	dockerStatus := docker.GetDockerStatus()
	rsp := &progress.ProgressRsp{ComStatus: dockerStatus, Progress: docker.GetStartingProgress()}
	svc.Rsp = rsp
	return svc.BaseService.Process()
}

func (svc *ComProgressService) startSpacePrograms() error {
	cmd := "bash"
	scripts := []string{"/usr/local/bin/start.sh"}
	logger.AppLogger().Debugf("startSpacePrograms, will RunCmd script: %v %v", cmd, strings.Join(scripts, " "))
	_, stdout, err := tools.RunCmd(cmd, scripts)
	if err != nil {
		logger.AppLogger().Warnf("%v %v, stdout:%v, err:%v", cmd, strings.Join(scripts, " "), stdout, err)
	} else {
		logger.AppLogger().Debugf("%v %v succ", cmd, strings.Join(scripts, " "))
	}
	return err
}
