FROM golang:latest

RUN env
COPY ~/.kube /root/.kube

EXPOSE 80
COPY out/server server

CMD ["./server"]
