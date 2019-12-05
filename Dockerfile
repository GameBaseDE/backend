FROM alpine:3.5

ENV GIN_MODE release

RUN out/server
