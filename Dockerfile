# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go

RUN  CGO_ENABLED=0 GOOS=linux go build -a -o redis-test main.go


FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/redis-test .
USER 65532:65532
ENTRYPOINT ["/redis-test"]