# awsvault

Provides convenience functions for working with AWS profiles via aws-vault.

## Features

- Opens an AWS profile in a sandboxed [Firefox Multi-Account Container](https://support.mozilla.org/en-US/kb/containers) or Google Chrome instance
- Open a subshell with AWS credentials for a profile
- List all AWS profiles
- Fuzzy-select a profile

## Install

```bash
go install github.com/joakimen/awsvault@latest
```

## Usage

```bash
$ awsvault
Helper functions for working with aws-vault

Usage:
  awsvault [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  exec        Start a subshell with temporary credentials
  help        Help about any command
  list        List profiles in ~/.aws/config
  login       Open a browser window and login using temporary credentials
  select      Select a single AWS profile interactively

Flags:
  -h, --help   help for awsvault

Use "awsvault [command] --help" for more information about a command.
```

### Login

```bash
$ awsvault login
# Select a profile interactively
# A browser window opens with temporary credentials
```

### Exec

```bash
$ awsvault exec
# Select a profile interactively
# A subshell opens with temporary credentials
```

### List

```bash
$ awsvault list
company1-dev
company1-prod
company9-dev
```

### Select

```bash
$ awsvault select
# Select a profile interactively
company1-dev
```
