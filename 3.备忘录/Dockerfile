FROM golang:1.20 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags="-w -s" -o ./main
RUN mkdir -p /app/publish  \
    && cp main /app/publish  \
    && cp -r conf.ini /app/publish

FROM busybox

WORKDIR /app

COPY --from=builder /app/publish .

# 指定运行时环境变量
ENV GIN_MODE=debug
EXPOSE 8080

ENTRYPOINT ["./main"]