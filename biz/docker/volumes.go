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

package docker

import (
	"agent/config"
	"agent/utils/logger"
	"agent/utils/tools"
	"strings"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

func ProcessVolumes(composeFile string, diskParts []string) error {
	logger.AppLogger().Debugf("ProcessVolumes, diskParts:%+v", diskParts)

	content, err := fileutil.ReadFromFile(composeFile)
	if err != nil {
		return err
	}
	s := string(content)

	placeholderInHost := config.Config.Box.Disk.DockerVolumePlaceholderInHost
	placeholderInContainer := config.Config.Box.Disk.DockerVolumePlaceholderInContainer

	lines := tools.StringToLines(s)
	logger.AppLogger().Debugf("ProcessVolumes, lines count:%+v", len(lines))
	newContent := ""
	for _, line := range lines {
		if strings.Index(line, placeholderInHost) > 0 {
			if diskParts == nil || len(diskParts) < 1 { // 绑定之后、磁盘尚未初始化, 不需要挂载磁盘目录.
				// logger.AppLogger().Debugf("ProcessVolumes, continue")
				continue

			} else { // 绑定之后、磁盘已经初始化
				for _, diskPart := range diskParts {
					n := strings.ReplaceAll(line, placeholderInHost, diskPart)
					n = strings.ReplaceAll(n, placeholderInContainer, diskPart)
					newContent += n
					newContent += "\n"
				}
			}

		} else {
			newContent += line
			newContent += "\n"
		}
	}
	if len(newContent) > 0 {
		return fileutil.WriteToFile(composeFile, []byte(newContent), true)
	}

	return nil
}
