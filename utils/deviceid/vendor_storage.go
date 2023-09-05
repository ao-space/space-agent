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

package deviceid

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	VendorStorageFile = "/dev/vendor_storage"

	VENDOR_SN_ID        = 1
	VENDOR_WIFI_MAC_ID  = 2
	VENDOR_LAN_MAC_ID   = 3
	VENDOR_BLUETOOTH_ID = 4
	VENDOR_USER_NAME1   = 16 // 设备内部型号, 8位字符, 不足时前面补空格.
	VENDOR_USER_NAME2   = 17

	VENDOR_REQ_TAG     = 0x56524551
	VENDOR_USER_LENGTH = 32

	// 	[root@EulixOS SN_MAC]# ./a.out
	// VENDOR_READ_IO = 40047601
	// VENDOR_WRITE_IO = 40047602
	VENDOR_READ_IO  = 0x40047601
	VENDOR_WRITE_IO = 0x40047602

	GPIOHANDLE_REQUEST_OUTPUT        = 0x2
	GPIO_GET_LINEHANDLE_IOCTL        = 0xc16cb403
	GPIOHANDLE_SET_LINE_VALUES_IOCTL = 0xc040b409
)

// func main() {

// 	// if len(os.Args) < 2 {
// 	// 	fmt.Printf("os.Args:%+v\n", os.Args)
// 	// 	return
// 	// }

// 	data, err := VendorStorageRead(VENDOR_SN_ID)
// 	if err != nil {
// 		fmt.Printf("failed ReadSn: %+v\n", err)
// 		return
// 	}
// 	fmt.Printf("data(%v): %+v\n", len(data), (string)(data))
// }

type gpiohandlerequest struct {
	Tag  uint32
	Id   uint16
	Len  uint16
	Data [128]byte
}

func VendorStorageRead(vendor_id uint16) ([]byte, error) {
	//file := "/dev/vendor_storage"

	f, err := os.Open(VendorStorageFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	handlereq := gpiohandlerequest{
		Tag: VENDOR_REQ_TAG,
		Id:  vendor_id,
		Len: VENDOR_USER_LENGTH,
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(f.Fd()), VENDOR_READ_IO, uintptr(unsafe.Pointer(&handlereq))); errno != 0 {
		return nil, fmt.Errorf("GPIO_GET_LINEHANDLE_IOCTL: %v", errno)
	}

	return handlereq.Data[:], nil
}
