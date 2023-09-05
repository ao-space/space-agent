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

package device

import (
	"agent/biz/model/device_ability"
	"agent/config"
	"agent/utils/deviceid"
	"agent/utils/hardware"
	"agent/utils/version"
	"fmt"
	"net/url"
	"strings"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

const (
	UrlQrCodeDomain = "https://ao.space/"
)

type NetworkClientInfo struct {
	ClientID  string `json:"clientId"`
	SecretKey string `json:"secretKey"`
}

type DeviceInfo struct {
	BoxUuid   string `json:"boxUuid"`   // 盒子端唯一id.
	Btid      string `json:"btid"`      // 二维码中的蓝牙id.
	BtidHash  string `json:"btidHash"`  // mdns广播使用的id.
	BoxQrCode string `json:"boxQrCode"` // 盒子的二维码数据.

	BoxRegKey         string `json:"boxRegKey"` // 平台端授权给盒子端的凭证. 空间平台使用，过期则重新申请
	BoxRegKeyExpireAt string `json:"expiresAt"`

	ApiBaseUrl    string             `json:"apiBaseUrl"` //空间平台的基础部分
	NetworkClient *NetworkClientInfo `json:"networkClient"`

	BoxRegisterTime int64 `json:"boxRegisterTime"`
	IsBoxRegistered bool  `json:"isBoxRegistered"`
}

// 考虑到配对并发的可能性很小，故暂时不考虑 race-condition
var deviceInfo DeviceInfo

func InitDeviceInfo() {
	deviceInfo.BoxRegisterTime = -1
	deviceInfo.IsBoxRegistered = false
	logger.AppLogger().Infof("InitBoxInfo InitBoxInfo, boxInfo.NetworkClient:%+v", deviceInfo.NetworkClient)

	b, err := readDeviceInfoFromFile()
	if err == nil {
		logger.AppLogger().Debugf("InitdeviceInfo, deviceInfoFile read succ.")
		deviceInfo = b
		writeSharedInfoFile()
		//UpdateMdns()
	} else {
		logger.AppLogger().Debugf("InitdeviceInfo, read error. write new to file")
		err1 := setDeviceInitData()
		if err1 == nil {
			writeDeviceInfoToFile(deviceInfo)
			writeSharedInfoFile()
			//UpdateMdns()
		}
	}

	if deviceInfo.NetworkClient == nil {
		deviceInfo.NetworkClient = &NetworkClientInfo{}
	}
	logger.AppLogger().Infof("Leave InitdeviceInfo, deviceInfo:%+v", deviceInfo)
}

func (b *DeviceInfo) SetAllInfo() error {
	// 方式2: 根据硬件生成
	btId, err := deviceid.GetBtId(config.Config.Box.CpuIdStoreFile, config.Config.Box.SnNumberStoreFile)
	if err != nil {
		err1 := fmt.Errorf("GET GetRPIBtId failed, no btid now, err:%v", err)
		logger.AppLogger().Warnf("%v", err1)
	}
	boxUuid, err := deviceid.GetProductId(config.Config.Box.CpuIdStoreFile)
	if err != nil {
		err1 := fmt.Errorf("GET GetRPIBtId failed, no boxUuid now, err:%v", err)
		logger.AppLogger().Errorf("%v", err1)
		return err1
	}
	btidHash := deviceid.GetBtIdHash(btId, config.Config.Box.Avahi.BtIdHashPrefix)
	qrCode := GetQrCode()
	logger.AppLogger().Debugf("setBoxInitData, qrCode:%v", qrCode)

	if len(b.BoxUuid) < 1 {
		b.BoxUuid = boxUuid
	}
	if len(b.Btid) < 1 {
		b.Btid = btId
	}
	if len(b.BtidHash) < 1 {
		b.BtidHash = btidHash
	}
	b.BoxQrCode = qrCode
	return nil
}

func setDeviceInitData() error {

	// 方式1: 随机生成
	// boxUuid := random.GenUUID()

	// 方式2: 根据硬件生成
	btId, err := deviceid.GetBtId(config.Config.Box.CpuIdStoreFile, config.Config.Box.SnNumberStoreFile)
	if err != nil {
		err1 := fmt.Errorf("GET GetRPIBtId failed, no btid now, err:%v", err)
		logger.AppLogger().Warnf("%v", err1)
	}
	boxUuid, err := deviceid.GetProductId(config.Config.Box.CpuIdStoreFile)
	if err != nil {
		err1 := fmt.Errorf("GET GetRPIBtId failed, no boxUuid now, err:%v", err)
		logger.AppLogger().Errorf("%v", err1)
		return err1
	}
	btidHash := deviceid.GetBtIdHash(btId, config.Config.Box.Avahi.BtIdHashPrefix)
	qrCode := GetQrCode()
	logger.AppLogger().Debugf("setBoxInitData, qrCode:%v", qrCode)

	if len(deviceInfo.BoxUuid) < 1 {
		deviceInfo.BoxUuid = boxUuid
	}
	if len(deviceInfo.Btid) < 1 {
		deviceInfo.Btid = btId
	}
	if len(deviceInfo.BtidHash) < 1 {
		deviceInfo.BtidHash = btidHash
	}
	deviceInfo.BoxQrCode = qrCode
	return nil
}

func UpdateSnNumber(snNumber string) error {
	err := fileutil.WriteToFile(config.Config.Box.SnNumberStoreFile, []byte(snNumber), true)
	if err != nil {
		logger.AppLogger().Errorf("Write deviceInfo SnNumberStoreFile failed, file:%v, err:%v", config.Config.Box.SnNumberStoreFile, err)
	} else {
		logger.AppLogger().Debugf("Write deviceInfo SnNumberStoreFile succ, file:%v, b:%+v", config.Config.Box.SnNumberStoreFile, snNumber)
	}
	err = setDeviceInitData()
	if err == nil {
		writeDeviceInfoToFile(deviceInfo)
		writeSharedInfoFile()
		//UpdateMdns()
	}
	return err
}

func UpdateApplyEmail(applyEmail string) error {
	err := fileutil.WriteToFile(config.Config.Box.ApplyEmailStoreFile, []byte(applyEmail), true)
	if err != nil {
		logger.AppLogger().Errorf("Write ApplyEmail failed, file:%v, err:%v", config.Config.Box.ApplyEmailStoreFile, err)
	} else {
		logger.AppLogger().Debugf("Write ApplyEmail succ, file:%v, applyEmail:%+v", config.Config.Box.ApplyEmailStoreFile, applyEmail)
	}
	return err
}

func GetApplyEmail() (string, error) {
	applyEmail, err := fileutil.ReadFromFile(config.Config.Box.ApplyEmailStoreFile)
	if err != nil {
		logger.AppLogger().Warnf("ReadFromFile ApplyEmail failed, file:%v, err:%v", config.Config.Box.ApplyEmailStoreFile, err)
	} else {
		logger.AppLogger().Debugf("ReadFromFile ApplyEmail succ, file:%v, applyEmail:%+v", config.Config.Box.ApplyEmailStoreFile, applyEmail)
	}
	return string(applyEmail), err
}

func GetQrCode() string {
	logger.AppLogger().Debugf("GetQrCode, device_ability.GetAbilityModel().DeviceModelNumber:%v", device_ability.GetAbilityModel().DeviceModelNumber)
	if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_SUPPORTED_FROM_MODEL_NUMBER {
		snNumber, err := deviceid.GetSnNumber(config.Config.Box.SnNumberStoreFile)
		if err != nil {
			err1 := fmt.Errorf("failed GetSnNumber, err:%v", err)
			// logger.AppLogger().Debugf("%v", err1)
			return err1.Error()
		}
		return fmt.Sprintf("%v?sn=%v", UrlQrCodeDomain, snNumber)
	} else if device_ability.GetAbilityModel().DeviceModelNumber <= device_ability.SN_GEN_CLOUD_DOCKER {
		snNumber, err := deviceid.GetSnNumber(config.Config.Box.SnNumberStoreFile)
		if err != nil {
			err1 := fmt.Errorf("failed GetSnNumber, err:%v", err)
			// logger.AppLogger().Debugf("%v", err1)
			return err1.Error()
		}
		logger.AppLogger().Debugf("GetQrCode, snNumber:%v", snNumber)

		ipaddr, err := fileutil.ReadFromFile(config.Config.Box.HostIpFile)
		if err != nil {
			err1 := fmt.Errorf("failed ReadFromFile HostIpFile, err:%v", err)
			logger.AppLogger().Errorf("%v", err1)
			// https://ao.space/?btid=4650ce2ef8b85d1704a73f9690995c9d47e&port=5678
			return fmt.Sprintf("%v?sn=%v&port=%v", UrlQrCodeDomain, snNumber, strings.TrimLeft(config.Config.Web.DefaultListenAddr, ":"))
		} else {
			// https://ao.space/?btid=4650ce2ef8b85d1704a73f9690995c9d47e&ipaddr=urlencode(192.168.124.100)&port=5678
			ip := string(ipaddr)
			port := "80"
			if strings.Contains(string(ipaddr), ":") {
				logger.AppLogger().Debugf("GetQrCode, string(ipaddr):%v", string(ipaddr))
				arr := strings.Split(string(ipaddr), ":")
				if len(arr) > 0 {
					ip = arr[0]
				}
				if len(arr) > 1 {
					port = arr[1]
				}
			}
			logger.AppLogger().Debugf("GetQrCode, ip:%v, port:%v", ip, port)
			// return fmt.Sprintf("%v?sn=%v"+`&`+"ipaddr=%v"+`&`+"port=%v", UrlQrCodeDomain, snNumber,
			// 	url.QueryEscape(s),
			// 	strings.TrimLeft(config.Config.Web.DefaultListenAddr, ":"))
			return UrlQrCodeDomain + `?sn=` + snNumber + `&ipaddr=` + url.QueryEscape(ip) + `&port=` + port
		}

	} else {
		btId, err := deviceid.GetBtId(config.Config.Box.CpuIdStoreFile, config.Config.Box.SnNumberStoreFile)
		if err != nil {
			err1 := fmt.Errorf("failed GetRPIBtId, err:%v", err)
			logger.AppLogger().Errorf("%v", err1)
			return err1.Error()
		}
		return fmt.Sprintf("%v?btid=%v", UrlQrCodeDomain, btId)
	}
}

func GetDeviceInfo() *DeviceInfo {
	// logger.AppLogger().Debugf("GetdeviceInfo, deviceInfo:%+v", deviceInfo)
	return &deviceInfo
}

func (b *DeviceInfo) SetRegKey(boxRegKey string, expiresAt string) {
	b.BoxRegKey = boxRegKey
	b.BoxRegKeyExpireAt = expiresAt
}

func (b *DeviceInfo) SetApiBaseUrl() {
	if len(deviceInfo.ApiBaseUrl) == 0 {
		b.ApiBaseUrl = config.Config.Platform.APIBase.Url
	} else {
		b.ApiBaseUrl = deviceInfo.ApiBaseUrl
	}
}

func (b *DeviceInfo) SetGT() {
	//b.NetworkClient = networkClient
}

func SetDeviceRegKey(boxRegKey string, expiresAt string) {
	deviceInfo.BoxRegKey = boxRegKey
	deviceInfo.BoxRegKeyExpireAt = expiresAt
	logger.AppLogger().Debugf("SetBoxRegKey, deviceInfo.BoxRegKey:%+v, deviceInfo.BoxRegKeyExpireAt:%+v ", deviceInfo.BoxRegKey, deviceInfo.BoxRegKeyExpireAt)
	writeDeviceInfoToFile(deviceInfo)
}

func SetApiBaseUrl(url string) {
	deviceInfo.ApiBaseUrl = url
	logger.AppLogger().Debugf("SetNetworkClient, deviceInfo.ApiBaseUrl=%+v", deviceInfo.ApiBaseUrl)
	writeDeviceInfoToFile(deviceInfo)
}

func GetApiBaseUrl() string {
	if len(deviceInfo.ApiBaseUrl) == 0 {
		return config.Config.Platform.APIBase.Url
	} else {
		return deviceInfo.ApiBaseUrl
	}
}

func SetNetworkClient(networkClient *NetworkClientInfo) {
	deviceInfo.NetworkClient = networkClient
	logger.AppLogger().Debugf("SetNetworkClient, deviceInfo.NetworkClient=%+v", deviceInfo.NetworkClient)
	writeDeviceInfoToFile(deviceInfo)
}

func (b *DeviceInfo) Registered() {
	b.BoxRegisterTime = time.Now().Unix()
	b.IsBoxRegistered = true
}

func (b *DeviceInfo) Unregistered() {
	b.BoxRegisterTime = 0
	b.IsBoxRegistered = false
}

func SetBoxRegistered() {
	deviceInfo.BoxRegisterTime = time.Now().Unix()
	deviceInfo.IsBoxRegistered = true
	logger.AppLogger().Debugf("SetBoxRegistered, deviceInfo=%+v", deviceInfo)
	writeDeviceInfoToFile(deviceInfo)
}
func SetDeviceUnregistered() {
	deviceInfo.BoxRegisterTime = 0
	deviceInfo.IsBoxRegistered = false
	logger.AppLogger().Debugf("SetBoxUnregistered, deviceInfo.BoxRegKey:%+v", deviceInfo.BoxRegKey)
	writeDeviceInfoToFile(deviceInfo)
}

func (b *DeviceInfo) Load() (*DeviceInfo, error) {

	err := fileutil.ReadFileJsonToObject(config.Config.Box.BoxInfoFile, b)
	if err != nil {
		logger.AppLogger().Errorf("Read BoxInfo file failed, file:%v, err:%v", config.Config.Box.BoxInfoFile, err)
		return b, err
	}
	logger.AppLogger().Debugf("readBoxInfoFromFile, binfo:%+v", b)

	if len(b.BoxUuid) < 1 || len(b.Btid) < 1 || len(b.BtidHash) < 1 || len(b.BoxQrCode) < 1 {
		logger.AppLogger().Errorf("readBoxInfoFromFile, binfo:%+v empty", b)
		return b, fmt.Errorf("read device fields empty from %v", config.Config.Box.BoxInfoFile)
	}

	// 根据硬件生成, 看看与本地缓存的是否一致
	btId, err := deviceid.GetBtId(config.Config.Box.CpuIdStoreFile, config.Config.Box.SnNumberStoreFile)
	if err != nil {
		logger.AppLogger().Errorf("GET GetBtId failed, no btid now, err:%v", err)
		return b, err
	}
	boxUuid, err := deviceid.GetProductId(config.Config.Box.CpuIdStoreFile)
	if err != nil {
		logger.AppLogger().Errorf("get GetProductId failed, no boxUuid now , err:%v", err)
		return b, err
	}

	if b.Btid != btId || b.BoxUuid != boxUuid {
		logger.AppLogger().Warnf("readBoxInfoFromFile, b.Btid{%v}!=btId{%v} || b.BoxUuid{%v}!=boxUuid{%v}",
			b.Btid, btId, b.BoxUuid, boxUuid)
		return b, fmt.Errorf("read device, boxuuid has changed")
	}

	return b, nil
}

func (b *DeviceInfo) Save() error {
	err := fileutil.WriteToFileAsJson(config.Config.Box.BoxInfoFile, b, "  ", true)
	if err != nil {
		logger.AppLogger().Errorf("Write BoxInfo file failed, file:%v, err:%v", config.Config.Box.BoxInfoFile, err)
	} else {
		logger.AppLogger().Debugf("Write BoxInfo file succ, file:%v, b:%+v", config.Config.Box.BoxInfoFile, b)
	}
	logger.AppLogger().Debugf("Write BoxInfo file succ")

	return err
}

func readDeviceInfoFromFile() (DeviceInfo, error) {
	var b DeviceInfo
	err := fileutil.ReadFileJsonToObject(config.Config.Box.BoxInfoFile, &b)
	if err != nil {
		logger.AppLogger().Errorf("Read BoxInfo file failed, file:%v, err:%v", config.Config.Box.BoxInfoFile, err)
		return b, err
	}
	logger.AppLogger().Debugf("readBoxInfoFromFile, binfo:%+v", b)

	if len(b.BoxUuid) < 1 || len(b.Btid) < 1 || len(b.BtidHash) < 1 || len(b.BoxQrCode) < 1 {
		logger.AppLogger().Errorf("readBoxInfoFromFile, binfo:%+v empty", b)
		return b, fmt.Errorf("read device fields empty from %v", config.Config.Box.BoxInfoFile)
	}

	// 根据硬件生成, 看看与本地缓存的是否一致
	btId, err := deviceid.GetBtId(config.Config.Box.CpuIdStoreFile, config.Config.Box.SnNumberStoreFile)
	if err != nil {
		logger.AppLogger().Errorf("GET GetBtId failed, no btid now, err:%v", err)
		return b, err
	}
	boxUuid, err := deviceid.GetProductId(config.Config.Box.CpuIdStoreFile)
	if err != nil {
		logger.AppLogger().Errorf("get GetProductId failed, no boxUuid now , err:%v", err)
		return b, err
	}

	if b.Btid != btId || b.BoxUuid != boxUuid {
		logger.AppLogger().Warnf("readBoxInfoFromFile, b.Btid{%v}!=btId{%v} || b.BoxUuid{%v}!=boxUuid{%v}",
			b.Btid, btId, b.BoxUuid, boxUuid)
		return b, fmt.Errorf("read device, boxuuid has changed")
	}

	return b, nil
}

func writeDeviceInfoToFile(b DeviceInfo) error {
	err := fileutil.WriteToFileAsJson(config.Config.Box.BoxInfoFile, b, "  ", true)
	if err != nil {
		logger.AppLogger().Errorf("Write BoxInfo file failed, file:%v, err:%v", config.Config.Box.BoxInfoFile, err)
	} else {
		logger.AppLogger().Debugf("Write BoxInfo file succ, file:%v, b:%+v", config.Config.Box.BoxInfoFile, b)
	}
	logger.AppLogger().Debugf("Write BoxInfo file succ")

	return err
}

func writeSharedInfoFile() {
	type SharedBoxInfo struct {
		BoxUuid    string `json:"boxUuid"`    // 盒子端唯一id.
		Btid       string `json:"btid"`       // 盒子蓝牙id.
		BoxVersion string `json:"boxVersion"` // 盒子版本号
	}

	boxVersion := version.GetInstalledAgentVersionRemovedNewLine()
	if hardware.RunningInDocker() {
		boxVersion = config.VersionNumber
	}
	s := &SharedBoxInfo{BoxUuid: deviceInfo.BoxUuid, Btid: deviceInfo.Btid, BoxVersion: boxVersion}
	err := fileutil.WriteToFileAsJson(config.Config.Box.PublicSharedInfoFile, s, "  ", true)
	if err != nil {
		logger.AppLogger().Errorf("Write SharedBoxInfo file failed, file:%v, err:%v", config.Config.Box.PublicSharedInfoFile, err)
	}
	logger.AppLogger().Debugf("Write SharedBoxInfo file succ")
}

func IsBoxRegistered() bool {
	return deviceInfo.IsBoxRegistered
}
