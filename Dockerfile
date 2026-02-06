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

FROM node:22.2.0-alpine3.19 as npm-builder

WORKDIR /work/

COPY ./web/boxdocker ./web/boxdocker

RUN	apk add --no-cache zip
RUN cd web/boxdocker && npm install && npm run build && mv dist boxdocker && \
        zip -r static_html.zip boxdocker && mkdir ../../res && mv static_html.zip ../../res

FROM golang:1.22.3-alpine3.20 as golang-builder

WORKDIR /work/

COPY . .

COPY --from=npm-builder /work/res /work/res

RUN apk add --no-cache make
RUN go env -w GO111MODULE=on && make -f Makefile

FROM debian:12

ENV LANG C.UTF-8
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		netbase \
		tzdata \
        supervisor \
		iputils-ping \
		docker-compose \
		curl \
		cron \
    ; \
	apt remove docker.io -y ; \
	rm -rf /var/lib/apt/lists/*

COPY --from=golang-builder /work/build/aospace /usr/local/bin/aospace
COPY --from=golang-builder /work/supervisord.conf /etc/supervisor/supervisord.conf

EXPOSE 5678

HEALTHCHECK --interval=60s --timeout=15s CMD curl -XGET http://localhost:5678/agent/status

CMD ["/usr/bin/supervisord","-n", "-c", "/etc/supervisor/supervisord.conf"]
