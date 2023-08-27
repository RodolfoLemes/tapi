FROM golang:1.19-alpine AS base
WORKDIR /app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0

# System dependencies
RUN apk update \
  && apk add --no-cache \
  ca-certificates \
  git \
  && update-ca-certificates

### Executable builder
FROM base AS builder
WORKDIR /app

# Application dependencies
COPY . /app
RUN go mod download \
  && go mod verify

RUN go build -o server -a .

### Production
FROM alpine:latest

RUN apk update \
  && apk add --no-cache \
  ca-certificates \
  curl \
  tzdata \
  && update-ca-certificates

# Copy executable
COPY --from=builder /app/server /usr/local/bin/server
COPY --from=builder /app/.env /usr/local/bin
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/server"]