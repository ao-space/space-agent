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

package system

import (
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/service/base"
	"fmt"
	"strings"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/sys/run"
)

type ShutdownService struct {
	base.BaseService
}

func (svc *ShutdownService) Process() dto.BaseRspStr {
	abilityModel := device_ability.GetAbilityModel()
	if !abilityModel.InnerDiskSupport {
		err := fmt.Errorf("unsupported function")
		return dto.BaseRspStr{Code: dto.AgentCodeUnsupportedFunction,
			Message: err.Error()}
	}

	go func() {
		time.Sleep(time.Microsecond * 1000)
		_, err := shutdown()
		if err != nil {
			logger.AppLogger().Warnf("shutdown err")
		}
	}()
	return svc.BaseService.Process()
}

func shutdown() (dto.BaseRspStr, error) {
	cmd := "shutdown"
	params := []string{"-Ph", "now"}
	stdOutput, errOutput, err := run.RunExe(cmd, params)
	if err != nil {
		err1 := fmt.Errorf("failed run %v %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			cmd, params, err, string(stdOutput), string(errOutput))
		logger.AppLogger().Warnf(err1.Error())
		rsp := dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err1.Error()}
		return rsp, err1
	} else {
		logger.AppLogger().Debugf("run cmd: %v %v, stdOutput:%v, errOutput:%v",
			cmd, strings.Join(params, " "), string(stdOutput), string(errOutput))
	}

	rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		Message: "OK"}
	return rsp, nil
}
