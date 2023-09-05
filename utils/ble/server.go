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
 * @Date: 2021-11-10 10:44:41
 * @LastEditors: wenchao
 * @LastEditTime: 2021-11-22 11:22:55
 * @Description:
 */
// +build

package ble

import (
	"fmt"

	"agent/utils/ble/option"
	"agent/utils/ble/service"

	"agent/utils/logger"

	"github.com/paypal/gatt"
)

var serviceName string
var serviceUUID string

var mtu int

func GetMtu() int {
	return mtu
}

func Start(serviceName_, serviceUUID_ string, onRecvCb service.OnRecvCallbackFunc, bleSendSpendMS uint) error {
	// fmt.Printf("before RegisterRecvCallBack \n")
	serviceName = serviceName_
	serviceUUID = serviceUUID_
	mtu = 200
	service.SetSendSpeed(bleSendSpendMS)

	service.RegisterRecvCallBack(onRecvCb)

	// fmt.Printf("before gatt.NewDevice \n")
	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		fmt.Printf("Failed to open device, err: %s\n", err)
		logger.AppLogger().Warnf("Failed to open device, err: %v", err)
		return fmt.Errorf("Failed to open device, err: %s", err)
	}

	// fmt.Printf("before d.Handle \n")
	// Register optional handlers.
	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) {
			fmt.Printf("Connect, ID:%v, MTU:%v, certral:%v\n", c.ID(), c.MTU(), c)
			mtu = c.MTU()
			logger.AppLogger().Infof("#### BLUETOOTH Connect, ID:%v, MTU:%v, certral:%v", c.ID(), c.MTU(), c)
			service.SetConnected(true)
		}),
		gatt.CentralDisconnected(func(c gatt.Central) {
			fmt.Printf("Disconnect: ID:%v, certral:%v\n", c.ID(), c)
			logger.AppLogger().Warnf("#### BLUETOOTH Disconnect, ID:%v, MTU:%v, certral:%v", c.ID(), c.MTU(), c)
			service.SetConnected(false)
		}),
	)

	// fmt.Printf("d.Init \n")
	d.Init(onStateChanged)
	// fmt.Printf("after d.Init \n")
	select {}
	// fmt.Printf("will return \n")
	return nil
}

// A mandatory handler for monitoring device state.
func onStateChanged(d gatt.Device, s gatt.State) {
	// fmt.Printf("State: %s\n", s)
	logger.AppLogger().Infof("BLUETOOTH onStateChanged, State: %s", s)

	switch s {
	case gatt.StatePoweredOn:
		// Setup GAP and GATT services for Linux implementation.
		// OS X doesn't export the access of these services.
		err := d.AddService(service.NewGapService(serviceName)) // no effect on OS X
		if err != nil {
			logger.AppLogger().Warnf("AddService NewGapService return err:%v", err)
		}
		err = d.AddService(service.NewGattService()) // no effect on OS X
		if err != nil {
			logger.AppLogger().Warnf("AddService NewGattService return err:%v", err)
		}

		// A simple count service for demo.
		s1 := service.NewCountService(serviceUUID)
		err = d.AddService(s1)
		if err != nil {
			logger.AppLogger().Warnf("AddService s1 return err:%v", err)
		}

		// A fake battery service for demo.
		// s2 := service.NewBatteryService()
		// d.AddService(s2)

		// Advertise device name and service's UUIDs.
		// d.AdvertiseNameAndServices(serviceName, []gatt.UUID{s1.UUID(), s2.UUID()})
		err = d.AdvertiseNameAndServices(serviceName, []gatt.UUID{s1.UUID()})
		if err != nil {
			logger.AppLogger().Warnf("AdvertiseNameAndServices return err:%v", err)
		}

		// Advertise as an OpenBeacon iBeacon
		err = d.AdvertiseIBeacon(gatt.MustParseUUID("AA6062F098CA42118EC4193EB73CCEB6"), 1, 2, -59)
		if err != nil {
			logger.AppLogger().Warnf("AdvertiseIBeacon  return err:%v", err)
		}

	default:
	}
}
