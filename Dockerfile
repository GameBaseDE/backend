FROM alpine:latest
RUN mkdir -p /root/.kube

EXPOSE 80
COPY out/server server
COPY gameservertemplates/ gameservertemplates
ENTRYPOINT ["./server"]
