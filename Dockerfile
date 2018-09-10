# Copyright (C) 2016-2018 Nicolas Lamirault <nicolas.lamirault@gmail.com>

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

FROM alpine:latest


LABEL summary="PiHole Exporter Docker image" \
      description="Prometheus Exporter for PiHole" \
      name="nlamirault/pihole_exporter" \
      version="stable" \
      url="https://github.com/nlamirault/pihole_exporter" \
      maintainer="Nicolas Lamirault <nicolas.lamirault@gmail.com>"

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN apk add --no-cache \
    ca-certificates

COPY . /go/src/github.com/nlamirault/pihole_exporter

RUN set -x \
    && apk add --no-cache --virtual .build-deps \
       go \
       git \
       gcc \
       libc-dev \
       libgcc \
    && cd /go/src/github.com/nlamirault/pihole_exporter \
    && go build -o /usr/bin/pihole_exporter . \
    && apk del .build-deps \
    && rm -rf /go \
    && echo "Build complete."

ENTRYPOINT ["pihole_exporter"]
