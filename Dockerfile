FROM golang:1.9

RUN mkdir -p /app

WORKDIR /app

ADD bin/main /app/main

ENTRYPOINT ["./main"]