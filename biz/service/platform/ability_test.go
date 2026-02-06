package platform

import (
	"agent/config"
	"testing"
)

func TestCheckPlatformAbilityDisabled(t *testing.T) {
	origEnabled := config.Config.PlatformEnabled
	config.Config.PlatformEnabled = false
	defer func() { config.Config.PlatformEnabled = origEnabled }()

	if CheckPlatformAbility("/any") {
		t.Fatalf("expected false when platform disabled")
	}
}
