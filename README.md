<div align="center">
    <img src=".github/banner.png" alt="Pocket Network logo" width="600"/>
    <h1>Portal API Database Driver</h1>
    <big>Database driver and struct definitions for use with the Portal API</big>
    <div>
    <br/>
        <a href="https://github.com/pokt-foundation/portal-db/v2/pulse"><img src="https://img.shields.io/github/last-commit/pokt-foundation/portal-db.svg"/></a>
        <a href="https://github.com/pokt-foundation/portal-db/v2/pulls"><img src="https://img.shields.io/github/issues-pr/pokt-foundation/portal-db.svg"/></a>
        <a href="https://github.com/pokt-foundation/portal-db/v2/issues"><img src="https://img.shields.io/github/issues-closed/pokt-foundation/portal-db.svg"/></a>
    </div>
</div>
<br/>

# Packages

## 1. Driver

Contains the following interfaces:

- **Driver**: contains all Read & Write methods.
- **Reader**: contains only Read methods and the Notification channel.
- **Writer**: contains only Write methods.

## 2. Postgres Driver

Contains all functionality to interact with Postgres.

- Provides a struct that satisfies the Driver interface.
- Typesafe Go code is generated from SQL schema by SQLC.
- Current Postgres version is `14.3`

## 3. Types

Contains all database structs and their associated methods which are used across the Portal API backend Go repos.

# Development

## Code Generation
⚠️ IMPORTANT - Before any commit the default `make` target **MUST** be run. ⚠️

In order for this target to properly generate the Go code and mocks the following dependencies **MUST** be installed:

- [SQLC](https://docs.sqlc.dev/en/stable/overview/install.html)
- [Mockery](https://vektra.github.io/mockery/latest/installation/)

## Modifying Functionality

- Start by updating the `postgresdrver/sqlc/schema.yml` and/or `postgresdrver/sqlc/query.yml` files
- Generate the SQLC Go code from these SQL files using `make`
- Add new functionality to the postgres driver methods and the Driver interface
- Update the driver integration tests in the `Portal DB` repository to cover all code changes
- Once merged to `main` a new version of the DB driver will be published, following semantic versioning standards based on the commit message(s)

## Pre-Commit Installation

Before starting development work on this repo, `pre-commit` must be installed.

In order to do so, run the command **`make init-pre-commit`** from the repository root.

Once this is done, the following checks will be performed on every commit to the repo and must pass before the commit is allowed:

### 1. Basic checks

- **check-yaml** - Checks YAML files for errors
- **check-merge-conflict** - Ensures there are no merge conflict markers
- **end-of-file-fixer** - Adds a newline to end of files
- **trailing-whitespace** - Trims trailing whitespace
- **no-commit-to-branch** - Ensures commits are not made directly to `main`

### 2. Go-specific checks

- **go-fmt** - Runs `gofmt`
- **go-imports** - Runs `goimports`
- **golangci-lint** - run `golangci-lint run ./...`
- **go-critic** - run `gocritic check ./...`
- **go-build** - run `go build`
- **go-mod-tidy** - run `go mod tidy -v`

### 3. Detect Secrets

Will detect any potential secrets or sensitive information before allowing a commit.

- Test variables that may resemble secrets (random hex strings, etc.) should be prefixed with `test_`
- The inline comment `pragma: allowlist secret` may be added to a line to force acceptance of a false positive

## Packages in Use

- [SQLC](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html) - Generates idiomatic type-safe Go code from SQL schema & queries.
- [Mockery](https://github.com/vektra/mockery) - Generates a mock from the Driver interface for testing purposes.

## Generating code

**Before committing any code to the repo, run the default Make target (`make`)**

This will used SQLC to generate Go code. It is also a useful way to check the `schema.sql` and `query.sql` files for SQL errors.

In addition, it will generate a mock of the `Driver` interface for testing purposes. This mock will automatically reflect changes made to the SQL schema files.
