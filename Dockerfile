# Create image from
FROM alpine:latest
RUN mkdir -p /root/.kube

# Final image
FROM scratch

EXPOSE 80
COPY out/server server
ENTRYPOINT ["server"]
