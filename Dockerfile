FROM scratch

ADD build/api /api

ENTRYPOINT ["./api"]