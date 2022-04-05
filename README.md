# Gocrypt

gocrypt is a little command line tool to encrypt/decrypt file in AES128/256/512.

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
```

# Install

Download the binary in the release section. There is no docker image, but you can install a binary in your Docker image if needed.

# Tests

Tests are done with [venom](https://github.com/ovh/venom).

```
cd tests
venom run
```