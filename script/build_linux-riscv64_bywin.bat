rem Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
rem
rem Licensed under the Apache License, Version 2.0 (the "License");
rem you may not use this file except in compliance with the License.
rem You may obtain a copy of the License at
rem
rem     http://www.apache.org/licenses/LICENSE-2.0
rem
rem Unless required by applicable law or agreed to in writing, software
rem distributed under the License is distributed on an "AS IS" BASIS,
rem WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
rem See the License for the specific language governing permissions and
rem limitations under the License.


@REM go get -u github.com/swaggo/swag/cmd/swag  @REM 下载的版本较新，CI 机器 go 版本 go 1.15.7-3，可能会 CI 编译失败!!!!

@REM swag init -g biz/web/routers/router.go @REM 去掉本行注释，可以重新生成 swagger 文档。

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=riscv64

go build  -ldflags "-s -w -X 'main.Version=aospace-agent-opensource' -X 'main.VersionNumber=0.9.0'" -gcflags="-N -l" -o build/system-agent
echo "build finished"
