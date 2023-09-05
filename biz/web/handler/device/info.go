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
 * @LastEditors: jeffery
 * @LastEditTime: 2022-04-13 14:59:11
 * @Description:
 */

package device

import (
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/dto/device"
	"agent/config"
	"agent/utils/disk/space"
	"agent/utils/docker/dockerfacade"
	"agent/utils/logger"
	"agent/utils/version"
	"fmt"
	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/dungeonsnd/gocom/encrypt/random"
	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/dungeonsnd/gocom/sys/run"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

// Version godoc
// @Summary get version info [for mirco service]
// @Description get version info
// @ID device.BoxDeviceVersion
// @Tags device
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=device.BoxDeviceVersion} "code=AG-200 success;"
// @Router /agent/v1/api/device/version [GET]
func Version(c *gin.Context) {

	// agent version
	boxVersion := version.GetInstalledAgentVersionRemovedNewLine()
	if len(boxVersion) < 3 {
		boxVersion = config.VersionNumber
	}

	// os version
	osVersionByte, err := fileutil.ReadFromFile("/proc/version")
	if err != nil {
		logger.AppLogger().Errorf("failed ReadFromFile /proc/version : %s", err)
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message:   fmt.Errorf("failed read /proc/version, %+v", c.Params).Error(),
			RequestId: random.GenUUID()})
		return
	}

	docker := dockerfacade.NewDockerFacade()
	docker.SetClientVersion(config.Config.Docker.APIVersion)
	// serviceDetail, err := docker.ListImages()
	serviceDetail, err := docker.ListContainers()
	if err != nil {
		c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message:   fmt.Errorf("failed ListImages, %+v", err).Error(),
			RequestId: random.GenUUID()})
	}

	serviceVersion := make([]*device.ServiceVersion, 0)

	for _, v := range serviceDetail {
		serviceName := ""

		sArray := strings.Split(v.Image, ":")
		if len(sArray) < 2 {
			continue
		}
		sArrayService := strings.Split(sArray[len(sArray)-2], "/")
		if len(sArrayService) < 1 {
			continue
		}

		if len(v.Names) > 0 {
			serviceName = v.Names[0]
		} else {
			serviceName = sArrayService[len(sArrayService)-1]
		}
		if strings.Index(serviceName, "/") == 0 {
			serviceName = serviceName[1:]
		}

		serviceVersion = append(serviceVersion, &device.ServiceVersion{Created: v.Created,
			ServiceName: serviceName,
			Version:     sArray[len(sArray)-1]})
	}

	// result
	versions := &device.BoxDeviceVersion{
		ServiceVersion: serviceVersion,
		ServiceDetail:  serviceDetail}
	versions.DeviceName = device.GetDeviceName()
	versions.DeviceNameEn = device.GetDeviceNameEn()

	versions.GenerationEn = device.GetGenerationEn()
	versions.GenerationEn = device.GetGenerationZh()

	versions.ProductModel = device.GetProductModel()
	versions.SnNumber = device.GetSnNumber()

	versions.SpaceVersion = boxVersion
	versions.OSVersion = cutOSVersion(string(osVersionByte))

	versions.DeviceLogoUrl = "" // TODO: logo link url. 不给图片的话客户端用内置的默认图。

	versions.DeviceAbility = device_ability.GetAbilityModel()

	rsp := &dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		Message:   "OK",
		RequestId: random.GenUUID(),
		Results:   versions}
	c.IndentedJSON(http.StatusOK, rsp)
}

// Info godoc
// @Summary get device info [for mirco service]
// @Description get device info
// @ID device.info
// @Tags device
// @Accept  plain
// @Produce  json
// @Success 200 {object} dto.BaseRspStr{results=device.StorageInfo} "code=0 success;"
// @Router /agent/v1/api/device/info [GET]
func Info(c *gin.Context) {

	var rsp *dto.BaseRspStr
	if device_ability.GetAbilityModel().InnerDiskSupport {

	} else {
		var all uint64
		var used uint64
		var free uint64
		var err error
		if device_ability.GetAbilityModel().RunInDocker {
			all, used, free, err = space.FolderUsage(config.SpaceMountPath)
		} else {
			all, used, free, err = space.AllPartsUsage()
		}
		if err != nil {
			logger.AppLogger().Warnf("/agent/v1/api/device/info, failed DiskUsage, err:%+v", err)
			c.JSON(http.StatusOK, dto.BaseRspStr{Code: dto.AgentCodeCallDiskUsageFailedStr,
				Message:   fmt.Errorf("failed DiskUsage, %+v", c.Params).Error(),
				RequestId: random.GenUUID()})
			return
		}

		rsp = &dto.BaseRspStr{Code: dto.AgentCodeOkStr,
			Message:   "OK",
			RequestId: random.GenUUID(),
			Results: &device.StorageInfo{
				Used:  int64(used),
				Free:  int64(free),
				Total: int64(all)}}
	}

	logger.AppLogger().Debugf("device info return rsp:%+v", rsp)
	c.IndentedJSON(http.StatusOK, rsp)
}

func DockerImages() (string, error) {
	params := []string{"images"}
	logger.AppLogger().Debugf("DockerImages, run cmd: docker %v", strings.Join(params, " "))
	stdOutput, errOutput, err := run.RunExe("docker", params)
	if err != nil {
		return "", fmt.Errorf("failed run ConnectWifi %v, err is :%v, stdOutput is :%v, errOutput is :%v",
			params, err, string(stdOutput), string(errOutput))
	}
	logger.AppLogger().Debugf("DockerImages, run cmd: docker %v, stdOutput is :%v, errOutput is :%v",
		strings.Join(params, " "), string(stdOutput), string(errOutput))

	return string(stdOutput), nil
}

func cutOSVersion(rawVersion string) string {
	abilityModel := device_ability.GetAbilityModel()
	if abilityModel.RunInDocker {
		prefix := "version "
		suffix := " ("
		i1 := strings.Index(rawVersion, prefix)
		if i1 < 1 {
			return ""
		}
		i2 := strings.Index(rawVersion[i1+len(prefix):], suffix)
		if i2 < 0 {
			return rawVersion[i1+len(prefix):]
		}
		return rawVersion[i1+len(prefix):][:i2]
	} else {
		prefix := "version "
		suffix := ".aarch64"
		r, err := regexp.Compile(prefix + ".*" + suffix)
		if err != nil {
			fmt.Printf("failed regexp.Compile r, err:%v\n", err)
			return ""
		}
		version := r.FindString(rawVersion)
		version = strings.ReplaceAll(version, prefix, "")
		version = strings.ReplaceAll(version, suffix, "")
		// delete .raspi
		j := strings.LastIndex(version, ".")
		if j > 0 {
			version = version[:j]
		}
		return version
	}
}

func parseImagesVersion(s string) ([]device.ServiceVersion, error) {
	serviceVersion := make([]device.ServiceVersion, 0)

	zp := regexp.MustCompile(`[\t\n\f\r]`)
	lines := zp.Split(s, -1)

	for _, v := range lines {

		obj := device.ServiceVersion{}
		err := encoding.JsonDecode([]byte(v), &obj)
		if err != nil {
			return serviceVersion, fmt.Errorf("parseImagesVersion, failed JsonDecode, err:%v", err)
		} else {
			fmt.Printf("%+v\n\n", obj)
			serviceVersion = append(serviceVersion, obj)
		}
	}

	return serviceVersion, nil
}
