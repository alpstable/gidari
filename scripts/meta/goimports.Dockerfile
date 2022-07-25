# we have to use latest over alpine because we need to be able to copy the
# go binaries to the ruby generate dockerfile, and alpine doesn't fully support
# that.
FROM golang:latest

# install goimports so that we can format go files and insert the imported
# dependencies programatically
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/99designs/gqlgen@latest
