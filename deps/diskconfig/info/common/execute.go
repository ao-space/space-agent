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

package common

import (
	"bytes"
	"os/exec"
)

func Execute(exe string, params ...string) (string, string, error) {
	cmd := exec.Command(exe, params...)

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Start()
	if err != nil {
		return "", "", err
	}

	if err := cmd.Wait(); err != nil {
		// fmt.Println("Wait: ", err.Error())
		return "", "", err
	}

	return stdOut.String(), stdErr.String(), nil
}
