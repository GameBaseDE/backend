FROM golang:latest

EXPOSE 80
COPY out/server server

RUN echo $KUBECONFIG > /root/.kube/config
CMD ["./server"]
