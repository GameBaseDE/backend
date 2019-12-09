FROM golang:latest

RUN env
COPY config /root/.kube/config

EXPOSE 80
COPY out/server server

CMD ["./server"]
