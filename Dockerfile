FROM golang

ENV PORT "8080"
ADD ./express-inventory /go/bin/expresso-inventory
ADD ./config.json /go/bin/config.json

ENTRYPOINT /go/bin/expresso-inventory

EXPOSE 8080
