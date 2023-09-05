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

import (
	"agent/config"
	"strings"

	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/file/fileutil"
)

func GetSystemdDns() (string, []string, error) {

	f := config.Config.Box.DnsConfigFile
	content, err := fileutil.ReadFromFile(f)
	if err != nil {
		return "", nil, err
	}
	s := string(content)

	DNS := ""
	k1 := "DNS="
	i := strings.Index(s, k1)
	if i > 0 {
		m := strings.Index(s[i+len(k1):], "\n")
		DNS = s[i+len(k1):][:m]
	}

	FallbackDNS := ""
	k2 := "FallbackDNS="
	j := strings.Index(s, k2)
	if j > 0 {
		m := strings.Index(s[j+len(k2):], "\n")
		FallbackDNS = s[j+len(k2):][:m]
	}

	DNS2 := []string{}
	arr := strings.Split(FallbackDNS, " ")
	for _, ele := range arr {
		if len(ele) > 6 {
			DNS2 = append(DNS2, ele)
		}
	}
	return DNS, DNS2, nil
}

func SetSystemdDnsManual(dns ...string) error {
	f := config.Config.Box.DnsConfigFile
	fBackup := config.Config.Box.DnsConfigFileBackup
	content := make([]byte, 0)
	var err error

	if fileutil.IsFileExist(fBackup) { // 恢复
		content, err = fileutil.ReadFromFile(fBackup)
		if err != nil {
			return err
		}
		err = fileutil.WriteToFile(f, content, true)
		if err != nil {
			return err
		}
	} else { //  备份
		content, err = fileutil.ReadFromFile(f)
		if err != nil {
			return err
		}
		err = fileutil.WriteToFile(fBackup, content, true)
		if err != nil {
			return err
		}
	}

	// 修改 dns
	if len(dns) > 0 {
		s := string(content)

		i := strings.Index(s, "DNS=")
		if i > 0 {
			m := strings.Index(s[i:], "\n")
			oldDNSLine := s[i:][:m]
			s = strings.ReplaceAll(s, oldDNSLine, "DNS="+dns[0])
		}

		j := strings.Index(s, "FallbackDNS=")
		if j > 0 {
			m := strings.Index(s[j:], "\n")
			oldFallbackDNSLine := s[j:][:m]
			s = strings.ReplaceAll(s, oldFallbackDNSLine, oldFallbackDNSLine+" "+strings.Join(dns[1:], " "))
		}

		if err = fileutil.WriteToFile(f, []byte(s), true); err != nil {
			return err
		}

		return restartSystemdResolved()
	}

	return nil
}

func SetSystemdDnsDefault() error {
	f := config.Config.Box.DnsConfigFile
	fBackup := config.Config.Box.DnsConfigFileBackup

	if fileutil.IsFileExist(fBackup) {
		content, err := fileutil.ReadFromFile(fBackup)
		if err != nil {
			return err
		}
		err = fileutil.WriteToFile(f, content, true)
		if err != nil {
			return err
		}
		return restartSystemdResolved()
	}
	return nil
}

func restartSystemdResolved() error {
	// systemctl restart systemd-resolved
	params := []string{"restart", "systemd-resolved"}
	if err := runCmd2("systemctl", params); err != nil {
		logger.AppLogger().Warnf("SetWirelessIpAuto, clear original ipv4.addresses failed. err:%v", err)
		return err
	}
	return nil
}
