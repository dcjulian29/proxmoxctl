# proxmoxctl

proxmoxctl is a full-featured command-line interface for managing Proxmox Virtual Environment (PVE) infrastructure.

[![Version](https://img.shields.io/github/v/release/dcjulian29/proxmoxctl)](https://github.com/dcjulian29/proxmoxctl/releases)
[![GitHub Issues](https://img.shields.io/github/issues-raw/dcjulian29/proxmoxctl.svg)](https://github.com/dcjulian29/proxmoxctl/issues)
[![Build](https://github.com/dcjulian29/proxmoxctl/actions/workflows/build.yml/badge.svg)](https://github.com/dcjulian29/proxmoxctl/actions/workflows/build.yml)

---

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Commands](#commands)
  - [config](#config)
- [Project Structure](#project-structure)
- [API Token Setup](#api-token-setup-in-proxmox)
- [Notes](#notes)

## Installation

#### Build from source code
```bash
git clone https://github.com/dcjulian29/proxmoxctl
cd proxmoxctl
go mod tidy
go build -o proxmoxctl .
```

## Configuration

Run the interactive setup once:

```bash
proxmoxctl config set
```

You'll be prompted for:
- **Server name** — a friendly label (e.g. `homelab`)
- **Server URL** — e.g. `https://192.168.1.10:8006`
- **API Token** — format: `USER@REALM!TOKENID=SECRET`

Config is saved to `~/.config/proxmoxctl/config.yaml`.

To view current config (token masked):
```bash
proxmoxctl config show
```

> **Tip:** You can also set values via environment variables:
>
> `PROXMOX_SERVER_URL`, `PROXMOX_API_TOKEN`

## Commands

### config

Manage connection settings stored in `~/.config/proxmoxctl/config.yaml`.

```bash
# Interactive setup
proxmoxctl config set

# Show current config (token is masked)
proxmoxctl config show
```


## Project Structure

```
proxmoxctl/
├── main.go
├── go.mod
├── cmd/
│   ├── root.go              # Root Cobra command, Viper init, --output flag
│   ├── config/
    |   ├── config.go        # config parent command
    |   ├── set.go           # Set the configuration for tool
│   │   └── show.go          # Show the configuration for the tool
└── internal/
    ├── output/
    │   └── output.go        # Shared table / JSON formatter (--output flag)
    └── settings/
        └── settings.go      # Viper config keys + save/require helpers
```

---

## API Token Setup in Proxmox

1. Go to **Datacenter → Permissions → API Tokens**
2. Add a token for your user (e.g. `user@pam!proxmoxctl`)
3. Assign permissions as needed — recommended roles:
   - `PVEAdmin` on `/` for full access (**USE CAREFULLY!!**)
   - Or scope roles per path for least-privilege setups (**Preferred Method**)
4. Copy the token secret — it is only shown once

Token format for this tool: `user@pam!proxmoxctl=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

---

## Notes

- **Environment variables** override config file values. Prefix any config key with `PROXMOX_` (e.g. `PROXMOX_API_TOKEN`).
- **JSON output** (`-o json`) is available on every read command and is suitable for piping into `jq` or other tools.
