# GC is the go compiler.
GC = go

export GO111MODULE=on

default:
	chmod +rwx scripts/*.sh
	$(GC) build -o gidari-cli cmd/gidari/cmd.go

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	scripts/build-storage.sh

# proto is a phony target that will generate the protobuf files.
.PHONY: proto
proto:
	protoc --proto_path=internal/proto --go_out=internal/proto internal/proto/db.proto

# test runs all of the unit tests locally. Each test is run 5 times to minimize flakiness.
.PHONY: tests
tests:
	$(GC) clean -testcache
	go test -v -count=5 -tags=utests ./...

# e2e runs all of the end-to-end tests locally.
.PHONY: e2e
e2e:
	chmod +rwx scripts/*.sh
	$(GC) clean -testcache
	./scripts/run-e2e-tests.sh

# mongo-integration-tests runs all of the mongo integration tests in a docker container. Each test is run 5 times
# to minimize flakiness.
.PHONY: mongo-integration-tests
mongo-integration-tests:
	chmod +rwx scripts/*.sh
	$(GC) clean -testcache
	./scripts/run-integration-tests.sh mdbinteg 5

# postgres-integration-tests runs all of the postgres integration tests in a docker container. Each test is run 5
# times to minimize flakiness.
.PHONY: postgres-integration-tests
postgres-integration-tests:
	chmod +rwx scripts/*.sh
	$(GC) clean -testcache
	./scripts/run-integration-tests.sh pginteg 5

# lint runs the linter.
.PHONY: lint
lint:
	scripts/lint.sh

# fmt runs the formatter.
.PHONY: fmt
fmt:
	./scripts/fmt.sh

# add-license adds the license to all the top of all the .go files.
.PHONY: add-license
add-license:
	./scripts/add-license.sh
