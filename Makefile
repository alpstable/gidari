PKGS=$(shell scripts/list_pkgs.sh ./pkg)

default:
	go build cmd/sherpa.go

.PHONY: containers
containers:
	chmod +rwx scripts/*.sh

	scripts/build-storage.sh

	sleep 15 # need to sleep to allow mongodb topologies to come up
	scripts/build-migrations.sh

.PHONY: proto
proto:
	protoc --proto_path=proto --go_out=proto proto/db.proto

.PHONY: test
test:
	go test ./... -v
