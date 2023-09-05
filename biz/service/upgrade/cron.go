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
	"agent/biz/db"
	"agent/utils/hardware"
	"agent/utils/logger"
	"github.com/robfig/cron/v3"
)

func CronForUpgrade() {
	err := db.CheckAndCreateDB()
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to create datafile, exit program. err: %s", err)
	}
	c := cron.New()
	var agent Agent
	if hardware.RunningInDocker() {
		agent = new(ContainerAgent)
	} else {
		agent = new(NativeAgent)
	}
	_, err = c.AddFunc("00 02 * * *", agent.AutoUpgrade)
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to config cron: %s", err)
	}
	c.Start()
}

//func WriteUpgradeSelfCron() error {
//	cronContent, err := ioutil.ReadFile("/etc/crontab")
//	if err != nil {
//		logger.AppLogger().Errorf("read cron file error:%v", err)
//		return err
//	}
//	if !strings.Contains(string(cronContent), "eulixspace-upgrade") {
//		task, err := db.ReadTask("")
//		if task.Status == upgrade.Installed || task.Status == upgrade.InstallErr {
//			err = dnfUpdateUpgradeToolsSelf()
//			if err != nil {
//				logger.UpgradeLogger().Errorf("%s self update error: %v", UpgradeName, err)
//			}
//		}
//		//} else if task.Status == upgrade.Installing {
//		//	time.Sleep(3 * time.Minute)
//		//	err = dnfUpdateUpgradeToolsSelf()
//		//	if err != nil {
//		//		logger.UpgradeLogger().Errorf("%s self update error: %v", UpgradeName, err)
//		//	}
//		//}
//		cronStr := "\n1 0 * * * root /usr/bin/dnf clean all && /usr/bin/dnf makecache && /usr/bin/dnf update -y eulixspace-upgrade \n"
//		cronBytes := []byte(cronStr)
//		f, err := os.OpenFile("/etc/crontab", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
//		if err != nil {
//			return err
//		}
//		defer f.Close()
//		if _, err := f.Write(cronBytes); err != nil {
//			return err
//		}
//		return nil
//	}
//	return nil
//}
