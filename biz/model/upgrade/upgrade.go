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
 * @Date: 2021-11-10 14:41:01
 * @LastEditors: xuyang
 * @LastEditTime: 2023-05-15 16:52:59
 * @Description:
 */

package upgrade

import (
	"reflect"
	"time"
)

const (
	Ing         = "ing"
	Done        = "done"
	Err         = "err"
	Downloading = "downloading"
	Downloaded  = "downloaded"
	Installing  = "installing"
	Installed   = "installed"
	DownloadErr = "download-err"
	InstallErr  = "install-err"
)

type Task struct {
	VersionId        string          `json:"versionId"`
	Status           string          `json:"status"`        // 整体流程状态："", downloading, downloaded, installing, installed, download-err，install-err
	DownStatus       string          `json:"downStatus"`    // 下载状态："", ing, done, err
	InstallStatus    string          `json:"installStatus"` // 安装状态："", ing, done, err
	StartDownTime    string          `json:"startDownTime"`
	StartInstallTime string          `json:"startInstallTime"`
	DoneDownTime     string          `json:"doneDownTime"`
	DoneInstallTime  string          `json:"doneInstallTime"`
	RpmPkg           VersionDownInfo `json:"rpmPkg"`
	CFile            VersionDownInfo `json:"cFile"` // docker-compose.yml
	ContainerImg     VersionDownInfo `json:"containerImg"`
	KernelImg        VersionDownInfo `json:"KernelImg"`
	NeedReboot       bool            `json:"reboot"`
}

type TaskGoal struct {
	VersionId string `json:"versionId"`
	TaskGoal  string `json:"taskGoal"` // onlyDownload， onlyInstall， downloadAndInstall
}

type Schedule struct {
	Rate   float64 `json:"rate"`
	Detail string  `json:"detail"`
}

type UpgradeConfig struct {
	AutoDownload bool `json:"autoDownload"`
	AutoInstall  bool `json:"autoInstall"`
}

type VersionFromPlatform struct {
	PkgName           string `json:"pkgName"`
	PkgType           string `json:"pkgType"`
	PkgVersion        string `json:"pkgVersion"`
	PkgSize           int64  `json:"pkgSize"`
	DownloadUrl       string `json:"downloadUrl"`
	UpdateDesc        string `json:"updateDesc"`
	Md5               string `json:"md5"`
	IsForceUpdate     bool   `json:"isForceUpdate"`
	MinAndroidVersion string `json:"minAndroidVersion"`
	MinIOSVersion     string `json:"minIOSVersion"`
	MinBoxVersion     string `json:"minBoxVersion"`
}

type VersionFromPlatformV2 struct {
	Id                int64     `json:"id"`
	PkgName           string    `json:"pkgName"`
	PkgType           string    `json:"pkgType"`
	PkgVersion        string    `json:"pkgVersion"`
	PkgSize           int64     `json:"pkgSize"`
	DownloadUrl       string    `json:"downloadUrl"`
	UpdateDesc        string    `json:"updateDesc"`
	Md5               string    `json:"md5"`
	IsForceUpdate     bool      `json:"isForceUpdate"`
	MinAndroidVersion string    `json:"minAndroidVersion"`
	MinIOSVersion     string    `json:"minIOSVersion"`
	MinBoxVersion     string    `json:"minBoxVersion"`
	CreateAt          time.Time `json:"createAt"`
	UpdateAt          time.Time `json:"updateAt"`
	Restart           bool      `json:"restart"`
	KernelVersion     string    `json:"kernelVersion"`
	KernelUrl         string    `json:"kernelDownloadUrl"`
	KernelMd5         string    `json:"kernelMd5"`
	KernelSize        string    `json:"kernelSize"`
}

type StartDownRes struct {
	VersionId string `json:"versionId" required:"true"`
	Anew      bool   `json:"anew"` // 重新下载：当已经有任务在下载的时候，用户仍然想重新开始下载，使用此选项
}

type StartUpgradeRes struct {
	VersionId string `json:"versionId" required:"true"`
}

type VersionDownInfo struct {
	VersionId  string    `json:"versionId"`
	Downloaded bool      `json:"downloaded"`
	PkgPath    string    `json:"pkgPath"`
	UpdateTime time.Time `json:"updateTime"`
}

type OverallInfo struct {
	VersionId  string    `json:"versionId"`
	Downloaded bool      `json:"downloaded"`
	PkgPath    string    `json:"pkgPath"`
	UpdateTime time.Time `json:"updateTime"`
	Restart    bool      `json:"restart"`
	KernelInfo `json:"kernelInfo"`
}

type KernelInfo struct {
	KernelVersion string `json:"kernelVersion"`
	KernelUrl     string `json:"kernelDownloadUrl"`
	KernelMd5     string `json:"kernelMd5"`
}

type TimeTransformer struct {
}

func (t TimeTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}

type AllInOneUpgradeReq struct {
	VersionId string `json:"versionId"`
	DataDir   string `json:"dataDir"`
}
