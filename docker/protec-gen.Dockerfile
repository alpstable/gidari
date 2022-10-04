FROM golang:alpine
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
ENTRYPOINT bin/protoc-gen-go