FROM golang:1.9

ENV APP_PATH=github.com/placer14/ob-crawler

WORKDIR /go/src/$APP_PATH

COPY ./src .

ENTRYPOINT ["/bin/bash"]
