# 平台依赖分析

简体中文 | [English](./PLATFORM_DEPENDENCIES.md)

本文档总结 space-agent 对平台的依赖点，以及通过 `PlatformEnabled` 开关实现平台可选运行的说明。测试进程会默认把日志写入系统临时目录（可通过环境变量 `AOSPACE_LOG_DIR` 指定）。

## 适用范围
space-agent 是产品的服务端组件。平台相关能力包括：注册、互联网通道配置、升级检查、平台切换。现在这些能力都受 `config.Config.PlatformEnabled` 控制（默认 `false`）。

## 平台依赖模块与代码路径

### 1) 设备注册 / 配对
- `biz/service/pair/register_box.go`
  - `ServiceRegisterBox()` 和 `GetDeviceRegKey()` 会请求平台注册接口。
  - 当 `PlatformEnabled=false` 时，注册流程直接跳过；`GetDeviceRegKey` 返回错误。
- `biz/service/pair/pairing.go`
  - 配对中会调用 `ServiceRegisterBox()`，平台禁用时注册步骤被跳过。

### 2) 互联网通道配置
- `biz/service/bind/internet/service/config/post_config.go`
  - 启用互联网通道会触发平台注册与网关平台切换。
  - 当 `PlatformEnabled=false` 且 `EnableInternetAccess=true` 时返回 `AG-403`。

### 3) 绑定/创建空间（互联网访问流程）
- `biz/service/bind/space/create/create.go`
  - 互联网访问流程调用 `registerDevice()` 触发平台注册。
  - 当 `PlatformEnabled=false` 且 `EnableInternetAccess=true` 时返回 `AG-403`。

### 4) 平台能力与 API 探测
- `biz/service/platform/ability.go`
  - 获取平台能力列表并检查 API 可用性。
  - 当禁用时直接返回 nil/false，不发起请求。

### 5) 升级检查
- `biz/service/upgrade/core.go`
  - `CheckLatestVersion()` 会请求平台版本接口。
  - 当禁用时立即返回错误。
- `main.go`
  - 平台能力获取、升级定时检查与升级结果检测在禁用时跳过。

### 6) 切换平台
- `biz/service/switch-platform/switch_platform.go`
  - 当 `PlatformEnabled=false` 时返回 `AG-403`。

### 7) 网络连通性检测
- `biz/alivechecker/entry.go`
  - 平台相关网络检查在禁用时跳过，第三方网络检查仍保留。

### 8) Docker 环境变量
- `biz/docker/env.go`
  - 仍会写入平台相关 env，但在平台禁用时通常不会被用到。

## 无平台运行（仅服务端 + 客户端）

1. 在配置中设置 `PlatformEnabled=false`。
2. 使用局域网绑定流程（不要启用互联网通道）。
3. 避免平台专属 API（切换平台、升级检查等）。这些会返回 `AG-403`。

## 平台禁用的限制
- 互联网通道无法启用（`/bind/internet/service/config` 返回 `AG-403`）。
- 平台升级检查不可用。
- 切换平台不可用。

## 已添加的单元测试
- `biz/service/pair/register_box_test.go`
- `biz/service/bind/internet/service/config/post_config_test.go`
- `biz/service/platform/ability_test.go`

## 备注
- 仍有部分组件（例如网关创建账号）在无平台情况下正常运行。如果下游服务需要平台返回字段（如 `boxRegKey`），需要配套调整。
