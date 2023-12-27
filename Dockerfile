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

FROM debian:experimental as builder

WORKDIR /work/

COPY . .

RUN apt-get update && apt-get install golang-go npm nodejs zip -y
RUN cd web/boxdocker && npm update && npm install && npm run build && mv dist boxdocker && \
        zip -r static_html.zip boxdocker && mv static_html.zip ../../res && cd ../../
RUN go env -w GO111MODULE=on && make -f Makefile

FROM xfan1024/openeuler:23.03-light

ENV LANG C.UTF-8
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN set -eux; \
        yum -y update; \
        yum install -y \
                ca-certificates \
                net-tools \
                tzdata \
                supervisor \
                iputils \
                docker-compose \
                curl \
                cronie \
        ; \
        yum remove docker-* -y ; \
        yum clean all;

COPY --from=builder /work/build/system-agent /usr/local/bin/system-agent
COPY --from=builder /work/supervisord.conf /etc/supervisor/supervisord.conf

EXPOSE 5678

HEALTHCHECK --interval=60s --timeout=15s CMD curl -XGET http://localhost:5678/agent/status

CMD ["/usr/bin/supervisord","-n", "-c", "/etc/supervisor/supervisord.conf"]
