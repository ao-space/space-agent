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

package upgrade

import (
	"agent/utils/tools"
	"fmt"
)

//func GetInstalledAgentVersion() (string, error) {
//	queryCmd := fmt.Sprintf(`dnf list installed %s | grep %s.aarch64  | awk '{print $2}'`, "eulixspace-agent", "eulixspace-agent")
//	stdout, stderr, err := tools.ExeCmd("bash", "-c", queryCmd)
//	if err != nil || stderr != "" {
//		return stdout, fmt.Errorf("get agent version though rpm %s: %s", err, stderr)
//	}
//	versionId := strings.TrimSpace(stdout)
//	return versionId, nil
//}

func CleanAllAndMakeCache() error {
	_, stdout, err := tools.RunCmd("dnf", []string{"clean", "all"})
	_, stdout, err = tools.RunCmd("dnf", []string{"makecache"})
	if err != nil {
		return fmt.Errorf("dnf clean all and make cache error: %v: %v", err, stdout)
	}
	return nil
}

func GetCurFirmwareRpmFileName(pkgName string) string {
	_, rpmFileName, err := tools.RunCmd("rpm", []string{"-q", pkgName})
	if err != nil {
		return ""
	}
	return rpmFileName
}

func dnfUpdateUpgradeToolsSelf() error {
	err := CleanAllAndMakeCache()
	if err != nil {
		return err
	}
	_, stdout, err := tools.RunCmd("dnf", []string{"update", "-y", UpgradeName})
	if err != nil {
		return fmt.Errorf("dnf update %s error: %v: %v", UpgradeName, err, stdout)
	}
	return nil
}
