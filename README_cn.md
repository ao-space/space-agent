# space-agent

[English](./README.md) | 简体中文

## 简介

space-agent 作为 AO.space（开源版）一体机的载体，主要为 AO.space 服务器提供统一的启动入口。

space-agent 负责绑定 AO.space 服务器和客户端，启动 AO.space 微服务并进行统一管理。

## 功能介绍

- 设备扫描与绑定
- 设备初始化
- 微服务启动与管理
- 分布式数字身份（DID）生成与管理
- 傲空间服务升级等功能

## 构建

我们会将傲空间（开源版）相关服务构建成容器镜像提供用户下载

如果你希望自己尝试在本地编译与构建

可以使用我们的 [Dockerfile](./Dockerfile) 来编译构建容器镜像

### 环境准备

- docker (>=18.09)
- git
- golang 1.18 +

### 源码下载

```shell
git clone git@github.com:ao-space/space-agent.git
```

### 容器镜像构建

进入模块根目录，执行命令

```shell
docker build -t local/space-agent:{tag} . 
````

其中 tag 参数可以根据实际情况修改，和服务器整体运行的 docker-compose.yml 保持一致即可。

## 运行

### Arch

- X86_64
- Arm64

### OS

- Linux:
  - EulixOS/OpenEuler
  - Ubuntu
  - Other(尚未验证)
- Windows
- MacOS

### Docker

- Docker Engine >= 18.09
- Docker Desktop

### 推荐运行配置

- RAM： 4G
- CPU： 2核

### 特别提示

您也可以在一些开发板上运行傲空间（开源版），例如*树莓派*等

### 快速开始

在保证你的环境中已正确安装并运行了docker

检查docker是否正确运行可以使用如下命令：

```shell
docker version
```

启动容器

*$AOSPACE_HOME_DIR* 表示你希望将数据存储在个目录

启动时可以自行替换

- Linux 环境

```shell
DATADIR="$HOME/aospace"
sudo docker run -d --name aospace-all-in-one  \
        --restart always  \
        --network=ao-space  \
        --publish 5678:5678  \
        --publish 127.0.0.1:5680:5680  \
        -v $DATADIR:/aospace  \
        -v /var/run/docker.sock:/var/run/docker.sock:ro  \
        -e AOSPACE_DATADIR=$DATADIR \
        -e RUN_NETWORK_MODE="host"  \
        ghcr.io/ao-space/space-agent:latest
```

其他环境的启动，可以参考[傲空间私有部署](https://ao.space/open/documentation/105001)

## 注意事项

### 查看aospace-agent 内置版本号

进入aosapce-agent 容器内执行

```shell
docker exec -it aospace-all-in-one system-agent --version
```

### 查看状态(可选)

查看盒子的 system-agent服务是否正常, 可调用如下接口，

`/agent/status`

```shell
{
  "status": "OK",
  "version": "dev"
}
```

### 如何开启 agent 的swagger ?

编辑 `/opt/tmp/system-agent.yml` 文件, 把 `debugmode: false` 改成 `debugmode: true` ，**然后重启 agent**。

重启 aospace-agent 容器可以使用以下命令

```shell
docker restart aospace-all-in-one
```

在电脑浏览器访问地址 `http://192.168.124.11:5678/swagger/index.html` 打开 swagger 界面，其中的 ip 地址是你盒子的局域网地址。

## 贡献指南

我们非常欢迎对本项目进行贡献。以下是一些指导原则和建议，希望能够帮助您参与到项目中来。

[贡献指南](https://github.com/ao-space/ao.space/blob/dev/docs/cn/contribution-guidelines.md)

## 联系我们

- 邮箱：<developer@ao.space>
- [官方网站](https://ao.space)
- [讨论组](https://slack.ao.space)

## 感谢您的贡献

最后，感谢您对本项目的贡献。我们欢迎各种形式的贡献，包括但不限于代码贡献、问题报告、功能请求、文档编写等。我们相信在您的帮助下，本项目会变得更加完善和强大。
