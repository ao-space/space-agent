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
 * @LastEditTime: 2022-06-06 11:28:18
 * @Description: 增加 VersionNumber
 */
package main

import (
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/config"
	"agent/utils/deviceid"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	AgentCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "Print the version number of system-agent")
	AgentCmd.ParseFlags(os.Args)

}

var Version = "dev"
var VersionNumber = ""
var versionFlag bool

var AgentCmd = &cobra.Command{
	Use:   "system-agent",
	Short: "aospace system agent version",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			btid, err := deviceid.GetBtId(config.Config.Box.CpuIdStoreFile, config.Config.Box.SnNumberStoreFile)
			boxUuid := ""
			if err != nil {
				fmt.Printf("failed GetBtId, err=%v\n\n", err)
			}
			boxUuid, _ = deviceid.GetProductId(config.Config.Box.CpuIdStoreFile)
			abilityModel := device_ability.GetAbilityModel()

			fmt.Printf("\nVersion:\n%v\n", Version)

			// fmt.Printf("\nLog.Path:\n%v\n", config.Config.Log.Path)
			fmt.Printf("\nWebBase:\n%v\n", config.Config.Platform.WebBase.Url)
			fmt.Printf("\nAPIBase:\n%v\n", device.GetApiBaseUrl())
			fmt.Printf("\nPairedStatus:\n%v (0: 已经配对, 1: 新盒子, 2: 已解绑)\n", clientinfo.GetAdminPairedStatus())
			fmt.Printf("\nBoxUuid:\n%v\n", boxUuid)
			fmt.Printf("\nBtid:\n%v\n", btid)
			fmt.Printf("\nDeviceModelNumber:\n%v\n", abilityModel.DeviceModelNumber)
			if len(abilityModel.SnNumber) > 0 {
				fmt.Printf("\nSnNumber:\n%v\n", abilityModel.SnNumber)
			} else {
				fmt.Printf("\nSnNumber:\n%v\n", "unavailable")
			}
			fmt.Printf("\nBox QR code:\n%v\n\n", device.GetQrCode())

			GracefullExit()
		}

	},
}
