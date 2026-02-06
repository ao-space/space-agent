# Platform Dependency Analysis

English | [简体中文](./PLATFORM_DEPENDENCIES_CN.md)

This document summarizes where the server (space-agent) depends on the platform, and how to run without the platform using the new `PlatformEnabled` switch. Test processes write logs to the system temp directory by default (override with `AOSPACE_LOG_DIR`).

## Scope
space-agent is the server component of the product. Platform features are used for registration, Internet tunnel setup, upgrade checks, and platform switching. These are now gated by `config.Config.PlatformEnabled` (default `false`).

## Platform-dependent modules and code paths

### 1) Device registration / pairing
- `biz/service/pair/register_box.go`
  - `ServiceRegisterBox()` and `GetDeviceRegKey()` call platform registration APIs.
  - When `PlatformEnabled=false`, registration is skipped (no-op) and `GetDeviceRegKey` returns an error.
- `biz/service/pair/pairing.go`
  - Calls `ServiceRegisterBox()` during pairing. With platform disabled, the registration step is skipped.

### 2) Internet tunnel configuration
- `biz/service/bind/internet/service/config/post_config.go`
  - Enabling Internet access triggers platform registration and gateway platform switch.
  - If `PlatformEnabled=false` and `EnableInternetAccess=true`, returns `AG-403` (unsupported).

### 3) Bind/space creation (Internet access flow)
- `biz/service/bind/space/create/create.go`
  - Internet-access flow calls `registerDevice()` which triggers platform registration.
  - If `PlatformEnabled=false` and `EnableInternetAccess=true`, returns `AG-403` (unsupported).

### 4) Platform ability & API gating
- `biz/service/platform/ability.go`
  - Fetches platform ability list and checks availability.
  - When disabled, returns nil/false and skips network calls.

### 5) Upgrade checks
- `biz/service/upgrade/core.go`
  - `CheckLatestVersion()` calls platform API.
  - When disabled, returns an error immediately.
- `main.go`
  - `platform.InitPlatformAbility`, `upgrade.CronForUpgrade`, and `upgrade.CheckUpgradeSucc` are skipped when disabled.

### 6) Switch platform
- `biz/service/switch-platform/switch_platform.go`
  - Platform switch is not available when `PlatformEnabled=false` and returns `AG-403`.

### 7) Network reachability checks
- `biz/alivechecker/entry.go`
  - Platform host checks are skipped when disabled; third-party network check remains.

### 9) Docker environment variables
- `biz/docker/env.go`
  - Writes platform-related env vars into containers (API base, web URL, etc.).
  - No functional change; these may be unused when platform is disabled.

## Running without the platform (server + client only)

1. Set `PlatformEnabled=false` in configuration or runtime config file.
2. Use LAN-only flows (do not enable Internet access during bind/create).
3. Avoid platform-only APIs (switch platform, upgrade checks). These return `AG-403` when disabled.

## Limitations when platform is disabled
- Internet tunnel enablement is blocked (`/bind/internet/service/config` returns `AG-403`).
- Tryout code verification is blocked.
- Upgrade checks from platform are blocked.
- Switch platform APIs are blocked.

## Tests added for platform-optional behavior
- `biz/service/pair/register_box_test.go`
- `biz/service/bind/internet/service/config/post_config_test.go`

## Notes
- Some components (e.g., gateway account creation) still run without platform. If any downstream service requires platform-only data (like `boxRegKey`), you may need to adjust those services accordingly.
