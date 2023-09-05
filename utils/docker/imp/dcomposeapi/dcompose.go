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

/*
 * @Author: jeffery
 * @Date: 2022-02-21 10:01:15
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-28 15:30:40
 * @Description:
 */
package dcomposeapi

import (
	"fmt"

	"github.com/dungeonsnd/gocom/sys/run"
)

const composeExe = "docker-compose"

func CreateNetwork(networkName string, stdOutput chan string, errOutput chan string) error {
	err := run.RunAndChangeDir(stdOutput, errOutput, "", "docker", "network", "create", networkName)
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir network create, err: %v", err)
	}
	return nil
}

func RemoveNetwork(networkName string, stdOutput chan string, errOutput chan string) error {
	err := run.RunAndChangeDir(stdOutput, errOutput, "", "docker", "network", "rm", networkName)
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir network rm, err: %v", err)
	}
	return nil
}

func PullImages(composeFile string, stdOutput chan string, errOutput chan string) error {
	err := run.RunAndChangeDir(stdOutput, errOutput, "", composeExe, "-f", composeFile, "pull")
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir docker-compose pull, err: %v", err)
	}
	return nil
}

func CreateContainers(composeFile string, stdOutput chan string, errOutput chan string) error {
	err := run.RunAndChangeDir(stdOutput, errOutput, "", composeExe, "-f", composeFile, "up", "--no-start")
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir docker-compose up, err: %v", err)
	}
	return nil
}

func StartContainers(composeFile string, stdOutput chan string, errOutput chan string) error {
	err := run.RunAndChangeDir(stdOutput, errOutput, "", composeExe, "-f", composeFile, "start")
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir docker-compose start, err: %v", err)
	}
	return nil
}

func DownContainers(composeFile string, stdOutput chan string, errOutput chan string) error {
	err := run.RunAndChangeDir(stdOutput, errOutput, "", composeExe, "-f", composeFile, "down")
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir docker-compose down, err: %v", err)
	}
	return nil
}

func StopSpecifiedContainers(composeFile string, stdOutput chan string, errOutput chan string, params ...string) error {
	for _, svc := range params {
		err := run.RunAndChangeDir(stdOutput, errOutput, "", composeExe, "-f", composeFile, "stop", svc)
		if err != nil {
			return fmt.Errorf("failed RunAndChangeDir docker-compose down, err: %v", err)
		}
	}
	return nil
}

func UpContainers(composeFile string, stdOutput chan string, errOutput chan string) error {
	err := run.RunAndChangeDir(stdOutput, errOutput, "", composeExe, "-f", composeFile, "up", "-d", "--remove-orphans")
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir docker-compose up, err: %v", err)
	}
	return nil
}

func UpContainersWithNoRecreate(composeFile string, stdOutput chan string, errOutput chan string) error {
	err := run.RunAndChangeDir(stdOutput, errOutput, "", composeExe, "-f", composeFile, "up", "-d", "--no-recreate", "--remove-orphans")
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir docker-compose up --no-recreate, err: %v", err)
	}
	return nil
}

func UpSpecifiedContainers(composeFile string, stdOutput chan string, errOutput chan string, params ...string) error {
	p := []string{"-f", composeFile, "up", "-d"}

	err := run.RunAndChangeDir(stdOutput, errOutput, "", composeExe, append(p, params...)...)
	if err != nil {
		return fmt.Errorf("failed RunAndChangeDir docker-compose stop, err: %v", err)
	}
	return nil
}
