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


@REM set host=192.168.2.33
set host=%1
set port=%2
set exe=system-agent
@REM script/build_linux-arm64_bywin.bat 
ssh -p %port% root@%host% systemctl stop system-agent 
scp -P %port% build/%exe% root@%host%:/usr/local/bin/
ssh -p %port% root@%host% chmod +x /usr/local/bin/%exe%
ssh -p %port% root@%host% rm -rf /opt/logs/system-agent
ssh -p %port% root@%host% systemctl start system-agent
@REM ssh -p %port% root@%host% /usr/local/bin/%exe%
@REM ssh -p %port% root@%host% tail -f /opt/logs/system-agent/system-agent.log
