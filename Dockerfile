FROM golang:1.11-alpine

ENV GIN_MODE release

COPY out/server server

CMD ["./server"]
