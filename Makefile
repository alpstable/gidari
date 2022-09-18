PKGS=$(shell scripts/list_pkgs.sh ./pkg)

default:
	go build cmd/gidari.go

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	chmod +rwx scripts/*.sh
	chmod +rwx third_party/docker/rs-init.sh

	scripts/build-storage.sh

# proto is a phony target that will generate the protobuf files.
.PHONY: proto
proto:
	protoc --proto_path=pkg/proto --go_out=pkg/proto pkg/proto/db.proto

# test runs all of the application tests locally.
.PHONY: test
test:
	go test ./... -v

# ctests are the integration tests in CI/CD.
.PHONY: ctests
ctests:
	chmod +rwx scripts/*.sh
	./scripts/run-ctests.sh
