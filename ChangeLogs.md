<!--
 * @Author: wenchao
 * @Date: 2021-12-03 08:54:59
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-03 11:12:07
 * @Description: 
-->
# System Agent Change Logs

### 0.6.0-alpha.6

- 优化配对初始化后启动速度,大约17秒。改成了配对前就启动容器了。
- 升级配置从 system-agent.yml 剥离出来。
- system -version 会导致 system-agent.yml 丢失的问题修复。

