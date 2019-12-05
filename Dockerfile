FROM golang:latest

EXPOSE 80
COPY out/server server

CMD ["./server"]
