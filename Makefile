# GC is the go compiler.
GC = go

export GO111MODULE=on

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	scripts/build-storage.sh

# proto is a phony target that will generate the protobuf files.
.PHONY: proto
proto:
	protoc --proto_path=proto --go_out=proto proto/db.proto

# test runs all of the unit tests locally. Each test is run 5 times to minimize flakiness.
.PHONY: tests
tests:
	$(GC) clean -testcache && go test -v -count=5 -tags=utests ./...

# e2e runs all of the end-to-end tests locally.
.PHONY: e2e
e2e:
	chmod +rwx scripts/*.sh
	$(GC) clean -testcache
	./scripts/run-e2e-tests.sh

# fmt runs the formatter.
.PHONY: fmt
fmt:
	gofumpt -l -w .

# lint runs the linter.
.PHONY: lint
lint:
	golangci-lint run --fix
	golangci-lint run --config .golangci.yml

# add-license adds the license to all the top of all the .go files.
.PHONY: add-license
add-license:
	./scripts/add-license.sh
