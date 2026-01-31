package pair

import (
	"agent/config"
	"testing"
)

func TestServiceRegisterBoxWhenPlatformDisabled(t *testing.T) {
	origEnabled := config.Config.PlatformEnabled
	config.Config.PlatformEnabled = false
	defer func() { config.Config.PlatformEnabled = origEnabled }()

	if err := ServiceRegisterBox(); err != nil {
		t.Fatalf("expected no error when platform disabled, got %v", err)
	}
}

func TestGetDeviceRegKeyWhenPlatformDisabled(t *testing.T) {
	origEnabled := config.Config.PlatformEnabled
	config.Config.PlatformEnabled = false
	defer func() { config.Config.PlatformEnabled = origEnabled }()

	if _, err := GetDeviceRegKey(""); err == nil {
		t.Fatalf("expected error when platform disabled")
	}
}
