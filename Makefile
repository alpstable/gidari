PKGS=$(shell scripts/list_pkgs.sh ./pkg)

default:
	scripts/build_meta.sh
	scripts/build_bazel.sh

build-bazel:
	scripts/build_bazel.sh

build-meta:
	scripts/build_meta.sh

generate:
	docker-compose -f "meta.docker-compose.yaml" run generate

list_pkgs:
	echo $(PKGS)

start-graphql:
	docker-compose -f "graphql.docker-compose.yaml" up -d --build
