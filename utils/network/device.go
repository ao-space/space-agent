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

package network

// nmcli con up static2
// or
// nmcli con up eth1
// or
// nmcli con down b0764f9d-0dbc-4f66-a2af-e3ff2465027c
func SetNetworkDeviceUp(device string) error {
	params := []string{"con", "up", device}
	return runCmd(params)
}

// nmcli con down static2
// or
// nmcli con down eth1
// or
// nmcli con up b0764f9d-0dbc-4f66-a2af-e3ff2465027c
func SetNetworkDeviceDown(device string) error {
	params := []string{"con", "down", device}
	return runCmd(params)
}
