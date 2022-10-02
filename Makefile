# PKGS returns all Go packages in the Gidari code base.
PKGS = $(or $(PKG),$(shell env GO111MODULE=on $(GC) list ./...))

# TESTPKGS returns all Go packages int the Gidari code base that contain "*_test.go" files.
TESTPKGS = $(shell env GO111MODULE=on $(GC) list -f \
            '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
            $(PKGS))

# GC is the go compiler.
GC = go

export GO111MODULE=on

default:
	chmod +rwx scripts/*.sh
	$(GC) build cmd/gidari.go

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	scripts/build-storage.sh

# proto is a phony target that will generate the protobuf files.
.PHONY: proto
proto:
	protoc --proto_path=proto --go_out=proto proto/db.proto

# test runs all of the application tests locally.
.PHONY: tests
tests:
	$(GC) clean -testcache
	@$(foreach dir,$(TESTPKGS), $(GC) test $(dir) -v;)

# ci are the integration tests in CI/CD.
.PHONY: ci
ci:
	$(GC) clean -testcache
	./scripts/run-ci-tests.sh

# lint runs the linter.
.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml

# fmt runs the formatter.
.PHONY: fmt
fmt:
	gofumpt -l -w .
	golangci-lint run --fix

# add-license adds the license to all the top of all the .go files.
.PHONY: add-license
add-license:
	./scripts/add-license.sh
