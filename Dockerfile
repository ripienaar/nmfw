FROM golang:1.19

ARG PROTOC_VERSION=21.12
ARG GEN_GO_VERSION=latest
ARG GEN_NMFW_VERSION=latest

RUN cd /tmp && \
    apt-get update && \
    apt-get install unzip && \
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    mv /tmp/bin/protoc /bin && \
    rm -rf /tmp/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${GEN_GO_VERSION} && \
    go install github.com/ripienaar/nmfw/protoc-gen-go-nmfw@${GEN_NMFW_VERSION}
