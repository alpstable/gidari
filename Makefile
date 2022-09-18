PKGS=$(shell scripts/list_pkgs.sh ./pkg)

default:
	go build cmd/gidari.go

# storage will build the storage containers for integration testing (postgres, redis, etc)
.PHONY: storage
storage:
	chmod +rwx scripts/*.sh
	chmod +rwx third_party/docker/rs-init.sh

	scripts/build-storage.sh

# migrations will run the liquibase migrations on the storage containers.
.PHONY: migrations
migrations:
	chmod +rwx scripts/*.sh

	scripts/build-migrations.sh

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	make storage
	make migrations

# proto is a phony target that will generate the protobuf files.
.PHONY: proto
proto:
	protoc --proto_path=pkg/proto --go_out=pkg/proto pkg/proto/db.proto

# test runs all of the application tests locally.
.PHONY: test
test:
	godotenv -f /etc/alpine-hodler/auth.env go test ./... -v

# ctests runs all of the application tests within a container, using the networking in the docker-compose file. this is
# useful for testing the applicaiton in CI/CD environments.
#
# This target is not intdeded to be run locally, for local use run `make test`
.PHONY: ctests
ctests:
	chmod +rwx scripts/*.sh
	./scripts/run-ctests.sh
