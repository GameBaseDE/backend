FROM golang:latest

RUN mkdir -p /root/.kube

EXPOSE 80
COPY out/server server
CMD ["./server"]
