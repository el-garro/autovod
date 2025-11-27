FROM golang:alpine AS builder

RUN apk add -U --no-cache ca-certificates

WORKDIR /build
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o autovod

FROM python:3-alpine

RUN apk add --no-cache ffmpeg && pip3 install --no-cache-dir yt-dlp

WORKDIR /
COPY --from=builder /build/autovod .
CMD ["/autovod"]