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
 * @Author: wenchao
 * @Date: 2021-12-28 10:42:18
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-24 11:22:52
 * @Description:
 */

package device

//func UpdateMdns() error {
//	logger.AppLogger().Debugf("UpdateMdns, GetDeviceInfo().BtidHash=%v", GetDeviceInfo().BtidHash)
//	// https://code.eulix.xyz/bp/cicada/-/issues/473
//	f := config.Config.Box.Avahi.AvahiConfigFile
//
//	s := res.GetContentAvahiCicadaService()
//	logger.AppLogger().Debugf("UpdateMdns, %v content: %v", f, string(s))
//	if len(s) < 6 {
//		return fmt.Errorf("UpdateMdns, GetContentAvahiCicadaService len error")
//	}
//
//	newstr := strings.ReplaceAll(string(s), config.Config.Box.Avahi.BtIdHashPlaceHolder,
//		GetDeviceInfo().BtidHash[:config.Config.Box.Avahi.BtIdHashLen])
//	// logger.AppLogger().Debugf("UpdateMdns, after Replace placeholder %v new content: %v", f, s)
//	if len(newstr) < 6 {
//		return fmt.Errorf("UpdateMdns, new content len error")
//	}
//
//	newModel := strings.ReplaceAll(newstr, config.Config.Box.Avahi.DeviceModelPlaceHolder,
//		strconv.Itoa(device_ability.GetAbilityModel().DeviceModelNumber))
//	logger.AppLogger().Debugf("UpdateMdns, after Replace placeholder %v new content: %v", f, string(newModel))
//
//	if len(newModel) < 6 {
//		return fmt.Errorf("UpdateMdns, new content len error")
//	}
//	err := fileutil.WriteToFile(f, []byte(newModel), true)
//
//	if err != nil {
//		logger.AppLogger().Warnf("UpdateMdns, failed WriteToFile %v", f)
//		return err
//	}
//	logger.AppLogger().Debugf("UpdateMdns, succ WriteToFile %v", f)
//	return nil
//}
