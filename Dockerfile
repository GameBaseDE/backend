FROM golang:latest
ENV GIN_MODE release

EXPOSE 8080
COPY out/server server

CMD ["./server"]
