FROM golang:1.19-bullseye AS builder

WORKDIR /usr/src/bothoi

RUN mkdir app

RUN apt-get update && apt-get install -y curl
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -tags prod -v -o app ./...

FROM debian:bullseye

WORKDIR /usr/local/bin/

RUN apt-get update && apt-get install -y ffmpeg python locales ca-certificates
RUN apt-get clean

COPY --from=builder /usr/local/bin/youtube-dl /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

RUN sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen && \
    locale-gen
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

COPY --from=builder /usr/src/bothoi/app ./

CMD ["bothoi"]
