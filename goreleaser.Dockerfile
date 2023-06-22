FROM gcr.io/distroless/static:latest
COPY k8s-wait-for-multi /
ENTRYPOINT ["/k8s-wait-for-multi"]