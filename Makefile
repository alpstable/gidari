# GC is the go compiler.
GC = go

export GO111MODULE=on

default:
	chmod +rwx scripts/*.sh
	$(GC) build -o gidari-cli cmd/gidari/cmd.go

# cli will build the cli binary.
.PHONY: cli
cli:
	(cd cmd/gidari && $(GC) build -o gidari-cli main.go && mv gidari-cli ../../)

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

	(cd cmd/gidari && $(GC) clean -testcache && go test -v -count=5 -tags=utests ./...)

# e2e runs all of the end-to-end tests locally.
.PHONY: e2e
e2e:
	chmod +rwx scripts/*.sh
	$(GC) clean -testcache
	./scripts/run-e2e-tests.sh

# repository-integration-tests runs all of the repository integration tests in a docker container.
# Each test is run 5 times to minimize flakiness.
.PHONY: repository-integration-tests
repository-integration-tests:
	chmod +rwx scripts/*.sh
	$(GC) clean -testcache
	./scripts/run-integration-tests.sh repinteg 5

# fmt runs the formatter.
.PHONY: fmt
fmt:
	gofumpt -l -w .

	(cmd cmd/gidari && gofumpt -l -w .)

# lint runs the linter.
.PHONY: lint
lint:
	golangci-lint run --fix
	golangci-lint run --config .golangci.yml

	(cd cmd/gidari && golangci-lint run --fix)
	(cd cmd/gidari && golangci-ling run --config ../../.golangci.yml)

# add-license adds the license to all the top of all the .go files.
.PHONY: add-license
add-license:
	./scripts/add-license.sh

	(cd cmd/gidari && ./scripts/add-license.sh)
