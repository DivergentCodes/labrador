# Each command is run in a separate sub-shell.
# Prefix commands with "@" to suppress command echo.
# Example: https://sohlich.github.io/post/go_makefile/

all: test build

# Setup local Git hooks.
githooks:
	./scripts/githooks.sh

# Generate Go docs.
docs:
	./scripts/docs.sh

# Serve HTML documentation using pkgsite.
docserve:
	./scripts/docserve.sh

# Run tests and generate HTML coverage report at [.coverage.html].
test:
	./scripts/test.sh

# Lint the code.
lint:
	./scripts/lint.sh

# Run SAST.
sast:
	./scripts/sast.sh

# Check the modules.
modcheck:
	./scripts/modcheck.sh

# Build single target binary using goreleaser config.
# Does not create a Docker image.
build:
	./scripts/build.sh

# Local installation of existing built executable.
_build_install:
	./scripts/install.sh

# Local build and installation of executable.
install: build _build_install

# Build point-in-time snapshot release.
# Creates a binary and a dev-snapshot Docker image. No Docker push.
snapshot:
	./scripts/snapshot.sh

# Create a release for supported platforms.
# Creates platform binaries and a x86_64 Docker image, all released to Github.
release:
	./scripts/release.sh
