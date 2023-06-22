FROM gcr.io/distroless/static:nonroot
COPY k8s-wait-for-multi /
ENTRYPOINT ["/k8s-wait-for-multi"]