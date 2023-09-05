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
	"testing"
)

func TestUpContainersWithSample(t *testing.T) {
	d := DockerFacade{}
	_, stdErr, err := d.UpContainers("./test-docker-compose.yml")
	if err != nil {
		t.Fatal("exec err: ", err)
	} else if len(stdErr) != 0 {
		t.Fatal("stdErr: ", stdErr)
	}
}

func TestUpContainers(t *testing.T) {
	d := DockerFacade{}

	_, stdErr, err := d.UpContainers("../../../res/docker-compose.yml")
	if err != nil {
		t.Fatal("exec err: ", err)
	} else if len(stdErr) != 0 {
		t.Fatal("stdErr: ", stdErr)
	}
}
