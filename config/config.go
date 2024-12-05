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

package config

import (
	hardware_util "agent/utils/hardware"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/jinzhu/configor"
	"gopkg.in/yaml.v2"
)

var Version string
var VersionNumber string
var runInDocker bool

// 以下配置项尽量都包含默认值. 让用户零配置启动.
// 系统OS镜像构建时，没有复制默认的 yml 配置文件. 原因是这里已经设置了默认值，不想在项目组重复设置默认值.

var Config = struct {
	DebugMode                                 bool `default:"false"` // 调试模式。会控制是否打开 swagger 等。
	OverwriteDockerCompose                    bool `default:"true"`  // 启动时是否覆盖 docker-compose.yml。"true" 表示覆盖。
	EnableSecurityChip                        bool `default:"true"`  // 是否启用加密芯片。
	EncryptLanSessionData                     bool `default:"true"`  // 加密局域网通信数据
	EnableBackupRestoreSupportWhenRunAsDocker bool `default:"false"` // 容器化部署时, 是否启动备份恢复功能.

	// 仅在 DebugMode=true 才生效, "0" 为不退出程序。测试中断初始化流程，system-agent退出程序的位置。
	// 可以是(磁盘初始化场景下1_XXX表示磁盘初始化, 2_XXX 表示磁盘扩容):
	// 1_FormatAndParted 或 2_FormatAndParted
	// 1_DiskRaid 或 2_DiskRaid
	// 1_DiskEncrypt 或 2_DiskEncrypt
	// 1_DiskMount 或 2_DiskMount
	// 1_MigrateDockers_0
	// 1_MigrateDockers_1
	// ...
	// 1_MigrateDockers_8,
	// 1_MmigrateDockers

	DebugExitProgramPlace string `default:"0"`
	Dsn                   string `default:"https://f0cdec64ac144e95b85c354daf8ef90d@sentry.eulix.xyz/4"` // sentry dsn

	APPName string `default:"system-agent"`

	Version     bool
	EnablePprof bool `default:"false"`

	EnvDefaultVal struct {
		SYSTEM_AGENT_URL_DEVICE_INFO string `default:"http://172.17.0.1:5680/agent/v1/api/device/info"`
		SYSTEM_AGENT_URL_BASE        string `default:"http://172.17.0.1:5680/agent/v1/api"`
	}

	AliveChecker struct {
		GateWay struct {
			Enable     bool   `default:"false"`
			UrlGateway string `default:"http://localhost:8080/space/status"`
		}
		// DockerAliveCheckIntervalSec    uint32 `default:"30"` // 容器保活检测的间隔(秒)
		// LogVersionInfoIntervalSec      uint32 `default:"60"` // 输出版本等信息的间隔(秒)
		// TestPlatformNetworkIntervalSec uint32 `default:"30"` // 测试与平台的网络连通性的间隔(秒)

		DockerAliveCheckIntervalSec    uint32 `default:"6"`    // 容器保活检测的间隔(秒)
		LogVersionInfoIntervalSec      uint32 `default:"3600"` // 输出版本等信息的间隔(秒)
		TestPlatformNetworkIntervalSec uint32 `default:"15"`   // 测试与平台的网络连通性的间隔(秒)
	}

	Log struct {
		// Path         string `default:"/var/log/"`
		AoLogDirBase         string `default:"/opt/logs/"` // 所有程序日志存储目录
		AoLogDirBaseSizLimit int64  `default:"3221225472"` // 所有程序日志存储目录总大小限制(单位是字节). 超过了会定时清理.
		AoLogMaxDayLimit     int    `default:"7"`          // 所有程序日志最多留存天数.
		AoLogDirCheckCronExp string `default:"0 0 4 * *"`  // 所有程序日志目录定时检测的 cron 表达式. 比如 0 0 4 * * 表示每天 4时 0分 0秒执行一次逻辑。
		Path                 string `default:"/opt/logs/system-agent/"`
		Filename             string `default:"system-agent"`
		MaxAge               uint32 `default:"7"`  // 文件最多保留多少天后被覆盖.
		RotationSize         int64  `default:"10"` // 多少MBytes生成新文件.
		RotationCount        uint32 `default:"20"` // 多少个文件之后覆盖最早的文件.
		Level                uint32 `default:"6"`  // InfoLevel=4, DebugLevel=5, logrus.TraceLevel=6 . 暂时没用到.
	}

	EnableKeyboard bool `default:"false"`

	BlueTooth struct {
		Enable             bool   `default:"true"`
		Service            string `default:"eulixspace-"`
		Adapter            string `default:"hci0"`
		BleSendSpendMS     uint   `default:"100"`
		BleSendSpendFastMS uint   `default:"5"`
	}

	Box struct {
		CpuIdStoreFile            string `default:"/etc/ao-space/hardware/cpuid.data"`
		SnNumberStoreFile         string `default:"/etc/ao-space/hardware/snnumber.data"`
		HostIpFile                string `default:"/etc/ao-space/hardware/host_ip.data"`
		ApplyEmailStoreFile       string `default:"/etc/ao-space/box/apply_email.data"`
		BoxInfoFile               string `default:"/etc/ao-space/box_info.json"`
		InternetServiceConfigFile string `default:"/etc/ao-space/internet_service_config.json"`
		SwithStatusFile           string `default:"/etc/ao-space/switch_status.json"`
		WifiNamePasswdFile        string `default:"/etc/ao-space/wifi.pwd"`         // wifi名称和密码,以便下次连接
		PingHost                  string `default:"www.baidu.com"`                  // 测试盒子外网连通性
		BoxMetaAdminPair          string `default:"/etc/ao-space/meta/admin/admin"` // 管理员是否配对.

		RandDockercomposePassword  string `default:"/etc/ao-space/box/rand_docker_compose_password.data"`
		RandDockercomposeRedisPort string `default:"/etc/ao-space/box/rand_redis_port.data"`

		PublicSharedInfoFile string `default:"/etc/ao-space/meta/shared/shared_info.json"`

		UpgradeCheckIntervalMs            uint `default:"1000"` // 检查 upgrade 升级状态的时间间隔.
		InitialEstimateTimeSecRunInDocker uint `default:"600"`  // 绑定预计时间(容器中运行时)

		RegisterBoxRetryIntervalSec uint32 `default:"20"` // 注册设备失败重试间隔
		RegisterBoxRetryTimes       uint32 `default:"12"` // 注册设备失败重试次数

		SysconfigNetworkIpAddressFile    string `default:"/etc/sysconfig/network-scripts"`
		SysconfigNetworkWIFIPasswordFile string `default:"/etc/sysconfig/network-scripts"`
		DnsConfigFile                    string `default:"/etc/systemd/resolved.conf"`
		DnsConfigFileBackup              string `default:"/etc/systemd/resolved.conf.backup"`

		SecurityChipAgentSockAddr string `default:"/opt/tmp/eulixspace-security-agent.sock"`

		SecurityChipAgentHttpAddr      string `default:"http://172.17.0.1:9200/security/v1/api"`
		SecurityChipAgentHttpLocalAddr string `default:"http://172.17.0.1:9200/security/v1/api"`

		BoxKey struct {
			RsaKeyFile    string `default:"/etc/ao-space/box_key.pem"`
			RsaPubKeyFile string `default:"/etc/ao-space/box_key_pub.pem"`
		}

		Disk struct {
			DiskInitialInfoFile  string `default:"/etc/ao-space/disk/disk_initial_info.json"`
			DeviceUuidRecordFile string `default:"/etc/ao-space/disk/disk_uuid_record.json"`

			MountPathRoot                      string `default:"/mnt/ao-space/data"`                       // 磁盘挂载的根路径。
			MountPathHddNamePrefix             string `default:"hdd_"`                                     // 挂载机械盘目录前缀
			MountPathM2NamePrefix              string `default:"nvme_"`                                    // 挂载m2盘目录前缀
			MountPathRaid1NamePrefix           string `default:"raid1_"`                                   // 挂载raid1盘目录前缀
			MountPathHddNameSuffix             string `default:"/mountpoint"`                              // 挂载机械盘目录后缀
			MountPathM2NameSuffix              string `default:"/mountpoint"`                              // 挂载m2盘目录后缀
			MountPathRaid1NameSuffix           string `default:"/mountpoint"`                              // 挂载raid1盘目录后缀
			StorageVolumePath                  string `default:"/home/eulixspace_file_storage/parts"`      // 挂载给 fileapi 存储容器的目录
			NoDisksFileStoragePath             string `default:"/home/eulixspace/data"`                    // 无独立磁盘时的 fileapi 数据存储目录
			NoDisksFileStoragePathDockerDeploy string `default:"/home/eulixspace_link/data"`               // 无独立磁盘时的 fileapi 数据存储目录(Docker 部署时)
			StorageDummy                       string `default:"dummy"`                                    // 无独立磁盘时的虚拟磁盘目录
			DiskEncryptKeyPath                 string `default:"/var/cache"`                               // 磁盘加密密钥临时存储位置
			DiskEncryptDiskMapperPrefix        string `default:"luks_"`                                    // 磁盘加密 mapper name 前缀
			FileStorageVolumePathPrefix        string `default:"ao_part_"`                                 // 挂载给 FileApi 容器的目录前缀 /home/eulixspace_file_storage/parts/ao_part_nvme_3841a657
			FileStorageInnerDataPath           string `default:"/eulixspace/data"`                         // fileapi 内部路径
			DiskSharedInfoFile                 string `default:"/etc/ao-space/meta/shared/disk_info.json"` // fileapi 共享磁盘信息
			DockerVolumePlaceholderInHost      string `default:"placeholder_file_storage_hosts"`           // docker-compose_gen2.yml 中配置的占位符, 代表 hdd/ssd 磁盘的对象文件存储宿主机的目录
			DockerVolumePlaceholderInContainer string `default:"placeholder_file_storage_container"`       // docker-compose_gen2.yml 中配置的占位符, 代表对象文件存储目录挂载到容器内部的目录
		}

		ClientKey struct {
			RsaPubKeyFile string `default:"/etc/ao-space/client_key_pub.pem"`
			RsaPriKeyFile string `default:"/etc/ao-space/client_key_pri.pem"`
			SharedSecret  string `default:"/etc/ao-space/shared_secret.key"`
		}

		UpgradeConfig struct {
			SettingsFile string `default:"/etc/ao-space/upgrade/settings.json"`
			// AutoDownload bool   `default:"true" json:"autoDownload"`
			// AutoInstall  bool   `default:"true" json:"autoInstall"`
		}

		Avahi struct {
			BtIdHashPlaceHolder     string `default:"placeholdercccccc"`
			BtIdHashPrefix          string `default:"eulixspace-"`
			BtIdHashXmlConfigPrefix string `default:"btidhash"`
			DeviceModelPlaceHolder  string `default:"placeholdermodel"`
			SSLPort                 string `default:"443"`
			WebPort                 string `default:"80"`
			AvahiServiceName        string `default:"avahi-daemon"`
			BtIdHashLen             uint16 `default:"6"`
		}

		Loki struct {
			Placeholder string `default:"placeholderlokihosts"`
		}

		RunInDocker struct {
			AoSpaceDataDirEnv          string `default:"AOSPACE_DATADIR"`  // 数据存储目录
			RunNetworkModeEnv          string `default:"RUN_NETWORK_MODE"` // docker-compose 中的 nginx 网络配置键名称
			AoSpaceSingleDockerModeEnv string `default:"AOSPACE_SINGLE_DOCKER_MODE"`
		}

		Cert struct {
			CertDir           string `default:"/etc/ao-space/certs/"`
			ACMERegisterEmail string `default:"service@ao.space"`
		}

		DID struct {
			RootPath string `default:"/etc/ao-space/did"`

			DBFileName string `default:"did.leveldb"`
		}

		SnNumberModelLength  int    `default:"3"`   // sn 序号最前面表示型号的字符长度
		SnNumberModelContent string `default:"002"` // sn 序号最前面"在线试用"型号的字符内容
	}

	Web struct {
		DefaultListenAddr                string `default:":5678"`
		DockerLocalListenAddr            string `default:"172.17.0.1:5680"` // TODO: 这里应该动态获取. 为了快速开发需要，暂时先固定这样. 多块 docker 网卡时可能会有问题.
		DockerLocalListenAddrRunInDocker string `default:":5680"`           // 在容器中运行时，暴露给网关等容器调用。
	}

	NetworkCheck struct {
		ThirdPartyHost struct {
			Url string `default:"www.baidu.com"`
		}
		CloudHost struct {
			Url string `default:"ao.space"`
		}
		CloudIpv4 struct {
			Url string `default:"121.41.7.103"`
		}
		CloudStatusHost struct {
			Url string `default:"https://services.ao.space/platform/status"`
		}
		CloudStatusIpv4 struct {
			Url string `default:"http://121.41.7.103"`
		}
		BoxStatusPath struct {
			Url string `default:"space/status"`
		}
	}

	PSPlatform struct {
		APIBase struct {
			Url string `default:"https://ao.space"` // 生产环境 https://services.ao.space
		}
	}

	AppStore struct {
		Sign struct {
			Url string `default:"https://auth.apps.ao.space"` // 生产环境 https://auth.apps.ao.space
		}
		Api struct {
			Url string `default:"https://api.apps.ao.space"` // 生产环境 https://api.apps.ao.space
		}
	}

	Platform struct {
		WebBase struct {
			Url string `default:"https://ao.space/"`
		}
		APIBase struct {
			Url string `default:"https://ao.space"`
		}

		RegistryBox struct {
			Path string `default:"/v2/platform/boxes"`
		}

		PresetBoxInfo struct {
			Path string `default:"/v2/service/trail/boxinfos"`
		}

		AuthBox struct {
			Path string `default:"/v2/platform/auth/box_reg_keys"`
		}

		Reset struct {
			Path string `default:"/v2/platform/boxes/"` //开发使用，前端暂无使用，暂不维护
		}

		NetworkRemoteApi struct {
			Path string `default:"/v2/platform/servers/network/detail"`
		}

		Status struct {
			Path string `default:"/v2/platform/status"`
		}

		DNSResolution struct {
			Path string `default:"/v2/platform/boxes/{box_uuid}/users/{user_id}/lan-subdomain"`
		}

		Migration struct {
			Path string `default:"/v2/platform/boxes/{box_uuid}/migration"`
		}

		Route struct {
			Path string `default:"/v2/platform/boxes/{box_uuid}/route"`
		}
		Ability struct {
			Path string `default:"/v2/platform/ability"`
		}
		LatestVersion struct {
			Path string `default:"/v1/api/package/box"`
		}
		LatestVersionV2 struct {
			Path string `default:"/v2/service/packages/box/latest"`
		}
	}

	GateWay struct {
		Revoke struct {
			Url string `default:"http://localhost:8080/space/v1/api/gateway/auth/revoke"`
		}
		SwitchPlatform struct {
			Url string `default:"http://localhost:8080/space/v1/api/space/platform"`
		}
		LanPort    uint16 `default:"80"`
		TlsLanPort uint16 `default:"443"`

		APIRoot struct {
			Url string `default:"http://localhost:8080/space"`
		}
	}

	Account struct {
		User struct {
			Url string `default:"http://localhost:8080/space/v1/api/user"`
		}
		Member struct {
			Url string `default:"http://localhost:8080/space/v1/api/member/list"`
		}
		AdminCreate struct {
			Url string `default:"http://localhost:8080/space/v1/api/admin/create"`
		}
		SpaceAdmin struct {
			Url string `default:"http://localhost:8080/space/v1/api/space/admin"`
		}
		NetworkChannelInfo struct {
			Url string `default:"http://localhost:8080/space/v1/api/device/network/channel/info"`
		}
		NetworkChannelWan struct {
			Url string `default:"http://localhost:8080/space/v1/api/device/network/channel/wan"`
		}

		AdminSetPassword struct {
			Url string `default:"http://localhost:8080/space/v1/api/admin/passwd/set"`
		}

		AdminPasswordCheck struct {
			Url string `default:"http://localhost:8080/space/v1/api/admin/passwd/check"`
		}

		AdminRevoke struct {
			Url string `default:"http://localhost:8080/space/v1/api/user/client/revoke"`
		}

		AdminInitial struct {
			Url string `default:"http://localhost:8080/space/v1/api/admin/inital/status"`
		}

		Migrate struct {
			Url string `default:"http://localhost:8080/space/v1/api/user/migration"`
		}
	}

	Upgrade struct {
		Url string `default:"http://aospace-upgrade:5681/upgrade/v1/api/start"`
	}

	GTClient struct {
		ConfigPath string `default:"/etc/ao-space/gt/aonetwork-client.yml"`
	}

	Docker struct {
		APIVersion                            string `default:"1.39"`
		ComposeFile                           string `default:"/opt/tmp/docker-compose.yml"`
		CustomComposeFile                     string `default:"/etc/ao-space/docker-compose.yml"`
		UpgradeComposeFile                    string `default:"/etc/ao-space/aospace-upgrade.yml"`
		NetworkName                           string `default:"ao-space"`
		NginxContainerName                    string `default:"aospace-nginx"`
		NetworkClientContainerName            string `default:"aonetwork-client"`
		DockerEngineReadyWaitingCheckInterval uint32 `default:"3"` // 容器引擎启动等待检测间隔(秒)
		DockerStorageFile                     string `default:"/etc/sysconfig/docker-storage"`
		DockerUpRetryIntervalSec              uint32 `default:"3"` // 容器启动失败重试间隔
		DockerUpRetryTimes                    uint32 `default:"3"` // 容器启动失败重试次数

		VolumeDirLink string `default:"/home/eulixspace_link"` // 容器挂载目录软链接
		VolumeDirReal string `default:"/home/eulixspace"`      // 容器实际挂载目录

		RegistryUrlNoAuth string `default:"hub.eulix.xyz/cicada-private/"`
	}
	RunTime struct {
		BasePath          string `default:"/var/system-agent/"`
		DBDir             string `default:".db"`
		PkgDir            string `default:"pkg"`
		UpgradeCollection string `default:"upgrade"`
		TaskResource      string `default:"task"`
		SocketFile        string `default:"upgrade.sock"`
	}

	Notification struct {
		UpgradeRecordFile string `default:"/etc/ao-space/upgrade_notification_push.json"`
	}

	Redis struct {
		Addr     string `default:"127.0.0.1:6379"`
		Password string `default:"mysecretpassword"`
		DBIndex  int    `default:"0"`
	}
}{}

var ConfFile = "/etc/ao-space/system-agent.yml"   // 默认配置文件
var ConfFileCustom1 = "/opt/tmp/system-agent.yml" // 用户自定义配置文件

var SpaceMountPath = `/aospace`

var flagConfFile *string

func init() {
	testing.Init()
	runInDocker = hardware_util.RunningInDocker()
	if runInDocker {
		ConfFile = SpaceMountPath + ConfFile
		ConfFileCustom1 = SpaceMountPath + ConfFileCustom1
	}

	v := false
	flagConfFile = flag.String("config", ConfFile, "configuration file")
	if !v {
		removeExistingConfigFile(*flagConfFile)
	}
	configor.New(&configor.Config{AutoReload: true,
		AutoReloadInterval: time.Second * 15,
		AutoReloadCallback: func(config interface{}) {
			// fmt.Printf("config file changed:\n%+v\n", config)
		}}).Load(&Config, ConfFileCustom1, *flagConfFile)

	Config.Version = v
	if !v {
		modifyConfigWhenRunInDocker()
		createLogFileDir()
		writeDefaultConfigFile(*flagConfFile)
	}
}

func modifyConfigWhenRunInDocker() {
	// fmt.Printf("+++++++++++++++++++++++ modifyConfigWhenRunInDocker \n")
	////////////////////////////////////////////////////////////
	// 注意!!!
	// docker 运行时, 需要把容器内运行的数据保存在宿主机上，以方便其他容器访问，同时也可以让用户
	// 修改保存路径。
	// 新增路径的配置项时，如果docker 中运行时, 需要在这里修改默认路径 !!!
	if runInDocker {
		p := []*string{
			&Config.Log.AoLogDirBase,
			&Config.Log.Path,
			&Config.Box.CpuIdStoreFile,
			&Config.Box.SnNumberStoreFile,
			&Config.Box.HostIpFile,
			&Config.Box.ApplyEmailStoreFile,
			&Config.Box.BoxInfoFile,
			&Config.Box.InternetServiceConfigFile,

			&Config.Box.SwithStatusFile,
			&Config.Box.WifiNamePasswdFile,
			&Config.Box.BoxMetaAdminPair,
			&Config.Box.RandDockercomposePassword,
			&Config.Box.RandDockercomposeRedisPort,

			&Config.Box.PublicSharedInfoFile,
			&Config.Box.BoxKey.RsaKeyFile,
			&Config.Box.BoxKey.RsaPubKeyFile,
			&Config.Box.Disk.DiskInitialInfoFile,
			&Config.Box.Disk.DeviceUuidRecordFile,
			&Config.Box.Disk.DiskSharedInfoFile,
			&Config.Box.ClientKey.RsaPubKeyFile,
			&Config.Box.ClientKey.RsaPriKeyFile,
			&Config.Box.ClientKey.SharedSecret,
			&Config.Box.UpgradeConfig.SettingsFile,
			&Config.Box.Cert.CertDir,
			&Config.Docker.ComposeFile,
			&Config.Docker.CustomComposeFile,
			&Config.RunTime.BasePath,
			&Config.Notification.UpgradeRecordFile,
			&Config.Box.Disk.StorageVolumePath,
			&Config.Box.Disk.NoDisksFileStoragePath,
			&Config.Box.Disk.NoDisksFileStoragePathDockerDeploy,
			&Config.GTClient.ConfigPath,
			&Config.Box.DID.RootPath}

		for _, v := range p {
			*v = SpaceMountPath + *v
		}
		// fmt.Printf("Config.Box.SnNumberStoreFile: %v \n", Config.Box.SnNumberStoreFile)

		// 调用地址修改
		if !strings.EqualFold(os.Getenv(Config.Box.RunInDocker.AoSpaceSingleDockerModeEnv), "true") {
			// aospace-all-in-one
			Config.EnvDefaultVal.SYSTEM_AGENT_URL_DEVICE_INFO = "http://aospace-all-in-one:5680/agent/v1/api/device/info"
			Config.EnvDefaultVal.SYSTEM_AGENT_URL_BASE = "http://aospace-all-in-one:5680/agent/v1/api"
			// DockerLocalListenAddr
			Config.Web.DockerLocalListenAddr = "aospace-all-in-one:5680"
			addr := []*string{&Config.AliveChecker.GateWay.UrlGateway,
				&Config.GateWay.Revoke.Url,
				&Config.GateWay.APIRoot.Url,
				&Config.Account.User.Url,
				&Config.Account.Member.Url,
				&Config.Account.AdminCreate.Url,
				&Config.Account.SpaceAdmin.Url,
				&Config.Account.NetworkChannelInfo.Url,
				&Config.Account.NetworkChannelWan.Url,
				&Config.Account.AdminSetPassword.Url,
				&Config.Account.AdminPasswordCheck.Url,
				&Config.Account.AdminRevoke.Url,
				&Config.Account.AdminInitial.Url,
				&Config.Account.Migrate.Url,
				&Config.Redis.Addr,
				&Config.GateWay.SwitchPlatform.Url}
			for _, v := range addr {
				*v = strings.ReplaceAll(*v, "localhost:8080", "aospace-gateway:8080")
				*v = strings.ReplaceAll(*v, "127.0.0.1:6379", "aospace-redis:6379")
			}
			Config.GateWay.LanPort = 12841
			Config.GateWay.TlsLanPort = 18569
		} else {
			Config.EnvDefaultVal.SYSTEM_AGENT_URL_DEVICE_INFO = "http://localhost:5680/agent/v1/api/device/info"
			Config.EnvDefaultVal.SYSTEM_AGENT_URL_BASE = "http://localhost:5680/agent/v1/api"
			Config.Web.DockerLocalListenAddr = "localhost:5680"
			Config.GateWay.LanPort = 80
			Config.GateWay.TlsLanPort = 443
		}

	}
	// fmt.Printf("######################## Config.Log.Path:%v\n",
	// 	Config.Log.Path)

	////////////////////////////////////////////////////////////
}

func createLogFileDir() {
	if !fileutil.IsFileExist(Config.Log.Path) {
		fileutil.CreateDirRecursive(Config.Log.Path)
	}
}

func readExistingConfigFile(f string) {
	b, err := fileutil.ReadFromFile(f)
	if err != nil {
		fmt.Printf("ReadFromFile config %v fail, err: %+v\n\n", f, err)
	} else {
		err1 := yaml.Unmarshal(b, &Config)
		if err1 != nil {
			fmt.Printf("Unmarshal config file %v fail, err: %+v\n\n", f, err1)
		}
	}
}

func removeExistingConfigFile(f string) {
	err := fileutil.WriteToFile(f, []byte{}, true)
	if err != nil {
		fmt.Printf("failed removeExistingConfigFile, WriteToFile file:%v, err: %+v\n",
			f, err)
		return
	}
}

func writeDefaultConfigFile(f string) {
	out, err := yaml.Marshal(Config)
	if err != nil {
		fmt.Printf("failed  yaml.Marshal: %+v\n", err)
		return
	}
	// fmt.Printf("@@@@ DefaultConfigFile:\n%+v\n\n", string(out))
	fileutil.WriteToFile(f, out, true)
}

func UpdateRedisConfig(addr, password string) {
	// fmt.Printf("#### UpdateRedisConfig, addr:%+v, password:%+v\n", addr, password)
	Config.Redis.Addr = addr
	Config.Redis.Password = password
	writeDefaultConfigFile(*flagConfFile)
}
