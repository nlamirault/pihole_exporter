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

FROM golang:alpine AS build

LABEL summary="PiHole Exporter Docker image" \
      description="Prometheus Exporter for PiHole" \
      name="nlamirault/pihole_exporter" \
      version="stable" \
      url="https://github.com/nlamirault/pihole_exporter" \
      maintainer="Nicolas Lamirault <nicolas.lamirault@gmail.com>"

RUN apk add --no-cache alpine-sdk bash
WORKDIR /go/src/github.com/nlamirault/pihole_exporter
COPY . .
RUN go build -o /app/pihole_exporter pihole_exporter.go && chmod +x /app/pihole_exporter

FROM alpine:latest
COPY --from=build /app/pihole_exporter /app/pihole_exporter
WORKDIR /app
ENTRYPOINT ["./pihole_exporter"]
CMD ["-h"]
EXPOSE 9311
