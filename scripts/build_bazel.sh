#!/usr/bin/env bash

# the mod cache can cause some goofy warnings if its not cleaned before we build
go clean -modcache
go mod tidy

bazel clean
bazel run //:gazelle

# Whenever we add new go dependencies we need to update the repos through bazel
bazel run //:gazelle-update-repos

bazel build //... --spawn_strategy=standalone

# run tests
bazel test //... \
	--test_verbose_timeout_warnings \
	--sandbox_debug \
	--spawn_strategy=standalone
