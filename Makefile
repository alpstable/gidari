PKGS=$(shell scripts/list_pkgs.sh ./pkg)

default:
	scripts/build_meta.sh

.PHONY: build-meta
build-meta:
	scripts/build_meta.sh

.PHONY: generate-meta
generate:
	docker-compose -f "meta.docker-compose.yaml" run generate

.PHONY: build-proto
proto:
	protoc --proto_path=data/proto --go_out=data/proto data/proto/db.proto
