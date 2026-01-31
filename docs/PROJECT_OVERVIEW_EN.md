# space-agent Project Overview

English | [简体中文](./PROJECT_OVERVIEW.md)

> Generated from the current repository (main references: `README.md`, `main.go`, `biz/web/routers/*`, `docs/swagger.yaml`, `config/config.go`, `Dockerfile`, `Makefile`, `script/*`).

## 1. Background and Purpose

space-agent is the core entry service for AO.space (open-source) all-in-one deployment. It handles device initialization and binding, key exchange, microservice startup/management, DID generation/management, and system upgrades. It provides a unified entry point for the server and exposes HTTP APIs for client/gateway pairing and management. See `README.md` for details.

## 2. System Architecture (Modules and Responsibilities)

Overall structure: Go monolith + Docker microservice orchestration + web static assets. Main flow:

- **Entry**: `main.go`
  - Initialize config and logging
  - Initialize device info/keys/client info
  - Start Web API (external + internal)
  - Start Docker microservices (or data migration)
  - Start alive checker, upgrade checks, log dir monitor, DID/LevelDB
- **Web API**: `biz/web/*`
  - Gin-based external/internal APIs, routes in `biz/web/routers/router.go`
  - Swagger docs in `docs/`
- **Business services**: `biz/service/*`
  - Pairing/binding/key exchange, DID docs, network config, system control, upgrade
- **Models & utils**: `biz/model/*`, `utils/*`
  - Device info, DID storage, network/hardware/container utilities
- **Docker management**: `biz/docker/*` + `utils/docker/*` + `res/*`
  - Compose templates and container lifecycle
- **Web static assets**: `web/boxdocker`
  - Built and packed into `res/static_html.zip`

## 3. Directory Structure (Key Paths)

- `main.go`: service entrypoint.
- `config/`: config definitions and defaults (`config/config.go`).
- `biz/`
  - `biz/web/`: HTTP API entry, routes, handlers.
  - `biz/service/`: core business logic (bind/pair/upgrade/system/DID).
  - `biz/docker/`: container orchestration and compose handling.
  - `biz/model/`: business models and DTOs.
  - `biz/alivechecker/`: container/network alive checks.
  - `biz/disk_space_monitor/`: log dir/disk monitoring.
  - `biz/notification/`: upgrade/ability change notifications.
- `utils/`: shared utilities (network, device ID, JWT, Docker, hardware, BLE, retry, etc.).
- `res/`: embedded resources (compose templates, static site zip, router config).
- `docs/`: Swagger docs (`docs/swagger.yaml`, `docs/swagger.json`).
- `web/boxdocker/`: front-end app (Vite + Vue).
- `script/`: build/run scripts.
- `Dockerfile`: multi-stage build for backend and frontend.
- `Makefile`: Go build settings.

## 4. Key Modules and Core Logic

- **Boot flow (`main.go`)**
  - Init logs and version
  - Init device info/keys and client info
  - Start Web APIs (external :5678, internal :5680)
  - Start Docker microservices or storage migration
  - Start alive checker, upgrade checks, log dir monitor, DID LevelDB

- **HTTP API (`biz/web/*`)**
  - External routes: `/agent` and `/agent/v1/api/*` for client pairing/binding/network/DID.
  - Internal routes: `/agent/v1/api/*` for internal/container/gateway calls (upgrade/device info).

- **Binding/Pairing (`biz/service/pair`, `biz/service/bind`)**
  - Key APIs: `initial`, `pairing`, `keyexchange`, `setpassword`, `bind/space/create`.

- **Docker management (`biz/docker/*`, `utils/docker/*`, `res/*`)**
  - Load embedded compose templates and start/maintain microservices.

- **Upgrade & system control (`biz/service/upgrade`, `biz/service/system`)**
  - Download/install/status APIs and reboot/shutdown.

- **DID docs & certificates (`biz/service/did/*`, `biz/model/did/*`)**
  - DID document generation/storage/update and LAN cert retrieval.

- **Alive checker (`biz/alivechecker/*`)**
  - Periodic container/network connectivity checks.

## 5. API List (Path, Method, Params, Response)

APIs are based on `docs/swagger.yaml` and route definitions.

### Info
- `GET /agent/status`: service status (`dto.BaseRspStr` + `status.Status`).
- `GET /agent/info`: service info (`dto.BaseRspStr` + `status.Info`).
- `GET /agent/logs`: logs (debug mode only).

### Pair / Bind / Keys
- `POST /agent/v1/api/initial`: init progress (`pair.PasswordInfo` -> `call.MicroServerRsp`).
- `POST /agent/v1/api/pairing`: pair client (`pair.PairingReq` -> `call.MicroServerRsp`).
- `POST /agent/v1/api/pubkeyexchange`: public key exchange.
- `POST /agent/v1/api/keyexchange`: symmetric key exchange.
- `POST /agent/v1/api/setpassword`: set admin password.
- `POST /agent/v1/api/bind/com/start`: start microservice containers.
- `GET /agent/v1/api/bind/com/progress`: startup progress.
- `GET /agent/v1/api/bind/init`: pre-bind init info.
- `POST /agent/v1/api/bind/space/create`: create space / finish binding.
- `POST /agent/v1/api/bind/password/verify`: verify admin password.
- `POST /agent/v1/api/bind/revoke`: unbind.
- `POST /agent/v1/api/admin/revoke`: admin unbind.
- `GET /agent/v1/api/pair/init`: wired binding init.
- `GET /agent/v1/api/pair/net/localips`: local IPs.
- `GET /agent/v1/api/pair/net/netconfig`: Wi-Fi list.
- `POST /agent/v1/api/bind/internet/service/config`: set Internet tunnel.
- `GET /agent/v1/api/bind/internet/service/config`: get Internet tunnel config.

### Device / Network / Cert
- `GET /agent/v1/api/device/ability`
- `GET /agent/v1/api/device/info`
- `GET /agent/v1/api/device/version`
- `GET /agent/v1/api/device/localips`
- `GET /agent/v1/api/device/netconfig`
- `POST /agent/v1/api/network/config`
- `GET /agent/v1/api/network/config`
- `POST /agent/v1/api/network/ignore`
- `GET /agent/v1/api/cert/get`

### DID
- `GET /agent/v1/api/did/document`
- `PUT /agent/v1/api/did/document/password`
- `PUT /agent/v1/api/did/document/method`

### Gateway passthrough
- `POST /agent/v1/api/passthrough`

### Upgrade / System / Switch platform
- `GET/POST /agent/v1/api/upgrade/*`
- `POST /agent/v1/api/system/reboot`
- `POST /agent/v1/api/system/shutdown`
- `POST /agent/v1/api/switch`
- `GET /agent/v1/api/switch/status`

### Other
- `GET /agent/v1/api/space/ready/check`

## 6. Local Dev, Tests, Deployment

### Development
- Go: `1.18+` (Dockerfile uses `1.20.6`).
- Frontend: `web/boxdocker` (`npm install && npm run build`).
- Build: `make -f Makefile`.
- Swagger: `swag init -g biz/web/http_server.go`.

### Tests
- Some unit tests exist (e.g., DID/JWT). Run `go test ./...` with environment readiness.
- Full test run may require Docker, nmcli, dnf, hardware files, etc. Consider skipping or injecting dependencies.

### Deployment
- Docker build: `docker build -t local/space-agent:{tag} .`
- Health check: `GET http://localhost:5678/agent/status`.

## 7. Notes and Extension Suggestions

- **Platform optionality**: see `docs/PLATFORM_DEPENDENCIES.md` (default `PlatformEnabled=false` for server+client only).
- **Single-container mode**: `AOSPACE_SINGLE_DOCKER_MODE` affects startup.
- **Internal API address**: default internal :5680 or 172.17.0.1:5680.
- **Static assets**: built into `res/static_html.zip`.
- **Upgrade**: depends on upgrade settings in `/etc/ao-space/upgrade/settings.json`.
- **Test log path**: test processes write logs to the system temp directory by default (override with `AOSPACE_LOG_DIR`).
