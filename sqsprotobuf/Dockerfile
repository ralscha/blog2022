FROM golang:1.22.4

ARG PROTOBUF_VERSION=21.12
ARG DEBIAN_FRONTEND=noninteractive

WORKDIR /app

RUN apt-get -qq -y update && \
    apt-get -qq -y upgrade && \
    apt-get install -y -qq unzip curl ca-certificates && \
    curl -sL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip -o protoc.zip && \		
   	unzip -j protoc.zip "bin/protoc" -d . && \
   	rm protoc.zip && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

ENTRYPOINT ["/app/protoc"]