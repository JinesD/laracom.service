FROM golang:1.14.13-alpine as builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add git
RUN mkdir -p /data/laracom.user.service

WORKDIR /app/laracom.user.service

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o laracom.user.service .

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add --no-cache bash ca-certificates
RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/laracom.user.service/laracom.user.service .

CMD ["./laracom.user.service"]