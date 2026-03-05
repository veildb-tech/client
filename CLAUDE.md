# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## About

VeilDB is a CLI utility (built with Cobra) for downloading database dumps from the VeilDB service. It handles authentication, RSA public-key encryption of download requests, and file retrieval.

## Build Commands

```bash
# Standard Linux build
go build -ldflags="-s -w" -o veildb

# Dev build (points to dev server: https://db-manager.bridge2.digital)
go build -tags dev -ldflags="-s -w" -o veildb

# Cross-compile
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/veildb
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/veildb
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/veildb.exe
```

The `-tags dev` build flag selects `services/config_dev.go` (dev server URL) instead of `services/config.go` (production URL `https://app.veildb.com`). Both files are otherwise identical.

## Architecture

The codebase follows a two-layer pattern: `cmd/` registers Cobra commands and delegates to `processes/`, which contains the business logic. Shared infrastructure lives in `services/`.

```
cmd/           - Cobra command definitions (login, download, save-key, install, config, list)
processes/     - Command implementations
  login/       - Auth flow: prompts for service URL (if unset), credentials, fetches token, saves PEM key
  download/    - Fetches DB/dump UIDs interactively, encrypts request with RSA public key, downloads file
  savekey/     - Creates/updates PEM public key files
  install/     - CLI installation helper
  config/      - Updates dump save path and/or service URL in config
services/
  envfile/     - Config file CRUD (~/.veildb/.env.json); defines Config, Workspace, Server structs
  encrypter/   - RSA PKCS1v15 encryption of download request data using stored PEM public key
  keypubfile/  - Reads PEM public key files from ~/.veildb/
  workspace/   - Workspace and server selection logic
  request/     - HTTP GET helper with Bearer token auth
  response/    - Error response parsing
  token/       - Token utilities
  predefined/  - Colored terminal output helpers (BuildOk, BuildError, BuildWarning)
  config.go    - Production constants + GetServiceUrl(), CurrentAppDir(), URL builders
  config_dev.go - Dev constants (build tag: dev)
util/helper.go - Misc utilities
```

## Config File

The app stores state in `~/.veildb/.env.json` with this structure:
```json
{
  "service_url": "https://app.veildb.com",
  "dump_path": "/optional/default/dump/path",
  "token": "jwt-token",
  "current_workspace": "workspace-name",
  "data": {
    "workspace-name": {
      "servers": {
        "server-name": {
          "key_file": "workspace_server.pem",
          "server_id": "uuid"
        }
      }
    }
  }
}
```

PEM public key files are stored alongside the config in `~/.veildb/`.

## Download Flow

1. Read saved config (token, workspace, server-to-PEM mapping)
2. Interactively select DB and dump (or accept `--db-uid`/`--dump-uid` flags)
3. Encrypt `{dbuuid, dumpuuid, dumpname, dumppath}` with RSA public key via `encrypter.EncryptData`
4. GET download link from API, then POST encrypted payload to that link to receive the `.sql.gz` file
