FROM golang:1.20 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/k8s-wait-for-multi

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11:latest
COPY --from=build /go/bin/k8s-wait-for-multi /
ENTRYPOINT ["/k8s-wait-for-multi"]