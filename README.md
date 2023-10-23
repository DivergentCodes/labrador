# Labrador

[![Go Reference](https://pkg.go.dev/badge/github.com/divergentcodes/labrador.svg)](https://pkg.go.dev/github.com/divergentcodes/labrador)

Labrador is a CLI tool to fetch secrets and other configuration values
from one or more remote services.

Labrador was created to explore safer, consistent, cross-platform ways of
handling secrets during each phase of the SDLC. The idea is to fetch secrets
from a central service at runtime, in a standard way, instead of copying secrets
to each environment and persisting them all over the place.

Example use cases:
- Secretless pipelines, where zero static secrets are stored in CI providers/pipelines
  (e.g. Github Actions, CircleCI, Jenkins, etc). Many CI providers support OIDC to
  cloud providers, so that not even cloud credentials need to be stored in pipelines.
- Smooth local development setup for projects, where one command sets up the
  environment configuration, instead of instructions with copy/paste steps.
  Developers can be sure that they are getting the latest configuration values.
- Modular app configuration, where each environment can combine an individual set of
  configuration values with another globally shared set of values.
- Portable app configuration, where secrets are fetched during runtime in the same
  way, regardless of the platform specifics.


## Contents

- [Quickstart](#quickstart)
- [Features](#features)
  - [Supported Value Stores](#supported-value-stores)
  - [CI/CD Pipeline Packages](#cicd-pipeline-packages)
- [Installation](#installation)
- [Example Usage](#example-usage)
  - [Fetch All AWS SSM Parameter Store Values at Given Base Path (Wildcard)](#fetch-all-aws-ssm-parameter-store-values-at-given-base-path-wildcard)
  - [Fetch Two Sets of AWS SSM Parameter Store Values](#fetch-two-sets-of-aws-ssm-parameter-store-values)
  - [Fetch an AWS Secrets Manager Value with multiple Key/Value Pairs](#fetch-an-aws-secrets-manager-value-with-multiple-keyvalue-pairs)

  - [Fetch GCP Secret Manager Values With or Without Explicit Versions](#fetch-gcp-secret-manager-values-with-or-without-explicit-versions)
  - [Fetch from Multiple Services At Once](#fetch-from-multiple-services-at-once)
  - [Save Fetched Values to an `.env` File](#save-fetched-values-to-an-env-file)
  - [Set Fetched Values as Environment Variables in the Current Shell](#set-fetched-values-as-environment-variables-in-the-current-shell)
  - [Use a Portable Config File for Consistent Value Fetching](#use-a-portable-config-file-for-consistent-value-fetching)
  - [Use Different Config Files for Local Development and CI/CD](#use-different-config-files-for-local-development-and-cicd)
- [Reference](#reference)
  - [Labrador Environment Variables](#labrador-environment-variables)
  - [AWS Environment Variables](#aws-environment-variables)
  - [GCP Environment Variables](#gcp-environment-variables)
- [Why Go (Golang)?](#why-go-golang)
- [Similar Projects](#similar-projects)

## Quickstart

Before running, set environment variables for accessing the services where the
values are stored. Some quick examples:

```sh
# AWS with a local profile.
export AWS_PROFILE="myprofile"
export AWS_REGION="us-east-1"

# GCP with a file containing service account key.
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

Download and extract the CLI.

```sh
curl -sL https://github.com/DivergentCodes/labrador/releases/latest/download/labrador_Linux_x86_64.tar.gz  | tar -zx
```

Fetch key/value pairs from AWS SSM Parameter Store, AWS Secret Manager,
and GCP Secret Manager, then save all of them in an `.env` file.

```sh
./labrador fetch \
  --aws-param "/path/to/test/params/*" \
  --aws-secret "my/test/aws-secret1" \
  --gcp-secret "projects/my-test-project/secrets/gcp-secret1" \
  --outfile ".env"
```

You can also copy the `.labrador.example.yaml` example configuration file
over to `.labrador.yaml`, customize it, and run `labrador fetch` or
`labrador export` without any arguments needed.

```yaml
aws:
  ssm_param:
  - /path/to/test/params/*
  sm_secret:
  - my/test/aws-secret1
gcp:
  sm_secret:
  - projects/my-test-project/secrets/gcp-secret1
```

Example of exporting into the current shell with the above configuation.

```sh
source <(./labrador export)
```


## Features

- **Configuration files**: Using an optional `.labrador.yaml` file, you
  can load all required environment variables with a single command.
  The file is portable and consistent across CI providers, environments,
  and anywhere that Labrador can run. Use alternate config files for local
  development, pipelines, and deployed environments.
- **Output options**: Use fetched values to create an `.env` file, set
  environment variables directly in the current shell, pass to another command,
  or print to the console. Key/value transforms include converting variable names
  to upper or lower case, and wrapping values in double quotes.
- **Wildcard paths**: For supported value stores, use a single wildcard resource
  path to recursively load all child values into the workflow
  ([example](#fetch-multiple-values-from-ssm-parameter-store-using-wildcard-paths)).
- **Multi-system fetching**: Labrador will pull from multiple remote stores in a
  single run. This can alleviate infrastructure migrations, multi-team situations,
  and other "real world quirks." Pulled values are coerced into a canonical format.


### Supported Value Stores

- **AWS SSM Parameter Store**: pull individual parameters, or recursively pull a wildcard path with all child variables, as individual environment variables.
- **AWS Secrets Manager**: pull all key/value pairs in a single secret are loaded as individual environment variables.
- **GCP Secret Manager**: pull each secret as a key/value pair, with or without explicit versions, and load as individual environment variables.

### CI/CD pipeline Packages

- [labrador-action](https://github.com/marketplace/actions/labrador-action):
    A Github action for pulling variables and secrets into workflows.


## Installation

Install and run the compiled binary.

```sh
curl -sL https://github.com/DivergentCodes/labrador/releases/latest/download/labrador_Linux_x86_64.tar.gz  | tar -zx
labrador version
```

Install and run the Docker container.

```sh
docker pull ghcr.io/divergentcodes/labrador
docker run -it --rm ghcr.io/divergentcodes/labrador version
```

Install and run from source.

```sh
git clone https://github.com/DivergentCodes/labrador
cd labrador
make install
labrador version
```
## Example Usage

### Fetch All AWS SSM Parameter Store Values at Given Base Path (Wildcard)

Instead of declaring each parameter individually, just point to a base path
and recursively fetch all child values. Add, update, or delete parameters in
AWS, without needing any environment configuration changes.

```sh
labrador fetch --aws-param "/path/to/params/*"
```

### Fetch Two Sets of AWS SSM Parameter Store Values

You can split configuration up into multiple sets of values, each with different
IAM permissions.

For instance, one base path in AWS Parameter Store can contain shared configuration
used across all environments, while another base path has parameters for a
specific instance.

```sh
labrador fetch --aws-param "/global/shared/params/*" --aws-param "/instance/params/*"
```

### Fetch an AWS Secrets Manager Value With multiple Key/Value Pairs

A single secret in AWS Secrets Manager can store multiple key/value pairs.
Labrador will pull the secret, extract each key/value, and return them as
individual variables.

```sh
labrador fetch --aws-secret "path/to/secret"
```

### Fetch GCP Secret Manager Values With or Without Explicit Versions

Labrador will default to a GCP secret's `versions/latest` version if
one isn't explicitly specified. Each secret can hold a single key/value
pair.

```sh
labrador fetch \
  --gcp-secret "projects/myproject/secrets/foo" \
  --gcp-secret "projects/myproject/secrets/bar/versions/1"
```

### Fetch from Multiple Services At Once

If your configuration is spread across multiple services (e.g. undergoing
a migration), you can fetch values from all of them at the same time.

```sh
labrador fetch --aws-param "/path/to/params/*" --gcp-secret "projects/myproject/secrets/mysecret"
```

### Save Fetched Values to an `.env` File

An [`.env` file](https://www.dotenv.org/docs/security/env.html)
is useful for ergonomically defining environment variables in development,
with tooling like
[Docker](https://docs.docker.com/compose/environment-variables/env-file/)
and [dotenv](https://www.npmjs.com/package/dotenv).

```sh
labrador fetch --aws-param "/path/to/params/*" --outfile ".env"
```

### Set Fetched Values as Environment Variables in the Current Shell

This example assumes a `.labrador.yaml` configuration file exists in the current
working directory.

```sh
source <(labrador export)
```

### Use a Portable Config File for Consistent Value Fetching

Instead of each developer manually setting development variables as a setup step
for a project, just have them run Labrador to pull dev values from a centralized
store. Values can be updated/rotated in the remote store and be instantly available
to all developers, removing tedious and error prone configuration changes.

These configuration files can safely be committed to the repository, providing
easier project setup for N developers, easier value updates, and no secrets stored
on developer machines.

```yaml
# .labrador.yaml
aws:
  region: us-east-1
  sm_secret:
  - app/dev/values
```

```sh
labrador fetch
```


### Use Different Config Files for Local Development and CI/CD

You can use separate configuration files per environment, to keep secrets and
configuration settings isolated, consistent, canonical, and reusable. Continue
to control access

Create one configuration file that will pull all needed values for local development.

```yaml
# .labrador.dev.yaml
aws:
  region: us-east-1
  ssm_param:
  - /app/shared/values/*
  - /app/dev/values/*
```

Create a second configuration file that will pull all needed values for builds in the pipeline.

```yaml
# .labrador.ci.yaml
aws:
  region: us-east-1
  ssm_param:
  - /app/shared/values/*
  - /app/ci/values/*
```

Run the development configuration to pull values needed for local development.

```sh
source <(labrador export --config .labrador.dev.yaml)
```

In the pipeline, run the CI configuration to pull values needed for builds or
deployments.

```sh
labrador fetch --config .labrador.ci.yaml
```


## Reference

### Labrador Environment Variables

All CLI arguments can also be configured as environment variables,
prefixed with `LAB_`.

Examples:
- `LAB_AWS_SM_SECRET=name/of/secret`
- `LAB_AWS_SSM_PARAM=/base/path/to/params/*`
- `LAB_OUT_FILE=file.env`
- `LAB_VERBOSE=1`


### AWS Environment Variables

For Labrador to access secrets stored in AWS, configure the region and one
of the authentication methods. When running in a CI pipeline, a step may
configure some of these variables automatically, like the
[aws-actions/configure-aws-credentials](https://github.com/marketplace/actions/configure-aws-credentials-for-github-actions)
Github Action.

The AWS region always needs to be defined. It can also be set as a CLI
argument or in the configuration file.
- `AWS_REGION`

Authenticating to AWS can be done using a profile, an API key, or an assumed role.
If running in a CI/CD pipeline, there may be an automation that configures the
session with AWS, like Github Actions' `aws-actions/configure-aws-credentials`.

- `AWS_PROFILE`
- `AWS_SECRET_ACCESS_KEY` and `AWS_ACCESS_KEY_ID`
- `AWS_ROLE_ARN` and `AWS_WEB_IDENTITY_TOKEN_FILE` and `AWS_ROLE_SESSION_NAME`.

### GCP Environment Variables

For Labrador to access secrets stored in GCP, configure the
[Application Default
Credentials (ADC)](https://cloud.google.com/docs/authentication/provide-credentials-adc). If the credentials file is not stored in the default
location, an environment variable can point to the file.

- `GOOGLE_APPLICATION_CREDENTIALS`


## Why Go (Golang)?

- Golang has official SDKs for most popular cloud platforms and value stores.
- A single binary release results in simpler distribution, smaller artifacts, and more versitile usage.
- The [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) packages are excellent for creating CLI tools.


## Similar Projects

- [segmentio/chamber](https://github.com/segmentio/chamber):
  CLI tool to fetch from AWS SSM Parameter Store.
- [hashicorp/vault-action](https://github.com/hashicorp/vault-action): Github
  Action to fetch from Hashicorp Vault.
- [aws-actions/aws-secretsmanager-get-secrets](https://github.com/aws-actions/aws-secretsmanager-get-secrets):
  Github Action to fetch from AWS Secrets Manager.
- [google-github-actions/get-secretmanager-secrets](https://github.com/google-github-actions/get-secretmanager-secrets):
  Github Action to fetch from GCP Secret Manager.
- [dotenv-org/dotenv-vault](https://github.com/dotenv-org/dotenv-vault):
  File based value storage that uses Dotenv's servers for synchronizing.
