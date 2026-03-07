FROM golang:1.25.3-alpine AS builder

WORKDIR /src
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o /url-shortener ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /url-shortener /usr/local/bin/url-shortener
COPY --from=builder /src/index.html /usr/local/bin/index.html

EXPOSE 8080
USER app

ENTRYPOINT ["/usr/local/bin/url-shortener"]
