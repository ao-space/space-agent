package config

import (
	"agent/biz/model/device"
	"agent/biz/model/dto"
	configdto "agent/biz/model/dto/bind/internet/service/config"
	"agent/config"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInternetServiceConfigPlatformDisabled(t *testing.T) {
	origEnabled := config.Config.PlatformEnabled
	config.Config.PlatformEnabled = false
	defer func() { config.Config.PlatformEnabled = origEnabled }()

	origAdminPairFile := config.Config.Box.BoxMetaAdminPair
	adminPairFile := filepath.Join(t.TempDir(), "admin.json")
	config.Config.Box.BoxMetaAdminPair = adminPairFile
	defer func() { config.Config.Box.BoxMetaAdminPair = origAdminPairFile }()

	adminInfo := map[string]string{
		"clientUUID": "client-uuid",
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

	origInternetConfigFile := config.Config.Box.InternetServiceConfigFile
	internetConfigFile := filepath.Join(t.TempDir(), "internet_config.json")
	config.Config.Box.InternetServiceConfigFile = internetConfigFile
	defer func() { config.Config.Box.InternetServiceConfigFile = origInternetConfigFile }()
	device.SetConfig(&device.InternetServiceConfig{EnableInternetAccess: false})

	svc := &InternetServiceConfig{}
	req := &configdto.ConfigReq{EnableInternetAccess: true}
	svc.Req = req
	rsp := svc.Process()
	if rsp.Code != dto.AgentCodeUnsupportedFunction {
		t.Fatalf("unexpected code: %s", rsp.Code)
	}
}
