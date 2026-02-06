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

package dockerfacade

import (
	"os"
	"os/exec"
	"testing"
)

func TestUpContainersWithSample(t *testing.T) {
	if os.Getenv("ENABLE_DOCKER_TESTS") != "1" {
		// TODO: Provide a Docker API stub or harness so this test can run in CI without Docker.
		t.Skip("ENABLE_DOCKER_TESTS is not set")
	}
	if _, err := exec.LookPath("docker-compose"); err != nil {
		// TODO: Replace docker-compose dependency with a mock to make this deterministic.
		t.Skip("docker-compose not found")
	}
	d := DockerFacade{}
	_, stdErr, err := d.UpContainers("./test-docker-compose.yml", nil)
	if err != nil {
		t.Fatal("exec err: ", err)
	} else if len(stdErr) != 0 {
		t.Fatal("stdErr: ", stdErr)
	}
}

func TestUpContainers(t *testing.T) {
	if os.Getenv("ENABLE_DOCKER_TESTS") != "1" {
		// TODO: Provide a Docker API stub or harness so this test can run in CI without Docker.
		t.Skip("ENABLE_DOCKER_TESTS is not set")
	}
	if _, err := exec.LookPath("docker-compose"); err != nil {
		// TODO: Replace docker-compose dependency with a mock to make this deterministic.
		t.Skip("docker-compose not found")
	}
	d := DockerFacade{}

	_, stdErr, err := d.UpContainers("../../../res/docker-compose.yml", nil)
	if err != nil {
		t.Fatal("exec err: ", err)
	} else if len(stdErr) != 0 {
		t.Fatal("stdErr: ", stdErr)
	}
}
