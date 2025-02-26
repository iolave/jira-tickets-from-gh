FROM golang:1.24rc1-bookworm

WORKDIR /home/app

ADD ./entrypoint.sh .
ENTRYPOINT sh ./entrypoint.sh

ADD ./go.mod ./go.sum .
RUN go mod download

ADD ./cmd cmd
ADD ./internal internal
ADD ./config.yml .

