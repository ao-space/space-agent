package bindinit

import (
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	dtopair "agent/biz/model/dto/pair"
	servicesinit "agent/biz/service/bind/init"
	"agent/biz/model/dto"
	"agent/config"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

type baseRsp struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Results json.RawMessage `json:"results"`
}

func TestBindInitHandlerBindsQueryAndStoresClientVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	restoreNetwork := servicesinit.SetGetConnectedNetworkFuncForTest(func() []*dtopair.Network {
		return []*dtopair.Network{}
	})
	restoreVersion := servicesinit.SetGetInstalledVersionFuncForTest(func() string {
		return "1.2.3"
	})
	defer restoreNetwork()
	defer restoreVersion()

	origEncrypt := config.Config.EncryptLanSessionData
	config.Config.EncryptLanSessionData = false
	defer func() { config.Config.EncryptLanSessionData = origEncrypt }()

	device.GetDeviceInfo().BoxUuid = "box-uuid"
	device_ability.GetAbilityModel().DeviceModelNumber = device_ability.SN_GEN_2

	router := gin.New()
	router.GET("/agent/v1/api/bind/init", Init)

	q := url.Values{}
	q.Set("clientUuid", "client-uuid-1")
	q.Set("clientVersion", "1.0.0")
	req := httptest.NewRequest("GET", "/agent/v1/api/bind/init?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("unexpected status code: %d", w.Code)
	}

	var rsp baseRsp
	if err := json.Unmarshal(w.Body.Bytes(), &rsp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if rsp.Code != dto.AgentCodeOkStr {
		t.Fatalf("unexpected response code: %s", rsp.Code)
	}

	v, found := clientinfo.GetClientVersion("client-uuid-1")
	if !found {
		t.Fatalf("client version not stored")
	}
	if v != "1.0.0" {
		t.Fatalf("client version mismatch: %s", v)
	}
}
