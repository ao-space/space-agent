package routers

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func routeSet(routes []gin.RouteInfo) map[string]struct{} {
	out := make(map[string]struct{}, len(routes))
	for _, route := range routes {
		out[fmt.Sprintf("%s %s", route.Method, route.Path)] = struct{}{}
	}
	return out
}

func mustContainAll(t *testing.T, registered map[string]struct{}, expected []string) {
	t.Helper()
	for _, api := range expected {
		if _, ok := registered[api]; !ok {
			t.Fatalf("missing route: %s", api)
		}
	}
}

func TestExternalRouter_RegistersExpectedAPIs(t *testing.T) {
	registered := routeSet(ExternalRouter().Routes())
	expected := []string{
		"POST /agent/v1/api/initial",
		"POST /agent/v1/api/pairing",
		"POST /agent/v1/api/pubkeyexchange",
		"POST /agent/v1/api/keyexchange",
		"POST /agent/v1/api/setpassword",
		"GET /agent/v1/api/device/ability",
		"POST /agent/v1/api/admin/revoke",
		"GET /agent/v1/api/pair/net/localips",
		"GET /agent/v1/api/pair/net/netconfig",
		"GET /agent/v1/api/pair/init",
		"GET /agent/v1/api/bind/init",
		"POST /agent/v1/api/bind/com/start",
		"GET /agent/v1/api/bind/com/progress",
		"POST /agent/v1/api/bind/space/create",
		"POST /agent/v1/api/bind/internet/service/config",
		"GET /agent/v1/api/bind/internet/service/config",
		"POST /agent/v1/api/bind/password/verify",
		"POST /agent/v1/api/bind/revoke",
		"GET /agent/v1/api/space/ready/check",
		"POST /agent/v1/api/network/config",
		"GET /agent/v1/api/network/config",
		"POST /agent/v1/api/network/ignore",
		"POST /agent/v1/api/passthrough",
		"POST /agent/v1/api/switch",
		"GET /agent/v1/api/switch/status",
		"GET /agent/v1/cert/get",
		"GET /agent/v1/api/did/document",
		"PUT /agent/v1/api/did/document/password",
		"PUT /agent/v1/api/did/document/method",
		"GET /agent/status",
		"GET /agent/info",
	}
	mustContainAll(t, registered, expected)
}

func TestInternalRouter_RegistersExpectedAPIs(t *testing.T) {
	registered := routeSet(InternalRouter().Routes())
	expected := []string{
		"GET /agent/v1/api/device/info",
		"GET /agent/v1/api/device/version",
		"GET /agent/v1/api/device/ability",
		"GET /agent/v1/api/device/localips",
		"GET /agent/v1/api/device/netconfig",
		"GET /agent/v1/api/upgrade/config",
		"POST /agent/v1/api/upgrade/config",
		"POST /agent/v1/api/upgrade/download",
		"POST /agent/v1/api/upgrade/install",
		"GET /agent/v1/api/upgrade/status",
		"POST /agent/v1/api/network/config",
		"GET /agent/v1/api/network/config",
		"POST /agent/v1/api/network/ignore",
		"POST /agent/v1/api/system/shutdown",
		"POST /agent/v1/api/system/reboot",
		"GET /agent/v1/api/cert/get",
		"POST /agent/v1/api/bind/internet/service/config",
		"GET /agent/v1/api/bind/internet/service/config",
		"GET /agent/v1/api/did/document",
		"PUT /agent/v1/api/did/document/password",
		"PUT /agent/v1/api/did/document/method",
	}
	mustContainAll(t, registered, expected)
}
