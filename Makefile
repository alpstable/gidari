PKGS=$(shell scripts/list_pkgs.sh ./pkg)

default:
	chmod +rwx scripts/*.sh

	scripts/build_meta.sh
	scripts/build-storage.sh

	sleep 15 # need to sleep to allow mongodb topologies to come up
	scripts/build-migrations.sh

.PHONY: build-meta
build-meta:
	scripts/build_meta.sh

.PHONY: generate
generate:
	docker-compose -f "meta.docker-compose.yaml" run generate

.PHONY: proto
proto:
	protoc --proto_path=data/proto --go_out=data/proto data/proto/db.proto

.PHONY: test
test:
	go test ./... -v
