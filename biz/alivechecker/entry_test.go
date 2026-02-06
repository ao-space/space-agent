package alivechecker

import (
	"agent/biz/alivechecker/model"
	"agent/config"
	"testing"
)

func TestTestNetworkUsesLatestResults(t *testing.T) {
	origPing := pingFn
	origCurl := curlFn
	origCurlHeader := curlHeaderFn
	origGetDomain := getAdminDomainFn
	origPlatformEnabled := config.Config.PlatformEnabled

	origCloudHost := config.Config.NetworkCheck.CloudHost.Url
	origThirdParty := config.Config.NetworkCheck.ThirdPartyHost.Url
	origCloudIpv4 := config.Config.NetworkCheck.CloudIpv4.Url
	origCloudStatus := config.Config.NetworkCheck.CloudStatusHost.Url
	origCloudStatusIpv4 := config.Config.NetworkCheck.CloudStatusIpv4.Url
	origBoxStatus := config.Config.NetworkCheck.BoxStatusPath.Url

	t.Cleanup(func() {
		pingFn = origPing
		curlFn = origCurl
		curlHeaderFn = origCurlHeader
		getAdminDomainFn = origGetDomain
		config.Config.PlatformEnabled = origPlatformEnabled

		config.Config.NetworkCheck.CloudHost.Url = origCloudHost
		config.Config.NetworkCheck.ThirdPartyHost.Url = origThirdParty
		config.Config.NetworkCheck.CloudIpv4.Url = origCloudIpv4
		config.Config.NetworkCheck.CloudStatusHost.Url = origCloudStatus
		config.Config.NetworkCheck.CloudStatusIpv4.Url = origCloudStatusIpv4
		config.Config.NetworkCheck.BoxStatusPath.Url = origBoxStatus
	})

	// set deterministic hosts for assertions
	config.Config.PlatformEnabled = true
	config.Config.NetworkCheck.CloudHost.Url = "cloud-host"
	config.Config.NetworkCheck.ThirdPartyHost.Url = "third-party"
	config.Config.NetworkCheck.CloudIpv4.Url = "cloud-ipv4"
	config.Config.NetworkCheck.CloudStatusHost.Url = "cloud-status"
	config.Config.NetworkCheck.CloudStatusIpv4.Url = "cloud-status-ipv4"
	config.Config.NetworkCheck.BoxStatusPath.Url = "box/status"

	pingFn = func(host string) (bool, error) {
		switch host {
		case "cloud-host":
			return true, nil
		case "third-party":
			return false, nil
		case "cloud-ipv4":
			return true, nil
		default:
			return false, nil
		}
	}

	curlCalls := map[string]int{}
	curlFn = func(host string) (bool, error) {
		curlCalls[host]++
		switch host {
		case "cloud-status":
			return false, nil
		case "admin-domain/box/status":
			return true, nil
		default:
			return true, nil
		}
	}
	curlHeaderFn = func(host string) (bool, error) {
		if host == "cloud-status-ipv4" {
			return true, nil
		}
		return false, nil
	}
	getAdminDomainFn = func() string {
		return "admin-domain"
	}

	TestNetwork()

	got := model.Get()
	if got.PingCloudHost != true {
		t.Fatalf("PingCloudHost expected true, got %v", got.PingCloudHost)
	}
	if got.PingThirdPartyHost != false {
		t.Fatalf("PingThirdPartyHost expected false, got %v", got.PingThirdPartyHost)
	}
	if got.PingCloudIpv4 != true {
		t.Fatalf("PingCloudIpv4 expected true, got %v", got.PingCloudIpv4)
	}
	if got.CurlCloudStatusHost != false {
		t.Fatalf("CurlCloudStatusHost expected false, got %v", got.CurlCloudStatusHost)
	}
	if got.CurlHttpHeaderCloudStatusIpv4 != true {
		t.Fatalf("CurlHttpHeaderCloudStatusIpv4 expected true, got %v", got.CurlHttpHeaderCloudStatusIpv4)
	}

	if curlCalls["admin-domain/box/status"] != 1 {
		t.Fatalf("expected admin-domain curl to be called once, got %d", curlCalls["admin-domain/box/status"])
	}
}
