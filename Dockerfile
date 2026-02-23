FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /owget .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /owget /usr/local/bin/owget
ENTRYPOINT ["owget"]
