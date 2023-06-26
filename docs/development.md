# Labrador Developer Documentation

## Local Development

### Requirements

- The `make` command.
- A configured Golang development environment.
- A configured Docker installation (optional).

### Building Only the Binary

```sh
make build

./dist/labrador_linux_amd64_v1/labrador version
```

### Building the Binary and Docker Image

```sh
make snapshot

./dist/labrador_linux_amd64_v1/labrador version

docker run -it --rm docker run -it --rm \
  ghcr.io/divergentcodes/labrador:dev-snapshot \
  version
```

### Build and Install the Local Binary

```sh
make install

labrador version
```

## Releasing

### Release Workflow Configuration

Most of what the Github Actions release workflow
(`.github/workflows/release.yaml`) does happens on the Github platform, and
the workflow is configured with all necessary Github permissions.

The Github Docker container registry is hosted through Github Packages, which
is organization wide in scope.

Note: A permission error will fail the release workflow if there are active
tag protection rules and no prior releases. Cut a single initial release,
and then create the tag protection rule `v*`.

### Triggering Releases

Releases are triggered by creating and pushing a Git tag in SemVer format.

```sh
git tag "v0.0.1-alpha01"; git push --tags
```

Tag rules:
- The tag should follow SemVer format, with optional trailing segment.
- The major/minor/patch versions should not have more than one leading `0`.
    The `pkg.go.dev` site won't accept docs that way.


Then the Github Actions workflow `.github/workflows/release.yaml` will:
- Login to the Docker container registry.
- Build everything.
- Push CLI binaries for multiple platforms to the repository Releases.
- Push a Docker image to the Docker container registry.
- Ask `pkg.go.dev` to pull new docs from the repository, if it is public.

### Release Artifacts

The Golang binary is compiled into several platform architectures, defined
in the `.goreleaser` file. Each file is bundled into a Github release,
along with checksums and a changelog.

A Docker image for the `x86_64` architecture is built and pushed to the Docker
container registry.


### Changelogs

GoReleaser will automatically generate a changelog and supply it along with the
release. Each new commit to the `main` branch will have an individual line
item in the changelog.

Changelog items are grouped by the commit message prefix:
- `feat:` prefixed commits are grouped under "Features" in the changelog.
- `bug:` prefixed commits are grouped under "Bug fixes" in the changelog.
- All other commits are grouped under "Others" in the changelog.

The commits with the following prefixes are omitted from the changelog:
- `ignore:`
- `docs:`
- `test:`

Commits with the following words anywhere are omitted from the changelog:
- `typo`
- `refactor`


## Package Documentation

Package documentation for the project is stored by `pkg.go.dev` at
https://pkg.go.dev/github.com/divergentcodes/labrador.

New versions are automatically fetched at the end of the release workflow.
Authenticating to `pkg.go.dev` is not necessary, since the site just pulls
whatever tagged/versioned documentation is present on Github. The
repository must be public for `pkg.go.dev` to pull documentation.

The `pkg.go.dev` site uses SemVer tags, and does not like multiple
leading zeros.

More information about adding documentation for a package can be found
at: https://pkg.go.dev/about#adding-a-package
