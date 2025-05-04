[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/gocrypt)](https://goreportcard.com/report/github.com/sgaunet/gocrypt)
[![GitHub release](https://img.shields.io/github/release/sgaunet/gocrypt.svg)](https://github.com/sgaunet/gocrypt/releases/latest)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/gocrypt/total)
[![GoDoc](https://godoc.org/github.com/sgaunet/gocrypt?status.svg)](https://godoc.org/github.com/sgaunet/gocrypt)
[![License](https://img.shields.io/github/license/sgaunet/gocrypt.svg)](LICENSE)

# Gocrypt

gocrypt is a simple and efficient command line tool to encrypt and decrypt files using AES-128 or AES-256 encryption.

## Features
- Encrypt and decrypt files using AES-128 or AES-256 (GCM mode)
- Key can be provided via a file or the `GOCRYPT_KEY` environment variable
- Easy integration in scripts and automation
- Cross-platform: works on Linux, macOS, and Windows

## Usage

```sh
$ gocrypt help
Tool to encrypt/decrypt files using AES128 or AES256.

Usage:
  gocrypt [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  dec         decrypt file in AES 128/256
  enc         encrypt file in AES 128/256
  help        Help about any command
  version     print version of gocrypt

Flags:
  -h, --help   help for gocrypt

Use "gocrypt [command] --help" for more information about a command.
```

### Example

```sh
# Encrypt a file with a 32-byte key (AES-256)
gocrypt enc --i input.txt --o encrypted.bin --k keyfile.txt

# Decrypt a file
gocrypt dec --i encrypted.bin --o decrypted.txt --k keyfile.txt
```

### Key Format
- For AES-128, the key must be exactly 16 bytes (characters)
- For AES-256, the key must be exactly 32 bytes (characters)
- The key can be provided in a file or via the `GOCRYPT_KEY` environment variable

## Demo

![Demo](doc/demo.gif)

## Recommendation
Even though this tool is maintained, for most use-cases you should consider using [age](https://github.com/FiloSottile/age) which is a more modern and secure encryption tool.

# Install

Download the binary from the [releases section](https://github.com/sgaunet/gocrypt/releases/latest). There is no official Docker image, but you can add the binary to your own Docker image if needed.

## With Homebrew

```sh
brew tap sgaunet/homebrew-tools
brew install sgaunet/tools/gocrypt
```

## Use in Docker

```Dockerfile
FROM sgaunet/gocrypt:latest as gocrypt

FROM alpine:latest
COPY --from=gocrypt /gocrypt /usr/local/bin/gocrypt
...
```

# Tests

The project includes automated tests for both small and large files. See the `tests/` directory for details.

```
task tests
```

# License

MIT License
