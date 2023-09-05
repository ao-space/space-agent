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

package network

import (
	"agent/utils/logger"
	"fmt"
	"strings"

	"github.com/dungeonsnd/gocom/sys/run"
)

func runCmd(params []string) error {
	cmd := "nmcli"
	logger.AppLogger().Debugf("will run cmd: %v %v", cmd, strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe(cmd, params)
	if err != nil {
		return fmt.Errorf("failed run %v %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			cmd, params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("run cmd return, %v %v, stdOutput is :%v, errOutput is :%v",
		cmd, strings.Join(params, " "), string(stdOutput), string(errOutput))
	return nil
}

func runCmd2(cmd string, params []string) error {
	logger.AppLogger().Debugf("will run cmd: %v %v\n", cmd, strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe(cmd, params)
	if err != nil {
		return fmt.Errorf("failed run %v %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			cmd, params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("run cmd return, %v %v, stdOutput is :%v, errOutput is :%v",
		cmd, strings.Join(params, " "), string(stdOutput), string(errOutput))
	return nil
}

func runCmdOutput(params []string) ([]byte, error) {
	cmd := "nmcli"
	logger.AppLogger().Debugf("will run cmd: %v %v", cmd, strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe(cmd, params)
	if err != nil {
		return nil, fmt.Errorf("failed run %v %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			cmd, params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("run cmd return, %v %v, stdOutput is :%v, errOutput is :%v",
		cmd, strings.Join(params, " "), string(stdOutput), string(errOutput))
	return stdOutput, nil
}
