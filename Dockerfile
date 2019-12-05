FROM golang:latest

EXPOSE 8080
COPY out/server server

CMD ["./server"]
