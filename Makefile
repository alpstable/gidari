default:
	go build cmd/gidari.go

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	chmod +rwx scripts/*.sh

	scripts/build-storage.sh

# hosts sets up the hosts file for integration tests. This only needs to be run once.
.PHONY: hosts
hosts:
	scripts/build-hosts.sh

# proto is a phony target that will generate the protobuf files.
.PHONY: proto
proto:
	protoc --proto_path=pkg/proto --go_out=pkg/proto pkg/proto/db.proto

# test runs all of the application tests locally.
.PHONY: tests
tests:
	godotenv -f /etc/alpine-hodler/auth.env go test ./... -v

# ci are the integration tests in CI/CD.
.PHONY: ci
ci:
	chmod +rwx scripts/*.sh
	./scripts/run-ci-tests.sh
