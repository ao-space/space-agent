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

package routers

import (
	"agent/biz/web/handler/bind/bindinit"
	"agent/biz/web/handler/bind/com/progress"
	"agent/biz/web/handler/bind/com/start"
	internetserviceconfig "agent/biz/web/handler/bind/internet/service/config"
	"agent/biz/web/handler/bind/password"
	"agent/biz/web/handler/bind/revoke"
	"agent/biz/web/handler/bind/space/create"
	"agent/biz/web/handler/certificate"
	"agent/biz/web/handler/device"
	"agent/biz/web/handler/did/document"
	"agent/biz/web/handler/did/document/method"
	did_document_password "agent/biz/web/handler/did/document/password"
	"agent/biz/web/handler/network"
	"agent/biz/web/handler/pair"
	pairadmin "agent/biz/web/handler/pair/admin"
	pairnet "agent/biz/web/handler/pair/net"
	"agent/biz/web/handler/passthrough"
	"agent/biz/web/handler/space"
	"agent/biz/web/handler/status"
	switchplatform "agent/biz/web/handler/switch-platform"
	"agent/biz/web/handler/system"
	"agent/biz/web/handler/upgrade"
	"agent/config"
	"agent/res"
	"agent/utils/logger"
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spkg/zipfs"
)

func ExternalRouter() *gin.Engine {
	router := gin.Default()

	addHtmlZipHandler(router)

	agent := router.Group("/agent")
	{
		v1 := agent.Group("/v1")
		{
			api := v1.Group("/api")
			{
				// pair apis
				api.POST("/initial", pair.Initial)
				api.POST("/pairing", pair.Pairing)
				api.POST("/pubkeyexchange", pair.PubKeyExchange)
				api.POST("/keyexchange", pair.KeyExchange)
				api.POST("/setpassword", pair.SetPassword)
				//if config.Config.DebugMode {
				//	api.POST("/reset", pair.Reset)
				//}

				// device apis
				dev := api.Group("/device")
				{
					dev.GET("/ability", device.Ability)
				}

				admin := api.Group("/admin")
				{
					admin.POST("/revoke", pairadmin.Revoke)
				}

				pairapis := api.Group("/pair")
				{
					pairapis.POST("/tryout/code", pair.TryOutCode)
					pairapis.POST("/init", pair.TryOutCode)
					pairapis.GET("/net/localips", pairnet.LocalIps)
					pairapis.GET("/net/netconfig", pairnet.NetConfig)
					pairapis.GET("/init", pair.Init)
				}

				bind := api.Group("/bind")
				{
					bind.GET("/init", bindinit.Init)
					bind.POST("/com/start", start.Start)
					bind.GET("/com/progress", progress.Progress)
					bind.POST("/space/create", create.Create)
					bind.POST("/internet/service/config", internetserviceconfig.PostConfig)
					bind.GET("/internet/service/config", internetserviceconfig.GetConfig)
					bind.POST("/password/verify", password.Verify)
					bind.POST("/revoke", revoke.Revoke)
				}

				api.GET("/space/ready/check", space.ReadyCheck)

				networkGroup := api.Group("/network")
				{
					networkGroup.POST("/config", network.PostNetworkConfig)
					networkGroup.GET("/config", network.GetNetworkConfig)
					networkGroup.POST("/ignore", network.NetworkIgnore)
				}

				api.POST("/passthrough", passthrough.Passthrough)
				api.POST("/switch", switchplatform.SwitchPlatform)
				api.GET("/switch/status", switchplatform.SwitchStatusQuery)

				certGroup := v1.Group("/cert")
				{
					certGroup.GET("/get", certificate.GetLanCert)
				}

				did := api.Group("/did")
				{
					did.GET("/document", document.GetDIDDocument)
					did.PUT("/document/password", did_document_password.UpdateDocumentPassword)
					did.PUT("/document/method", method.UpdateDocumentMethod)
				}
			}
		}
		agent.GET("/status", status.Status)
		agent.GET("/info", status.Info)

		if config.Config.DebugMode {
			agent.GET("/logs", status.Logs)
			agent.GET("/log", status.Logs)
		}
	}

	return router
}

func InternalRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/agent/v1/api")
	{

		deviceGroup := v1.Group("/device")
		{
			deviceGroup.GET("/info", device.Info)
			deviceGroup.GET("/version", device.Version)
			deviceGroup.GET("/ability", device.Ability)

			deviceGroup.GET("/localips", pairnet.LocalIpsDevice)
			deviceGroup.GET("/netconfig", pairnet.NetConfigDevice)
		}

		upgradeApp := v1.Group("/upgrade")
		{
			upgradeApp.GET("/config", upgrade.GetUpgradeConfig)
			upgradeApp.POST("/config", upgrade.SetUpgradeConfig)
			upgradeApp.POST("/download", upgrade.StartDownload)
			upgradeApp.POST("/install", upgrade.StartUpgrade)
			upgradeApp.GET("/status", upgrade.GetTaskStatus)
		}

		networkGroup := v1.Group("/network")
		{
			networkGroup.POST("/config", network.PostNetworkConfig)
			networkGroup.GET("/config", network.GetNetworkConfig)
			networkGroup.POST("/ignore", network.NetworkIgnore)
		}

		systemGroup := v1.Group("/system")
		{
			systemGroup.POST("/shutdown", system.SystemShutdown)
			systemGroup.POST("/reboot", system.SystemReboot)
		}
		certGroup := v1.Group("/cert")
		{
			certGroup.GET("/get", certificate.GetLanCert)
		}

		bindGroup := v1.Group("/bind")
		{
			bindGroup.POST("/internet/service/config", internetserviceconfig.PostConfig)
			bindGroup.GET("/internet/service/config", internetserviceconfig.GetConfig)
		}

		did := v1.Group("/did")
		{
			did.GET("/document", document.GetDIDDocument)
			did.PUT("/document/password", did_document_password.UpdateDocumentPassword)
			did.PUT("/document/method", method.UpdateDocumentMethod)
		}

	}
	return router
}

func addHtmlZipHandler(router *gin.Engine) error {

	buf := res.GetContentStaticHtmlZip()
	reader := bytes.NewReader(buf)
	fs, err := zipfs.NewFromReaderAt(reader, int64(len(buf)), nil)
	if err != nil {
		err1 := fmt.Errorf("Failed NewFromReaderAt, err: %v", err)
		fmt.Printf("%+v\n", err1)
		logger.AppLogger().Errorf("%+v", err1)
	}

	router.GET("/", func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/boxdocker/index.html")
		}
	}())

	// router.Use(gin.WrapH(zipfs.FileServer(fs)))
	router.GET("/boxdocker/*any", gin.WrapH(zipfs.FileServer(fs)))

	// go http.ListenAndServe(":5678", zipfs.FileServer(fs))

	return nil
}
