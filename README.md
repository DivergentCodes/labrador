# Labrador

[![Go Reference](https://pkg.go.dev/badge/github.com/divergentcodes/labrador.svg)](https://pkg.go.dev/github.com/divergentcodes/labrador)

Pull secrets from remote services into CI/CD pipelines.

Labrador is a CLI tool to fetch variables and secrets from remote
services. The primary use case is enabling secretless pipelines in
conjunction with automated OIDC authentication to cloud providers,
where secrets are dynamically fetched instead of statically stored.

Values are recursively pulled from one or more services, and output
to the terminal or a file.

Labrador is focused on pulling values, not on managing or writing values.

## Installation

### Binary from Github Releases

Download and run the binary CLI from Github releases.

```sh
curl -sL https://github.com/DivergentCodes/labrador/releases/latest/download/labrador_Linux_x86_64.tar.gz  | tar -zx
labrador version
```

### Docker from Github Container Registry

Pull and run the Docker image from the Github container registry.

```sh
docker pull ghcr.io/divergentcodes/labrador
docker run -it --rm ghcr.io/divergentcodes/labrador version
```

### Binary from Source

Download and run the binary CLI from Github releases.

```sh
git clone https://github.com/DivergentCodes/labrador
cd labrador
make install
labrador version
```

## Environment Variables

All CLI arguments can also be configured as environment variables,
prefixed with `LAB_`.

### AWS

For Labrador to access secrets stored in AWS, configure the region and one
of the authentication methods. When running in a CI pipeline, a step may
configure some of these variables automatically, like the `aws-actions/configure-aws-credentials` Github Action.

- `AWS_REGION`
- `AWS_PROFILE`
- `AWS_SECRET_ACCESS_KEY` and `AWS_ACCESS_KEY_ID`
