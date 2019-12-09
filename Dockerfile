FROM golang:latest

EXPOSE 80
COPY out/server server

RUN mkdir /root/.kube && echo $KUBECONFIG > /root/.kube/config
CMD ["./server"]
