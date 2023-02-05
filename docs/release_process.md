# Release Process

Releases will use [semantic versioning](https://semver.org/).

## Patch Release

Patch releases are for backwards compatible bug fixes.

1. Create a draft release from the [release template](#release-template) but do not release it yet. The title should be "vX.Y.Z".
2. Put the release notes on the discord server for review: https://discord.com/invite/3jGYQz74s7
3. Cherry-pick the issues with label `x.y.z` in the patch from `main` into `release/x.y`.
4. Change the version on `release/X.Y` to `x.y.z` and commit the change with the message "updating release version x.y.z"
5. Create and push the tag using `git tag -a vx.y.z -m "vx.y.z"; git push --tags`
6. Tag the the draft release with the `vx.y.z` tag.
7. Publish the release.
8. Change the version on `release/x.y` to `x.y.z+1-prerelease` and commit the change with the message "updating release version x.y.z+1-prerelease".
9. Put the release notes on the discord server: https://discord.com/invite/3jGYQz74s7
10. Close the version project

All issues with label `x.y.z` will be included in the patch release.

## Minor Release

Minor releases are for adding functionality in a backwards compatible manner.

1. Create a draft release from the [release template](#release-template) but do not release it yet. The title should be "vX.Y.0".
2. Put the release notes on the discord server for review: https://discord.com/invite/3jGYQz74s7
3. Change the version on `main` to `x.y.0` and commit the change with the message "updating release version x.y.0"
4. Create a new branch from `main` called `release/x.y`
6. On `release/x.y`, create and push the tag using `git tag -a vx.y.0 -m "vx.y.0"; git push --tags`
7. Tag the the draft release with the `x.y.0` tag.
8. Change the version on `release/x.y` to `x.y.1-prerelease` and commit the change with the message "updating version x.y.1-prelease"
9. Publish the release.
10. Change the version on `main` to `x.y+1.0-prerelease` and commit the change with the message "updating release version x.y+1.0-prerelease"
11. Close the version project

All issues with label `x.y.0` wil be included in the minor release.

## Release Template

```
{Description of the release}

### Migration Steps
* [ACTION REQUIRED]
*

### Breaking Changes
*
*

### New Features
*
*

### Bug Fixes
*
*

### Performance Improvements
*
*

### Other Changes
*
*
```
