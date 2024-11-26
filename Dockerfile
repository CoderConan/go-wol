FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache tzdata
COPY go.mod go.sum ./
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
COPY . .
RUN go build -o go-wol .

FROM alpine:3.19
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/go-wol .
CMD ["./go-wol"]
