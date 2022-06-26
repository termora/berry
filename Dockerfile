FROM golang:latest AS builder

WORKDIR /build
COPY . ./
RUN go mod download -x
ENV CGO_ENABLED 0
RUN go build -o termora -v -ldflags="-buildid= -X github.com/termora/berry/common.Version=`git rev-parse --short HEAD`" ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /build/termora termora

CMD ["/app/termora", "bot"]
