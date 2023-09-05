#!/bin/bash
# Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# c=`command -v swag |wc -l`
# if [ $c -eq 0 ]; then
#     go get -u github.com/swaggo/swag/cmd/swag
# fi
# if [ $c -eq 1 ]; then
#     swag init -g biz/web/http_server.go
# fi

swag init -g biz/web/http_server.go
ARCH=arm64
GOOS=linux
GOARCH=$ARCH
go build ldflags "-s -w -X 'main.Version=EulixSpace-RaspberryPi' -X 'main.VersionNumber=1.0.0-0'" -gcflags="-N -l" -o build/system-agent
echo "build finished"
