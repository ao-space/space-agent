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

package version

import (
	"agent/utils/tools"
	"fmt"
	"strings"

	"agent/utils/logger"
)

const AgentName = "eulixspace-agent"

// GetInstalledAgentVersion from rpm
func GetInstalledAgentVersion() (string, error) {
	queryCmd := fmt.Sprintf(`dnf list installed %s | grep %s.aarch64  | awk '{print $2}'`, AgentName, AgentName)
	stdout, stderr, err := tools.ExeCmd("bash", "-c", queryCmd)
	if err != nil || stderr != "" {
		return stdout, fmt.Errorf("get agent version though rpm %s: %s", err, stderr)
	}

	return stdout, nil
}

// https://pm.eulix.xyz/bug-view-1074.html
func GetInstalledAgentVersionRemovedNewLine() string {
	boxVersion, err := GetInstalledAgentVersion()
	if err != nil {
		logger.AppLogger().Errorf("putEnvIntoGateway : %s", err)
		return ""
	}

	boxVersion = strings.ReplaceAll(boxVersion, "\r", "\n")
	boxVersion = strings.ReplaceAll(boxVersion, "\n\n", "\n")
	arr := strings.Split(boxVersion, "\n")
	if len(arr) > 1 && len(arr[len(arr)-2]) > 0 {
		boxVersion = arr[len(arr)-2]
	} else if len(arr) > 0 {
		boxVersion = arr[len(arr)-1]
	}

	return boxVersion
}
