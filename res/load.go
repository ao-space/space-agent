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

package res

import (
	"agent/biz/model/device_ability"
	"agent/config"
	"agent/utils/tools"
	"os"
	"strings"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"

	_ "embed"
)

//go:embed docker-compose_run_as_docker.yml
var Content_docker_compose_run_as_docker_yml []byte

//go:embed aospace-upgrade.yml
var Content_aospace_upgrade_yml []byte

//go:embed agent-routers.json
var Content_agent_routers []byte

//go:embed pre-up-containers.json
var Content_pre_up_containers []byte

//go:embed static_html.zip
var Content_static_html_zip []byte

func GetContentDockerCompose() []byte {
	logger.AppLogger().Debugf("GetContentDockerCompose")

	var current_yml_content []byte
	m := device_ability.GetAbilityModel()
	if m != nil {
		// fmt.Printf("--------------- DeviceModelNumber:%+v\n", m.DeviceModelNumber)
		//if m.DeviceModelNumber >= device_ability.SN_GEN_2 {
		//	current_yml_content = Content_docker_compose_gen2_yml
		//} else
		if m.DeviceModelNumber <= device_ability.SN_GEN_CLOUD_DOCKER {
			current_yml_content = Content_docker_compose_run_as_docker_yml
		}

	}

	if m.DeviceModelNumber <= device_ability.SN_GEN_CLOUD_DOCKER {
		current_yml_content = disposeDockerComposeWhenRunInDocker(current_yml_content)
	}

	return current_yml_content
}

func GetContentUpgradeComposeFile() []byte {
	current_yml_content := Content_aospace_upgrade_yml
	logger.AppLogger().Debugf("upgrade compose file:%s", string(current_yml_content))
	m := device_ability.GetAbilityModel()
	if m.DeviceModelNumber <= device_ability.SN_GEN_CLOUD_DOCKER {
		current_yml_content = disposeDockerComposeWhenRunInDocker(current_yml_content)
	}
	return current_yml_content
}

func disposeDockerComposeWhenRunInDocker(content_docker_compose []byte) []byte {

	// 从环境变量中获取 docker-compose.yml 中的 volumes 宿主机的挂载目录.
	hostDataPath := "/run/desktop/mnt/host/c/aospace" // 默认目录, macOS/Linux 必须要用户传入.
	envkey := config.Config.Box.RunInDocker.AoSpaceDataDirEnv
	dataDir := os.Getenv(envkey)
	if len(dataDir) > 0 {
		dataDir = strings.ReplaceAll(dataDir, "\\", "/")
		hostDataPath = dataDir
	}
	hostDataPath = fileutil.AddPathSepIfNeed(hostDataPath)
	isWindowsHost := strings.Contains(hostDataPath, "/run/desktop/mnt/host")

	// 挂载目录修改
	newCompose := ""
	s := string(content_docker_compose)
	preLine := ""
	lines := tools.StringToLines(s)
	for _, v := range lines {
		// logger.AppLogger().Debugf("line:%v", v)

		if strings.Index(strings.TrimSpace(v), "#") == 0 { // 注释不添加到新的 yml 文件.
			continue
		}

		// 挂载目录修改,  windows 修改成类似这样 /run/desktop/mnt/host/c/aospace/.... , macOS 类似于 /Users/nist/aospace
		if strings.Contains(preLine, "volumes:") || strings.Index(strings.TrimSpace(preLine), "- /") == 0 {
			if strings.Index(strings.TrimSpace(v), "- /") == 0 { //路径情况时, 比如 - /home/eulixspace_link/data/third_party:/data

				if strings.Index(strings.TrimSpace(v), `- /var/run/docker.sock:/var/run/docker.sock`) == 0 {
					if isWindowsHost { // 只有 windows 才需要修改 docker.sock 路径
						// - "/var/run/docker.sock:/var/run/docker.sock" 改成 - "//var/run/docker.sock:/var/run/docker.sock"
						v = strings.Replace(v, `- /var/run/docker.sock:/var/run/docker.sock`, `- //var/run/docker.sock:/var/run/docker.sock`, 1)
					}

				} else {
					// 路径增加传入的数据目录， 比如 windows 时是 "/run/desktop/mnt/host/c/aospace".
					v = strings.Replace(v, "- /", "- "+hostDataPath, 1)
				}
			}
		} else if strings.Contains(v, "APP_ROOT_PATH:") { // 开发者选项安装的容器数据目录环境变量 APP_ROOT_PATH 对应的值加上传入的数据目录
			arr := strings.Split(v, ":")
			if len(arr) < 2 {
				continue
			}
			mountP := hostDataPath + arr[1]
			mountP = strings.Replace(mountP, "//", "/", 1)
			mountP = strings.Replace(mountP, "/ /", "/", 1)
			mountP = strings.Replace(mountP, " /", "/", 1)
			v = strings.Replace(v, arr[1], " "+mountP, 1)
		}

		newCompose = newCompose + v + "\n"
		preLine = v
	}

	return []byte(newCompose)
}

func GetContentAgentRouters() []byte {
	return Content_agent_routers
}

func GetContentPreUpContainers() []byte {
	return Content_pre_up_containers
}

func GetContentStaticHtmlZip() []byte {
	return Content_static_html_zip
}
