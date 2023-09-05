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

package docker

import (
	"github.com/asaskevich/EventBus"
)

var busDockerPowerOn EventBus.Bus

const (
	signalBusDockerPowerOn = "busDockerPowerOn"
)

func PublishDockerPowerOn() {
	busDockerPowerOn.Publish(signalBusDockerPowerOn, 1)
}

func SubscribeDockerPowerOn(handler interface{}) {
	busDockerPowerOn.Subscribe(signalBusDockerPowerOn, handler)
}

func SubscribeAsyncDockerPowerOn(handler interface{}) {
	busDockerPowerOn.SubscribeAsync(signalBusDockerPowerOn, handler, false)
}

func UnsubscribeDockerPowerOn(handler interface{}) {
	busDockerPowerOn.Unsubscribe(signalBusDockerPowerOn, handler)
}

func init() {
	busDockerPowerOn = EventBus.New()
}
