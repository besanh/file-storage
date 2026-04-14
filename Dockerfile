FROM --platform=$BUILDPLATFORM golang:1.26.1 AS builder

ARG TARGETOS
ARG TARGETARCH

COPY . /src
WORKDIR /src

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates  \
    netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
COPY configs /app/configs
COPY cert /app/cert

WORKDIR /app

EXPOSE 8001
EXPOSE 9001
CMD ["./server", "-conf", "configs"]
