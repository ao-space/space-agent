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

package space

import (
	"agent/utils/tools"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/dungeonsnd/gocom/sys/run"
)

func FolderUsage(folder string) (uint64, uint64, uint64, error) {
	return GetFolderSize(folder)
}

func AllPartsUsage() (uint64, uint64, uint64, error) {
	total, err := Total()
	if err != nil {
		return 0, 0, 0, err
	}
	avail, err := Avail()
	if err != nil {
		return 0, 0, 0, err
	}
	return total, total - avail, avail, nil
}

// root@61828da6640b:/aospace# df  --block-size=1 /aospace
// Filesystem        1B-blocks         Used    Available Use% Mounted on
// C:\            510770802688 341898379264 168872423424  67% /aospace
//
// root@61828da6640b:/aospace# df -h /aospace
// Filesystem      Size  Used Avail Use% Mounted on
// C:\             476G  319G  158G  67% /aospace

// return total, used, avail
func GetFolderSize(folder string) (uint64, uint64, uint64, error) {

	params := []string{"--block-size=1", folder}
	logger.AppLogger().Debugf("GetFolderSize, run cmd: df")
	stdOutput, errOutput, err := run.RunExe("df", params)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed run GetFolderSize df, err is :%v, stdOutput is :%v, errOutput is :%v",
			err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("GetFolderSize, run cmd: df %v, , stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))

	lines := tools.StringToLines(string(stdOutput))
	for _, line := range lines {
		if strings.Contains(line, folder) {
			arr := strings.Fields(line)
			if len(arr) < 4 {
				return 0, 0, 0, fmt.Errorf("failed run GetFolderSize df, line :%v", line)
			}

			total, err := strconv.ParseInt(arr[1], 10, 64)
			if err != nil {
				logger.AppLogger().Warnf("strconv.ParseInt, arr[1]:%v, err:%v", arr[1], err)
			}
			logger.AppLogger().Debugf("total=%v", total)

			used, err := strconv.ParseInt(arr[2], 10, 64)
			if err != nil {
				logger.AppLogger().Warnf("strconv.ParseInt, arr[2]:%v, err:%v", arr[2], err)
			}
			logger.AppLogger().Debugf("used=%v", used)

			avail, err := strconv.ParseInt(arr[3], 10, 64)
			if err != nil {
				logger.AppLogger().Warnf("strconv.ParseInt, arr[3]:%v, err:%v", arr[3], err)
			}
			logger.AppLogger().Debugf("avail=%v", avail)
			return uint64(total), uint64(used), uint64(avail), nil
		}
	}

	return 0, 0, 0, fmt.Errorf("failed run GetFolderSize df, not found  folder:%v", folder)
}

// [root@EulixOS ~]# df -l --block-size=1
// Filesystem        1B-blocks       Used    Available Use% Mounted on
// /dev/mmcblk0p4 107647705088    4878336 103247577088   1% /home
// 取 Available 字段
func Avail() (uint64, error) {

	params := []string{"-l", "--block-size=1"}
	logger.AppLogger().Debugf("Avail, run cmd: df")
	stdOutput, errOutput, err := run.RunExe("df", params)
	if err != nil {
		return 0, fmt.Errorf("failed run AllPartsUsage df, err is :%v, stdOutput is :%v, errOutput is :%v",
			err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("Avail, run cmd: df %v, , stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)

	var avail int64
	for _, line := range lines {
		// if strings.Contains(line, "mmc") && strings.Index(line, "/home") > 0 {
		if strings.Index(line, "/home") > 0 {
			fileds := strings.Fields(line)
			if len(fileds) > 4 {
				availStr := fileds[3]
				logger.AppLogger().Debugf("availStr=%v", availStr)
				avail, err = strconv.ParseInt(availStr, 10, 64)
				logger.AppLogger().Debugf("avail=%v", avail)
			}
			break
		}
	}

	return uint64(avail), nil
}

// [root@EulixOS ~]# lsblk -b
// NAME        MAJ:MIN RM         SIZE RO TYPE MOUNTPOINT
// mmcblk0     179:0    0 127865454592  0 disk
func Total() (uint64, error) {

	params := []string{"-b"}
	logger.AppLogger().Debugf("Total, will run cmd: lsblk")
	stdOutput, errOutput, err := run.RunExe("lsblk", params)
	if err != nil {
		return 0, fmt.Errorf("failed run Total lsblk, err is :%v, stdOutput is :%v, errOutput is :%v",
			err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("Total, run return, cmd: lsblk %v, errOutput is :%v",
		strings.Join(params, " "), string(errOutput))

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(string(stdOutput), -1)

	// 树莓派
	var total int64
	for _, line := range lines {
		if strings.Contains(line, "mmc") && strings.Index(line, "disk") > 0 {
			fileds := strings.Fields(line)
			if len(fileds) > 4 {
				totalStr := fileds[3]
				logger.AppLogger().Debugf("totalStr=%v", totalStr)
				total, err = strconv.ParseInt(totalStr, 10, 64)
				logger.AppLogger().Debugf("total=%v", total)
			}

			break
		}
	}
	if total > 0 {
		return uint64(total), nil
	}

	// x86
	for _, line := range lines {
		if strings.Index(line, "disk") > 0 {
			fileds := strings.Fields(line)
			if len(fileds) > 4 {
				totalStr := fileds[3]
				logger.AppLogger().Debugf("totalStr=%v", totalStr)
				total, err = strconv.ParseInt(totalStr, 10, 64)
				logger.AppLogger().Debugf("total=%v", total)
			}

			break
		}
	}

	return uint64(total), nil
}

type LsblkChildren struct {
	Name         string `json:"name"`
	Fsavail      string `json:"fsavail"`
	Fssize       string `json:"fssize"`
	Fstype       string `json:"fstype"`
	Fsused       string `json:"fsused"`
	FsusePercent string `json:"fsusePercent"`
	Mountpoint   string `json:"mountpoint"`
	Uuid         string `json:"uuid"`
	Ptuuid       string `json:"ptuuid"`
	Partuuid     string `json:"partuuid"`
	Model        string `json:"model"`
	Size         uint64 `json:"size"`

	Children []*LsblkChildren `json:"children"`
}

func LsblkToObject() ([]*LsblkChildren, error) {
	cmd := "lsblk"
	params := []string{"-fbJM", "-o", "NAME,FSAVAIL,FSSIZE,FSTYPE,FSUSED,FSUSE%,MOUNTPOINT,UUID,PTUUID,PARTUUID,MODEL,SERIAL,SIZE"}
	stdOutput, errOutput, err := run.RunExe(cmd, params)
	out := string(stdOutput)
	if err != nil {
		return nil, fmt.Errorf("TotalOfOneDisk, failed run %v %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			cmd, strings.Join(params, " "), err, out, string(errOutput))
	}
	logger.AppLogger().Debugf("Total, run return, %v %v, errOutput is :%v",
		cmd, strings.Join(params, " "), string(errOutput))

	// var obj LsblkMapObject
	var obj map[string][]*LsblkChildren
	err = encoding.JsonDecode(stdOutput, &obj)
	if err != nil {
		return nil, fmt.Errorf("TotalOfOneDisk, failed JsonDecode %v, err is :%v",
			string(errOutput), err)
	}

	v := obj["blockdevices"]
	return v, nil
}
