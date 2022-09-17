PKGS=$(shell scripts/list_pkgs.sh ./pkg)

default:
	go build cmd/gidari.go

.PHONY: containers
containers:
	chmod +rwx scripts/*.sh

	#scripts/update-etc-hosts.sh
	scripts/build-storage.sh

	sleep 15 # need to sleep to allow mongodb topologies to come up
	scripts/build-migrations.sh

.PHONY: proto
proto:
	protoc --proto_path=pkg/proto --go_out=pkg/proto pkg/proto/db.proto

.PHONY: test
test:
	go test ./... -v
