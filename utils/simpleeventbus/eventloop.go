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

package simpleeventbus

import (
	"sync"
)

type EvFunc func(event string)

type EventLoop struct {
	chEvent chan bool

	evList    []string
	lock      sync.Mutex
	evHandler map[string]EvFunc
}

func NewEventLoop() *EventLoop {
	return &EventLoop{chEvent: make(chan bool, 128),
		evList:    make([]string, 0),
		evHandler: make(map[string]EvFunc)}
}

func (loop *EventLoop) RegisterEvent(event string, handler EvFunc) {
	loop.lock.Lock()
	defer loop.lock.Unlock()
	loop.evHandler[event] = handler
}

func (loop *EventLoop) UnregisterEvent(event string) {
	loop.lock.Lock()
	defer loop.lock.Unlock()
	delete(loop.evHandler, event)
}

func (loop *EventLoop) Poll() {
	for range loop.chEvent {
		loop.disposeEvent()
	}
}

func (loop *EventLoop) disposeEvent() {
	loop.lock.Lock()
	defer loop.lock.Unlock()

	for i := 0; i < len(loop.evList); i++ {
		event := loop.evList[i]
		handler, exist := loop.evHandler[event]
		if exist {
			handler(event)
		}
		loop.evList = append(loop.evList[:i], loop.evList[i+1:]...)
		i--
	}
}

func (loop *EventLoop) PostEvent(event string) {
	loop.lock.Lock()
	defer loop.lock.Unlock()

	loop.evList = append(loop.evList, event)
	loop.chEvent <- true
}
