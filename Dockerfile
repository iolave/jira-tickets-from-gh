FROM golang:1.24rc1-bookworm

ADD ./entrypoint.sh .
ENTRYPOINT sh ./entrypoint.sh

ADD ./go.mod ./go.sum .
ADD ./cmd cmd
ADD ./internal internal
ADD ./config.yml .

RUN mkdir data
