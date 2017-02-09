FROM golang

ENV PORT "8080"
ADD ./warehouse /go/bin/warehouse
ADD ./config.json /go/bin/config.json

ENTRYPOINT /go/bin/warehouse

EXPOSE 8080
