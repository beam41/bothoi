FROM golang:1.19-alpine AS builder

WORKDIR /usr/src/bothoi

RUN mkdir app

RUN apk update && apk upgrade --no-cache
RUN apk add --no-cache curl gcc
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -tags prod -v -o app ./...

FROM alpine:latest

WORKDIR /usr/local/bin/

RUN apk update && apk upgrade --no-cache
RUN apk add --no-cache ffmpeg python3 ca-certificates
RUN ln -s /usr/bin/python3 /usr/bin/python

COPY --from=builder /usr/local/bin/youtube-dl /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

COPY --from=builder /usr/src/bothoi/app ./

CMD ["bothoi"]
