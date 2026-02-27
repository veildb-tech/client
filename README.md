# VeilDB Client

[Documentation](https://veildb.gitbook.io/) • [Main Project](https://github.com/veildb-tech) • [Website](https://veildb.com)

Lightweight CLI application for securely downloading anonymized database dumps.

The Client is installed on developers' machines and provides simple access to prepared backups.

**This repository contains sources for the Client CLI. You can download a built binary from the Service Dashboard or directly from the documentation page https://veildb.gitbook.io/veildb-docs/user-guide/getting-started/installation#client-side**

---

## Part of VeilDB

This repository is **part** of the VeilDB platform.

- Main project overview: [https://github.com/veildb-tech](https://github.com/veildb-tech)
- Documentation: [https://veildb.gitbook.io/](https://veildb.gitbook.io/)

---

## Responsibilities

- Authenticate with VeilDB Service
- List available database dumps
- Download the latest anonymized versions

---

## Why Client?

Instead of manually sharing dumps:
- No SSH access required
- No unsafe file transfers
- No confusion about versions

Developers always download the correct, anonymized version.

---

## Typical Flow

1. Login
2. Select database
3. Download the latest processed dump

VeilDB is a CLI utility designed for convenient work on downloading the necessary database dumps


## Commands

```bash
Usage:
  veildb [command]

Available Commands:
  config      Configure application settings
  download    Downloading a dump of the database.
  help        Help about any command
  list        List commands
  login       Creating/updating a token and creating/updating a public key in the configuration file required for downloading database dumps.
  save-key    Creating/updating a PEM public key.

Flags:
  -h, --help      help for veildb
  -V, --version   Display this application version

Use "veildb [command] --help" for more information about a command.
```


## Technical Documentation:

### Flags
Flag ```-ldflags="-s -w"``` - is reduce the resulting binary size.
Flag ```-tags dev``` - served for compiling a binary file with settings for the dev server.

### Builds

#### Creating a build for Linux
To create a build for Linux, you need to go to the root of the directory where the code is contained and run the command: ```go build -ldflags="-s -w" -o veildb```
And then place the executable file in the bin directory (if it doesn’t exist, then create it)

#### Creating a build for Mac
To create a build for Mac, you need to go to the root of the directory where the code is contained and run the command: ```GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/veildb```

#### ARM
```GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/veildb```
And then place the executable file in the bin directory (if it doesn’t exist, then create it)

#### Creating a build for Windows
To create a build for Windows, you need to go to the root of the directory where the code is contained and run the command: ```GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/veildb.exe```



