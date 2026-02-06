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

package upgrade

import (
	"agent/config"
	"fmt"
	"os"
	"testing"
)

func TestCheckLatestVersion(t *testing.T) {
	if os.Getenv("ENABLE_PLATFORM_TESTS") != "1" {
		// TODO: Add a platform API mock (or fixture server) so this test runs in CI without external dependency.
		t.Skip("ENABLE_PLATFORM_TESTS is not set")
	}
	origEnabled := config.Config.PlatformEnabled
	config.Config.PlatformEnabled = true
	t.Cleanup(func() { config.Config.PlatformEnabled = origEnabled })
	versionDesc, err := CheckLatestVersion()
	if err != nil {
		t.Fatalf("failed to get latest version")
	}
	if versionDesc.PkgVersion != "" {
		fmt.Printf("version:%s", versionDesc.PkgVersion)
	}
}
