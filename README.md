# Labrador

[![Go Reference](https://pkg.go.dev/badge/github.com/divergentcodes/labrador.svg)](https://pkg.go.dev/github.com/divergentcodes/labrador)

Pull secrets from remote services into CI/CD pipelines.

Labrador is a CLI tool to fetch variables and secrets from remote
services. Values are recursively pulled from one or more services,
and output to the terminal or a file.

The primary use case is enabling secretless pipelines in
conjunction with automated OIDC authentication to cloud providers,
where secrets are dynamically fetched instead of statically stored.

Additional CI/CD pipeline support:
- [labrador-action](https://github.com/marketplace/actions/labrador-action):
    A Github action for pulling variables and secrets into workflows.


## Quickstart

```sh
curl -sL https://github.com/DivergentCodes/labrador/releases/latest/download/labrador_Linux_x86_64.tar.gz  | tar -zx
./labrador fetch --aws-param "/path/to/params/*" --outfile "local.env"
```

You can also copy the example configuration file, `.labrador.example.yaml`
over to `.labrador.yaml` and customize it.


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

Examples:
- `LAB_VERBOSE=1`
- `LAB_AWS_SSM_PARAM=/base/path/to/params/*`
- `LAB_OUT_FILE=file.env`


### AWS Environment Variables

For Labrador to access secrets stored in AWS, configure the region and one
of the authentication methods. When running in a CI pipeline, a step may
configure some of these variables automatically, like the
[aws-actions/configure-aws-credentials](https://github.com/marketplace/actions/configure-aws-credentials-for-github-actions)
Github Action.

The AWS region always needs to be defined.
- `AWS_REGION`

Authenticating to AWS can be done using a profile, an API key, or an assumed role.
If running in a CI/CD pipeline, there may be an automation that configures the
session with AWS, like Github Actions' `aws-actions/configure-aws-credentials`.

- `AWS_PROFILE`
- `AWS_SECRET_ACCESS_KEY` and `AWS_ACCESS_KEY_ID`
- `AWS_ROLE_ARN` and `AWS_WEB_IDENTITY_TOKEN_FILE` and `AWS_ROLE_SESSION_NAME`.
