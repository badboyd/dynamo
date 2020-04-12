FROM golang:1.13-alpine AS builder
WORKDIR /go/src/github.com/badboyd/dynamo/
COPY . /go/src/github.com/badboyd/dynamo/
RUN go build -o ./dist/server ./cmd/server/server.go

FROM alpine:3.5
RUN apk add --update ca-certificates
RUN apk add --no-cache tzdata && \
  cp -f /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
  apk del tzdata
WORKDIR /app
COPY --from=builder /go/src/github.com/badboyd/dynamo/dist/server .
ENTRYPOINT ["./server"]