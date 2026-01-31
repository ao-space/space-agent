# space-agent 项目概览

简体中文 | [English](./PROJECT_OVERVIEW_EN.md)

> 基于当前仓库代码生成（主要参考 `README.md`、`main.go`、`biz/web/routers/*`、`docs/swagger.yaml`、`config/config.go`、`Dockerfile`、`Makefile`、`script/*` 等）。

## 1. 项目背景与用途

space-agent 是 AO.space（开源版）一体化部署的核心入口服务，负责设备初始化与绑定、密钥交换、微服务启动与管理、DID 生成与管理、以及系统升级等功能。它为 AO.space 服务端提供统一的启动与运行管理入口，并提供 HTTP API 供客户端/网关进行配对与管理。项目说明见 `README.md`。

## 2. 整体系统架构（模块划分、主要职责）

整体为 Go 单体服务 + Docker 微服务编排 + Web 前端资源构建的结构，核心链路如下：

- **启动入口**：`main.go`
  - 初始化配置与日志
  - 设备信息/密钥/客户端信息初始化
  - 启动 Web API（外部 + 内部）
  - 启动 Docker 微服务（或数据迁移）
  - 启动存活检测、升级检查、日志目录监控、DID/LevelDB
- **Web API 层**：`biz/web/*`
  - Gin 框架实现外部/内部 API，路由见 `biz/web/routers/router.go`
  - Swagger 文档生成配置见 `docs/`
- **业务服务层**：`biz/service/*`
  - 配对/绑定/密钥交换、DID 文档、网络配置、系统控制、升级等
- **设备与模型层**：`biz/model/*`, `utils/*`
  - 设备信息、DID 文档存储、网络/硬件/容器等通用能力
- **Docker 管理层**：`biz/docker/*` + `utils/docker/*` + `res/*`
  - 负责微服务容器的编排、启动、升级、配置生成
- **前端静态资源**：`web/boxdocker`
  - 通过 Dockerfile 构建并打包到 `res/static_html.zip`

## 3. 目录结构说明（逐级说明关键目录和文件作用）

- `main.go`：服务入口，启动各核心模块。
- `config/`：配置定义与默认值（`config/config.go`）。
- `biz/`
  - `biz/web/`：HTTP API 入口、路由和处理器。
    - `biz/web/routers/`：外部/内部路由定义与 Server 启动逻辑。
    - `biz/web/handler/`：具体 API 处理器。
  - `biz/service/`：核心业务逻辑实现（绑定、配对、升级、系统控制、DID 等）。
  - `biz/docker/`：容器启动、compose 文件处理、容器状态监听等。
  - `biz/model/`：业务模型与 DTO，DID 与设备数据存储。
  - `biz/alivechecker/`：容器/网络存活检测。
  - `biz/disk_space_monitor/`：日志目录与磁盘监控。
  - `biz/notification/`：升级/能力变更等通知逻辑。
- `utils/`：通用能力（网络、设备 ID、JWT、Docker 接口、硬件信息、BLE、重试等）。
- `res/`：内置资源（docker-compose 模板、静态站点 zip、路由配置等）。
- `docs/`：Swagger 文档与生成入口（`docs/swagger.yaml`、`docs/swagger.json`）。
- `web/boxdocker/`：前端工程（Vite + Vue），打包生成静态资源。
- `script/`：构建与运行脚本（含 swagger 生成与交叉编译示例）。
- `Dockerfile`：多阶段构建，打包服务与前端。
- `Makefile`：本地 Go 编译配置。

## 4. 关键模块与核心逻辑说明

- **启动流程（`main.go`）**
  - 初始化日志与版本信息。
  - 初始化设备信息、设备密钥、客户端信息。
  - 按环境变量 `AOSPACE_SINGLE_DOCKER_MODE` 决定是否启动平台能力、升级任务、微服务启动逻辑。
  - 启动 Web API 服务（外部 5678，内部 5680）。
  - 启动 Docker 微服务或执行存储迁移。
  - 启动 alive checker、升级检查、日志目录监控、DID LevelDB。

- **HTTP API（`biz/web/*`）**
  - 外部路由：`/agent` 与 `/agent/v1/api/*`，用于客户端绑定、配对、网络配置、DID 管理等。
  - 内部路由：`/agent/v1/api/*`，用于容器内部或网关调用（如升级、设备信息）。

- **绑定/配对流程（`biz/service/pair`、`biz/service/bind`）**
  - 包含 `initial`、`pairing`、`keyexchange`、`setpassword`、`bind/space/create` 等关键接口。

- **Docker 容器管理（`biz/docker/*`、`utils/docker/*`、`res/*`）**
  - 读取嵌入式 compose 模板与配置，启动/维护 AO.space 微服务容器。
  - 通过 `res/load.go` 动态调整挂载路径、环境变量等。

- **升级与系统管理（`biz/service/upgrade`、`biz/service/system`）**
  - 提供下载、安装、状态查询与升级配置管理接口。
  - 支持关机、重启等系统控制操作。

- **DID 文档与证书（`biz/service/did/*`、`biz/model/did/*`）**
  - DID 文档生成、存储与更新。
  - LAN 证书获取接口。

- **存活检测（`biz/alivechecker/*`）**
  - 定期检测容器与网络连通状态。

## 5. API 接口清单（接口路径、方法、参数、返回值）

以下接口来自 `docs/swagger.yaml` 与路由定义（外部与内部共用前缀 `/agent` 或 `/agent/v1/api`）。参数与返回值类型请参考 swagger definitions。

### 信息类
- `GET /agent/status`：返回服务状态（`dto.BaseRspStr` + `status.Status`）。
- `GET /agent/info`：返回服务信息（`dto.BaseRspStr` + `status.Info`）。
- `GET /agent/logs`：返回日志（调试模式才启用）。

### 配对 / 绑定 / 密钥
- `POST /agent/v1/api/initial`：查询初始化进度（`pair.PasswordInfo` -> `call.MicroServerRsp`）。
- `POST /agent/v1/api/pairing`：客户端绑定（`pair.PairingReq` -> `call.MicroServerRsp`）。
- `POST /agent/v1/api/pubkeyexchange`：公钥交换（`pair.PubKeyExchangeReq` -> `pair.PubKeyExchangeRsp`）。
- `POST /agent/v1/api/keyexchange`：对称密钥交换（`pair.KeyExchangeReq` -> `pair.KeyExchangeRsp`）。
- `POST /agent/v1/api/setpassword`：设置管理员密码（`pair.PasswordInfo` -> `call.MicroServerRsp`）。
- `POST /agent/v1/api/bind/com/start`：启动微服务容器（无 body -> `dto.BaseRspStr`）。
- `GET /agent/v1/api/bind/com/progress`：查询启动进度（`progress.ProgressRsp`）。
- `GET /agent/v1/api/bind/init`：绑定前初始化信息（`bindinit.InitReq` -> `pair.InitResult`）。
- `POST /agent/v1/api/bind/space/create`：创建空间/完成绑定（`create.CreateReq` -> `create.CreateRsp`）。
- `POST /agent/v1/api/bind/password/verify`：验证管理员密码（`password.VerifyReq` -> `password.VerifyRsp`）。
- `POST /agent/v1/api/bind/revoke`：解除绑定（`revoke.RevokeReq` -> `revoke.RevokeRsp`）。
- `POST /agent/v1/api/admin/revoke`：管理员解绑（`pair.RevokeReq` -> `call.MicroServerRsp`）。
- `GET /agent/v1/api/pair/init`：有线绑定初始化（`pair.InitResult`）。
- `GET /agent/v1/api/pair/net/localips`：获取本地 IP（`[]pair.Network`）。
- `GET /agent/v1/api/pair/net/netconfig`：获取 Wi‑Fi 列表（`[]pair.WifiListRsp`）。
- `POST /agent/v1/api/bind/internet/service/config`：设置 Internet tunnel（`config.ConfigReq` -> `config.ConfigRsp`）。
- `GET /agent/v1/api/bind/internet/service/config`：获取 Internet tunnel 配置（query 参数 -> `config.GetConfigRsp`）。

### 设备 / 网络 / 证书
- `GET /agent/v1/api/device/ability`：设备能力（`device_ability.DeviceAbility`）。
- `GET /agent/v1/api/device/info`：设备信息（`device.StorageInfo`）。
- `GET /agent/v1/api/device/version`：设备版本（`device.BoxDeviceVersion`）。
- `GET /agent/v1/api/device/localips`：设备本地 IP（`[]pair.Network`）。
- `GET /agent/v1/api/device/netconfig`：设备 Wi‑Fi 列表（`[]pair.WifiListRsp`）。
- `POST /agent/v1/api/network/config`：配置网络（`network.NetworkConfigReq` -> `dto.BaseRspStr`）。
- `GET /agent/v1/api/network/config`：查询网络配置（`network.NetworkStatusRsp`）。
- `POST /agent/v1/api/network/ignore`：忽略网络（`network.NetworkIgnoreReq`）。
- `GET /agent/v1/api/cert/get`：获取 LAN 证书（`certificate.LanCert`）。

### DID
- `GET /agent/v1/api/did/document`：获取 DID 文档（`document.GetDocumentReq` -> `document.GetDocumentRsp`）。
- `PUT /agent/v1/api/did/document/password`：更新 DID 文档密码（`password.UpdateDocumentPasswordReq` -> `dto.BaseRspStr`）。
- `PUT /agent/v1/api/did/document/method`：更新 DID 验证方法（`method.UpdateDocumentMethodReq` -> `method.UpdateDocumentMethodRsp`）。

### 网关/透传
- `POST /agent/v1/api/passthrough`：网关透传（Header `Request-Id` + `dto.LanInvokeReq` -> `string`）。

### 升级 / 系统 / 切换平台
- `GET /agent/v1/api/upgrade/config`：获取升级配置（`upgrade.UpgradeConfig`）。
- `POST /agent/v1/api/upgrade/config`：设置升级配置（`upgrade.UpgradeConfig`）。
- `POST /agent/v1/api/upgrade/download`：下载升级包（`upgrade.StartDownRes` -> `upgrade.Task`）。
- `POST /agent/v1/api/upgrade/install`：安装升级（`upgrade.StartUpgradeRes` -> `upgrade.Task`）。
- `GET /agent/v1/api/upgrade/status`：查询升级任务状态（`upgrade.Task`）。
- `POST /agent/v1/api/system/reboot`：重启系统（`dto.BaseRspStr`）。
- `POST /agent/v1/api/system/shutdown`：关机（`dto.BaseRspStr`）。
- `POST /agent/v1/api/switch`：切换平台（`switchplatform.SwitchPlatformReq` -> `switchplatform.SwitchPlatformResp`）。
- `GET /agent/v1/api/switch/status`：查询切换状态（query `transId` -> `switchplatform.SwitchStatusQueryResp`）。

### 其他
- `GET /agent/v1/api/space/ready/check`：空间就绪检查（`space.ReadyCheckRsp`）。

## 6. 本地开发、测试与部署说明

### 本地开发
- Go 环境：`go 1.18+`（`Dockerfile` 构建使用 `go 1.20.6`）。
- 前端：`web/boxdocker` 使用 `npm install && npm run build`（Dockerfile 中集成）。
- 构建命令：`make -f Makefile` 或 `go build ...`。
- Swagger 生成：`script/build.sh` 中执行 `swag init -g biz/web/http_server.go`。

### 测试
- 项目包含部分单元测试（如 `biz/model/did/*_test.go`、`utils/jwt/*_test.go` 等）。
- 常规运行：`go test ./...`（需保证依赖和环境齐备）。
  - 全量单测建议使用 `go test ./...`，但注意部分测试/代码路径依赖系统环境（如 Docker、dnf、nmcli、硬件信息/文件路径等），在 CI 或非目标设备上可能需要跳过或进行依赖注入替换。

### 部署
- Docker 构建：`docker build -t local/space-agent:{tag} .`（见 `README.md`）。
- 运行示例（端口 `5678` 对外暴露，`5680` 为内部端口）：
  - 参考 `README.md` 中 `docker run` 示例。
- 健康检查：`GET http://localhost:5678/agent/status`（见 `Dockerfile` Healthcheck）。

## 7. 常见注意事项和扩展建议

- **调试模式**：`config.Config.DebugMode` 影响 Swagger 与日志接口暴露（见 `biz/web/routers/server.go` 与 README 中说明）。
- **平台依赖与可选性**：平台依赖点与可选模式见 `docs/PLATFORM_DEPENDENCIES_CN.md` / `docs/PLATFORM_DEPENDENCIES.md`（默认 `PlatformEnabled=false`，仅用服务端+客户端即可运行）。
- **单容器模式**：环境变量 `AOSPACE_SINGLE_DOCKER_MODE` 会改变启动流程（`main.go`）。
- **内部 API 地址**：默认内部监听 `:5680` 或 `172.17.0.1:5680`，请结合 `config.Config.Web` 与部署方式使用。
- **静态前端资源**：`web/boxdocker` 构建后打包进 `res/static_html.zip`，由服务内置提供。
- **升级流程**：升级接口依赖后台任务与配置文件（`/etc/ao-space/upgrade/settings.json`）。
- **建议扩展**：
  - 将 swagger 生成纳入 CI 或 Makefile 目标，确保接口文档与代码同步。
  - 为关键业务（绑定、升级、Docker 启动）增加更完善的端到端测试。
  - 对内部/外部 API 增加权限或鉴权策略（目前以配置与环境控制为主）。
