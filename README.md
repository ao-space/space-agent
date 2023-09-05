# space-agent

English | [简体中文](./README_cn.md)

## Introduce

space-agent, as the carrier of AO.space (open source version) all-in-one, mainly provides a unified entrance for AO.space server to start.

space-agent is responsible for binding AO.space server and client, launching AO.space microservices and unified management.

## Feature

- Device Discovery and Binding
- Device Initialization
- Key Exchange
- Microservice Initiation and Management
- Decentralized Identity (DID) generation and management
- AO.space service upgrade and other functions

## Build

We will build a docker image of AO.space (Open Source) related services and make it available for download.

If you want to try to compile and build locally yourself

you can use our [Dockerfile](. /Dockerfile) to compile and build the container image

### Prepare Environment

- docker (>=18.09)
- git
- golang 1.18 +

### Download source code

```shell
git clone git@github.com:ao-space/space-agent.git
```

### Build image

go into project root directory , run this command

```shell
docker build -t local/space-agent:{tag} . 
````

The tag parameter can be modified to be consistent with the docker-compose.yml that is running on the server as a whole.

## Run

### Arch

- X86_64
- Arm64

### OS

- Linux:
  - EulixOS/OpenEuler
  - Ubuntu
  - Others
- Windows
- MacOS

### Docker

- Docker Engine >= 18.09
- Docker Desktop

### Recommended Hardware Configuration

- RAM： 4G
- CPU： 2 cores

### special hint

You can also run AO.space (open source version) on some development boards such as *Raspberry Pi* etc.

### Getting start

After ensuring that docker is properly installed and running in your environment

To check if docker is running correctly you can use the following command:

```shell
docker version
```

#### run container

*$AOSPACE_HOME_DIR* is  the directory where you want the data to be stored

You can replace it yourself at startup

- Linux

```shell
        sudo docker run -d --name aospace-all-in-one  \
        --restart always  \
        --network=ao-space  \
        --publish 5678:5678  \
        --publish 127.0.0.1:5680:5680  \
        -v $AOSPACE_HOME_DIR:/aospace  \
        -v /var/run/docker.sock:/var/run/docker.sock:ro  \
        -e AOSPACE_DATADIR=$AOSPACE_HOME_DIR \
        -e RUN_NETWORK_MODE="host"  \
        hub.eulix.xyz/ao-space/space-agent:dev
```

if you want to run ao.space on other os, refer to [ao.space self-hosting doc](https://ao.space/docs/install-opensource-linux)

## Notes

### get aospace agent version

```shell
system-agent -version
```

### Query Status

To see if the aospace agent service is working, call the following interface.

`/agent/status`

```shell
{
  "status": "OK",
  "version": "dev"
}
```

### how to enable swagger doc service?

modify `/opt/tmp/system-agent.yml` , change `debugmode: false` to `debugmode: true` ，**and restart agent**。

restart command

```shell
docker restart aospace-all-in-one
```

Open the swagger interface by accessing the address `http://{your-host-ip}:5678/swagger/index.html` in your computer's browser,
where the ip address is the LAN address of your box.

## Contribution Guidelines

Contributions to this project are very welcome. Here are some guidelines and suggestions to help you get involved in the project.

[Contribution Guidelines](https://github.com/ao-space/ao.space/blob/dev/docs/en/contribution-guidelines.md)

## Contact us

- Email: <developer@ao.space>
- [Official Website](https://ao.space)
- [Discussion group](https://slack.ao.space)

## Thanks for your contribution

Finally, thank you for your contribution to this project. We welcome contributions in all forms, including but not limited to code contributions, issue reports, feature requests, documentation writing, etc. We believe that with your help, this project will become more perfect and stronger.
