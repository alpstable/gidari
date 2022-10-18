# Release Process

Releases will use [semantic versioning](https://semver.org/).

Releasing is a process that will require adaptation to the use case, the rules here are not set in stone.

## Patch Release

Assume the patch release is `X.Y.Z`

1. Cherry-pick the issues with label `x.y.z` in the patch from `main` into `release/X.Y`.
2. Create a PR to update
	1. the `release/X.Y` branch's `version/version.go` to `vX.Y.Z`
	2. the install instructions in README.md on `releaze/X.Y` to `@X.Y.Z`
3. Create a PR to update the install instructions in README.md on the `main` branch to `@vX.Y.Z`
4. Draft a new release [here](https://github.com/alpstable/gidari/releases/new)
	- Create a new tag `vX.Y.Z`
	- Target should be `release/X.Y`
	- Describe the release, include the the name of the label and link to the issues with that label filtered.
5. Publish the release

All issues with label `x.y.z` will be included in the patch release.

## Minor Release

Assume the minor release is `X.Y.0`

1. Create a PR updating the install instructions in README.md of the `main` branch to `@vX.Y.0`
2. Create a new branch from `main` called `release/X.Y`
3. Create a PR updating the `main` branch's `version/version.go` to `vX.(Y+1).0-prerelease`
4. Draft a new release [here](https://github.com/alpstable/gidari/releases/new)
	- Create a new tag `vX.Y.0`
	- Target should be `release/X.Y`
	- Describe the release, including the name of the label and link to the issues with that label filtered
5. Publish the release

All issues with label `x.y` wil be included in the minor release.
