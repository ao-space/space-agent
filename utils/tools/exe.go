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
 * @Date: 2021-12-09 10:25:41
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-30 11:27:28
 * @Description:
 */
package tools

import (
	"agent/biz/model/dto"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/sys/run"
)

func ExeCmd(name string, arg ...string) (string, string, error) {
	// logger.AppLogger().Debugf("exec: %s %s", name, arg)
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	return outStr, errStr, err
}

func RunCmd(cmd string, params []string) (dto.BaseRspStr, string, error) {
	if params == nil {
		params = []string{}
	}
	logger.AppLogger().Debugf("will run cmd: %v %v", cmd, strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe(cmd, params)
	if err != nil {
		err1 := fmt.Errorf("failed run %v %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			cmd, params, err, string(stdOutput), string(errOutput))
		logger.AppLogger().Warnf(err1.Error())
		rsp := dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err1.Error()}
		return rsp, "", err1
	}
	logger.AppLogger().Debugf("run return, cmd: %v %v, stdOutput:%v, errOutput:%v",
		cmd, strings.Join(params, " "), string(stdOutput), string(errOutput))

	rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		Message: "OK"}
	return rsp, string(stdOutput), nil
}
