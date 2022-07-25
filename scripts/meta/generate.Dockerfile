#Deriving the latest base image
FROM ruby:latest

WORKDIR /usr/src/meta

# copy the /go partition from the web_go_generate container, this gives us
# access to goimports
COPY --from=web_goimports /usr/local/go /usr/local/go
COPY --from=web_goimports /go /go
COPY Gemfile ./
COPY Gemfile.lock ./
COPY test ./

# we have to actually set the go path and root in order to use the goimports
# successfully
ENV GOPATH=/go
ENV GOROOT=/usr/local/go/bin
ENV PATH="${GOROOT}:${GOPATH}:${PATH}"

RUN bundle install
