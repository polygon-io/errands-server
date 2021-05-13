FROM golang:1.14

COPY errands-server /go/bin/errands-server

ENTRYPOINT [ "errands-server" ]