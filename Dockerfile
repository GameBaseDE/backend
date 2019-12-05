FROM golang:latest
ENV GIN_MODE release

COPY out/server server

CMD ["./server"]
