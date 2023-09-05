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
 * @Date: 2021-11-10 10:46:41
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-25 13:59:07
 * @Description:
 */
package service

import (
	"time"

	"agent/utils/logger"

	"github.com/paypal/gatt"
)

var chSendQueue chan []byte
var chRecvQueue chan []byte
var bleConnected bool
var bleSendSpendMS uint

type OnRecvCallbackFunc func([]byte)

func init() {
	chSendQueue = make(chan []byte, 128)
	chRecvQueue = make(chan []byte, 128)
	bleConnected = false
	bleSendSpendMS = 150
}

func SetSendSpeed(_bleSendSpendMS uint) {
	bleSendSpendMS = _bleSendSpendMS
}

func SetConnected(conn bool) {
	logger.AppLogger().Debugf("before SetConnected, bleConnected:%v", bleConnected)
	if conn {
		logger.AppLogger().Debugf("conn:%v", conn)
		bleConnected = conn
	} else {
		logger.AppLogger().Debugf("before SendData 0xff, bleConnected:%v", bleConnected)
		SendData([]byte{0xff})
		logger.AppLogger().Debugf("after SendData 0xff, bleConnected:%v", bleConnected)
	}
	logger.AppLogger().Debugf("after SetConnected, bleConnected:%v", bleConnected)
}

func SendData(data []byte) {
	logger.AppLogger().Debugf("BLUETOOTH SendData, bleConnected:%v, len(data):%+v", bleConnected, len(data))
	if bleConnected {
		chSendQueue <- data
	}
}

func emptyChan(c chan []byte) {
	// for {
	// 	logger.AppLogger().Debugf("emptyChan, before get, len chRecvQueue:%v", len(chRecvQueue))
	// 	if len(chRecvQueue) < 1 {
	// 		break
	// 	}
	// 	_, ok := <-chRecvQueue
	// 	logger.AppLogger().Debugf("emptyChan, after get")
	// 	if !ok {
	// 		break
	// 	}
	// 	logger.AppLogger().Debugf("emptyChan loop")
	// }
	// logger.AppLogger().Debugf("emptyChan return")

	for {
		select {
		case x, ok := <-chRecvQueue:
			if ok {
				logger.AppLogger().Debugf("emptyChan, Value was read, len = %v", len(x))
			} else {
				logger.AppLogger().Debugf("emptyChan, Channel closed")
				return
			}
		default:
			logger.AppLogger().Debugf("emptyChan, No value ready, moving on.")
			return
		}
	}
}

func RegisterRecvCallBack(cb OnRecvCallbackFunc) {
	go func() {
		for {
			data := <-chRecvQueue
			logger.AppLogger().Debugf("callback, bleConnected:%v, len(data):%x, cb:%+v",
				bleConnected, len(data), cb)

			if bleConnected == false {
				logger.AppLogger().Debugf("len(chRecvQueue):%x", len(chRecvQueue))
				emptyChan(chRecvQueue)
				logger.AppLogger().Debugf("bleConnected == false, continue, len(chRecvQueue):%x", len(chRecvQueue))
				continue
			}

			if cb != nil {
				cb(data)
			}
		}
	}()
}

func NewCountService(serviceUUID string) *gatt.Service {
	// logger.AppLogger().Debugf("#### NewCountService \n")
	logger.AppLogger().Debugf("BLUETOOTH NewCountService")

	n := 0
	// s := gatt.NewService(gatt.MustParseUUID("09fc95c0-c111-11e3-9904-0002a5d5c51b"))
	s := gatt.NewService(gatt.MustParseUUID(serviceUUID))
	s.AddCharacteristic(gatt.MustParseUUID("11fac9e0-c111-11e3-9246-0002a5d5c51b")).HandleReadFunc(
		func(rsp gatt.ResponseWriter, req *gatt.ReadRequest) {
			logger.AppLogger().Infof("#### HandleReadFunc, can write here")
			// fmt.Fprintf(rsp, "012345678901234567890123456789012345678901234567890000000000")
			n++
		})

	s.AddCharacteristic(gatt.MustParseUUID("16fe0d80-c111-11e3-b8c8-0002a5d5c51b")).HandleWriteFunc(
		func(r gatt.Request, data []byte) (status byte) {
			logger.AppLogger().Debugf("-------- BLUETOOTH recv data len %v", len(data))
			// logger.AppLogger().Debugf("-------- BLUETOOTH recv data:%x, string(data):%v", data, string(data))
			// log.Println("Wrote:", string(data))
			chRecvQueue <- data
			return gatt.StatusSuccess
		})

	s.AddCharacteristic(gatt.MustParseUUID("1c927b50-c116-11e3-8a33-0800200c9a66")).HandleNotifyFunc(
		func(r gatt.Request, n gatt.Notifier) {
			logger.AppLogger().Debugf("HandleNotifyFunc Entry, before for loop")

			// cnt := 0
			for !n.Done() {
				// logger.AppLogger().Debugf("waiting on chSendQueue")

				// fmt.Fprintf(n, "Count: %d", cnt)
				// cnt++

				data := <-chSendQueue
				logger.AppLogger().Debugf("len(chSendQueue):%x, len(data):%v", len(chSendQueue), len(data))
				if n.Done() {
					emptyChan(chSendQueue)
					logger.AppLogger().Debugf("HandleNotifyFunc, n.Done()==true, break")
					break
				}
				if len(data) < 2 {
					emptyChan(chSendQueue)
					logger.AppLogger().Debugf("HandleNotifyFunc, DISCONNECT SIGNAL, break, data:%x", data)
					break
				}
				if bleConnected == false {
					emptyChan(chSendQueue)
					logger.AppLogger().Debugf("HandleNotifyFunc, bleConnected == false, break, len(chSendQueue):%v, data:%x", len(chSendQueue), data)
					break
				}
				// logger.AppLogger().Debugf("n.Write(data) before, data len:%v", len(data))
				n, err := n.Write(data)
				if err != nil {
					logger.AppLogger().Warnf("@@@@ failed BLUETOOTH send, break, err:%+v, sent data:%x, data:%v", err, data, string(data))
					break
				} else {
					logger.AppLogger().Debugf("======== succ BLUETOOTH sent, n:%+v", n)
					// logger.AppLogger().Debugf("======== succ BLUETOOTH sent, n:%+v, sent data:%x", n, data)
					// logger.AppLogger().Debugf("data:%v", string(data))
				}

				// 蓝牙带宽有限，在这减缓发送速度
				// logger.AppLogger().Debugf("sleep")
				time.Sleep(time.Millisecond * time.Duration(bleSendSpendMS))
			}

			bleConnected = false
			logger.AppLogger().Debugf("HandleNotifyFunc Entry, end for loop")

		})

	return s
}
