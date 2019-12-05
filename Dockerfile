FROM alpine:3.5

ENV GIN_MODE release

COPY out/server server

RUN ./server
