# Gocrypt

gocrypt is a command line tool to encrypt/decrypt file in AES128/256/512.

Usage is quite simple :

```
$ gocrypt -h
  -d    decrypt option
  -e    encrypt option
  -h    Print help
  -i string
        File to encrypt/decrypt
  -k string
        Name of file containing the key
  -o string
        Output file
  -v    Get version
```

# Install

Download the binary in the release section. There is no docker image, but you can install a binary in your Docker image if needed. If you want to create a docker image from scratch, you will need to do a multi stage docker build in order to download the binary.

## With homebrew

```
brew tap sgaunet/tools
brew install gocrypt
```

# Tests

Tests are done with [venom](https://github.com/ovh/venom).

```
cd tests
venom run
```

# Development

This project is using :

* golang
* [task for development](https://taskfile.dev/#/)
* docker
* [docker buildx](https://github.com/docker/buildx)
* docker manifest
* [goreleaser](https://goreleaser.com/)
* [venom](https://github.com/ovh/venom) : Tests
* [pre-commit](https://pre-commit.com/)

There are hooks executed in the precommit stage. Once the project cloned on your disk, please install pre-commit:

```
brew install pre-commit
```

And install the hooks

```
pre-commit install
```
