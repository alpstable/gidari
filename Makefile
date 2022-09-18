PKGS=$(shell scripts/list_pkgs.sh ./pkg)

default:
	go build cmd/gidari.go

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	chmod +rwx scripts/*.sh

	scripts/update-etc-hosts.sh
	scripts/build-storage.sh

	sleep 15 # need to sleep to allow mongodb topologies to come up

	scripts/build-migrations.sh

# proto is a phony target that will generate the protobuf files.
.PHONY: proto
proto:
	protoc --proto_path=pkg/proto --go_out=pkg/proto pkg/proto/db.proto

# test runs all of the application tests locally.
.PHONY: test
test:
	go test ./... -v

# ctests runs all of the application tests within a container, using the networking in the docker-compose file. this is
# useful for testing the applicaiton in CI/CD environments.
#
# This target is not intdeded to be run locally, for local use run `make test`
.PHONY: ctest
ctests:
	docker-compose -f "third_party/docker/storage.docker-compose.yaml" up
	docker-compose -f "third_party/docker/storage.docker-compose.yaml" run ctests

# ctest-local runs all of the application tests within a container, using the networking in the docker-compose file.
# This is the same as `make ctest` but it will use the local docker-compose file.
.PHONY: ctest-local
	chmod +rwx scripts/*.sh
	./scripts/run-ctests.sh
