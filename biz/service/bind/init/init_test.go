package init

import (
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	bindinitdto "agent/biz/model/dto/bind/bindinit"
	dtopair "agent/biz/model/dto/pair"
	"agent/config"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitServiceProcessSetsFieldsAndClientVersion(t *testing.T) {
	restoreNetwork := SetGetConnectedNetworkFuncForTest(func() []*dtopair.Network {
		return []*dtopair.Network{{Ip: "1.2.3.4"}}
	})
	restoreVersion := SetGetInstalledVersionFuncForTest(func() string {
		return "1.2.3"
	})
	defer restoreNetwork()
	defer restoreVersion()

	origEncrypt := config.Config.EncryptLanSessionData
	config.Config.EncryptLanSessionData = false
	defer func() { config.Config.EncryptLanSessionData = origEncrypt }()

	origAdminPairFile := config.Config.Box.BoxMetaAdminPair
	adminPairFile := filepath.Join(t.TempDir(), "admin.json")
	config.Config.Box.BoxMetaAdminPair = adminPairFile
	defer func() { config.Config.Box.BoxMetaAdminPair = origAdminPairFile }()

	adminInfo := map[string]string{
		"clientUUID": "client-uuid-2",
		"status":     "0",
		"boxName":    "Box-A",
	}
	data, err := json.Marshal(adminInfo)
	if err != nil {
		t.Fatalf("failed to marshal admin info: %v", err)
	}
	if err := os.WriteFile(adminPairFile, data, 0o644); err != nil {
		t.Fatalf("failed to write admin info: %v", err)
	}

	device.GetDeviceInfo().BoxUuid = "box-uuid"
	ability := device_ability.GetAbilityModel()
	ability.DeviceModelNumber = device_ability.SN_GEN_2
	ability.RunInDocker = false

	svc := &InitService{}
	req := &bindinitdto.InitReq{ClientUuid: "client-uuid-2", ClientVersion: "2.0.0"}
	svc.Req = req
	rsp := svc.Process()

	if rsp.Code != dto.AgentCodeOkStr {
		t.Fatalf("unexpected response code: %s", rsp.Code)
	}

	initRsp, ok := svc.Rsp.(*dtopair.InitResult)
	if !ok {
		t.Fatalf("unexpected response type: %T", svc.Rsp)
	}

	if initRsp.BoxUuid != "box-uuid" || initRsp.ProductId != "box-uuid" {
		t.Fatalf("unexpected box uuid/product id: %s/%s", initRsp.BoxUuid, initRsp.ProductId)
	}
	if initRsp.Paired != 0 || !initRsp.PairedBool {
		t.Fatalf("unexpected paired status: %d/%v", initRsp.Paired, initRsp.PairedBool)
	}
	if initRsp.ClientUuid != "client-uuid-2" || initRsp.BoxName != "Box-A" {
		t.Fatalf("unexpected client/box info: %s/%s", initRsp.ClientUuid, initRsp.BoxName)
	}
	if initRsp.SpaceVersion != "1.2.3" {
		t.Fatalf("unexpected space version: %s", initRsp.SpaceVersion)
	}
	if len(initRsp.Networks) != 1 || initRsp.Networks[0].Ip != "1.2.3.4" {
		t.Fatalf("unexpected networks: %+v", initRsp.Networks)
	}

	v, found := clientinfo.GetClientVersion("client-uuid-2")
	if !found || v != "2.0.0" {
		t.Fatalf("client version not stored: %v/%s", found, v)
	}
}
