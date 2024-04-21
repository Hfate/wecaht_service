FROM golang:1.20 as builder

WORKDIR /go/src/github.com/flipped-aurora/gin-vue-admin/server
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="549398044@qq.com"

# 安装 tzdata 包和设置时区
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /go/src/github.com/flipped-aurora/gin-vue-admin/server

COPY --from=builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/server ./
COPY --from=builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/resource ./resource/
COPY --from=builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/config.docker.yaml ./

EXPOSE 8888
ENTRYPOINT ["./server", "-c", "config.docker.yaml"]