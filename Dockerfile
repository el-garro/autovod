FROM golang:alpine AS builder

RUN apk add -U --no-cache ca-certificates

WORKDIR /build
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o autovod

FROM alpine

RUN wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /bin/yt-dlp
RUN chmod a+rx ~/.local/bin/yt-dlp

RUN apk update
RUN apk add --no-cache ffmpeg

WORKDIR /
COPY --from=builder /build/autovod .
CMD ["/autovod"]