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
  - [backup](#backup)
  - [clone](#clone)
  - [config](#config)
  - [group](#group)
  - [lxc](#lxc--containers)
  - [snapshot](#snapshot)
  - [status](#status)
  - [storage](#storage)
  - [user](#user)
  - [version](#version)
  - [vm](#vm--kvm-virtual-machines)
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

### backup

Manage on-demand backups, backup files in storage, restore operations, and scheduled backup jobs.

#### On-demand backups

```bash
# Back up a single VM
proxmoxctl backup create 100 --storage local

# Back up multiple guests
proxmoxctl backup create 100,101,102 --storage backup-nfs

# Back up all guests on the node
proxmoxctl backup create all --storage backup-nfs

# Custom mode and compression
proxmoxctl backup create 100 --storage local --mode stop --compress lzo

# Auto-prune backups older than 7 days
proxmoxctl backup create 100 --storage local --remove-older 7
```

#### Listing and inspecting backup files

```bash
# List all backups on a storage
proxmoxctl backup list --storage local

# Filter by VM ID
proxmoxctl backup list --storage local --vmid 100

# Inspect a specific backup file
proxmoxctl backup show --storage local \
  --file vzdump-qemu-100-2024_01_15-03_00_01.vma.zst
```

#### Restore

```bash
# Restore a VM from backup
proxmoxctl backup restore --vmid 100 --storage local \
  --file vzdump-qemu-100-2024_01_15-03_00_01.vma.zst

# Restore an LXC container and start it immediately
proxmoxctl backup restore --vmid 300 --type lxc --storage local \
  --file vzdump-lxc-300-2024_01_15-03_00_01.tar.zst --start
```

#### Delete a backup file

```bash
proxmoxctl backup delete --storage local \
  --file vzdump-qemu-100-2024_01_15-03_00_01.vma.zst
```

#### Scheduled backup jobs

```bash
# List all scheduled jobs
proxmoxctl backup jobs list

# Show a specific job
proxmoxctl backup jobs show job-abc123

# Create a scheduled job (runs at 2am daily)
proxmoxctl backup jobs create \
  --vmids 100,101,200 \
  --storage backup-nfs \
  --schedule "0 2 * * *" \
  --compress zstd \
  --max-files 7 \
  --mailto admin@example.com

# Back up all guests on a schedule
proxmoxctl backup jobs create --all --storage backup-nfs --schedule "0 3 * * 0"

# Modify a job
proxmoxctl backup jobs modify job-abc123 --max-files 14 --enabled true

# Delete a job
proxmoxctl backup jobs delete job-abc123
```

**Flags:** `--storage` (required), `--mode` (snapshot|suspend|stop), `--compress` (zstd|lzo|gzip), `--vmid`, `--file`, `--type` (qemu|lxc), `--schedule`, `--max-files`, `--mailto`, `--start`, `--all`, `--force`, `--node`

### clone

Clone KVM VMs or LXC containers. Supports full clones, linked clones, and cloning from a named snapshot.

```bash
# Full clone of VM 100 → new VM 200
proxmoxctl clone vm 100 --newid 200 --name webserver-copy

# Linked clone (shares base disk; source must be a template)
proxmoxctl clone vm 100 --newid 201 --name webserver-linked --linked

# Clone from a specific snapshot
proxmoxctl clone vm 100 --newid 202 --name from-snap --snapname before-upgrade

# Clone to a specific storage
proxmoxctl clone vm 100 --newid 203 --name migrated --storage ceph-pool

# Clone an LXC container
proxmoxctl clone lxc 300 --newid 400 --hostname new-container
```

**VM flags:** `--newid` (required), `--name`, `--snapname`, `--linked`, `--full`, `--storage`, `--pool`, `--node`

**LXC flags:** `--newid` (required), `--hostname`, `--snapname`, `--storage`, `--pool`, `--node`

### config

Manage connection settings stored in `~/.config/proxmoxctl/config.yaml`.

```bash
# Interactive setup
proxmoxctl config set

# Show current config (token is masked)
proxmoxctl config show
```

### group

Manage Proxmox user groups.

```bash
# List all groups
proxmoxctl group list

# Show group details (including members)
proxmoxctl group show admins

# Create a group
proxmoxctl group create developers --comment "Dev team"

# Modify a group
proxmoxctl group modify developers --comment "Development team"

# Delete a group
proxmoxctl group delete developers
proxmoxctl group delete developers --force
```

**Flags:** `--comment`, `--force`

### lxc — Containers

```bash
# List containers
proxmoxctl lxc list

# Show container status
proxmoxctl lxc status 101

# Create a container
proxmoxctl lxc create --vmid 300 --hostname mycontainer \
  --template local:vztmpl/debian-12-standard_12.2-1_amd64.tar.zst \
  --memory 1024 --cores 2 --disk local-lvm:8

# Modify
proxmoxctl lxc modify 300 --memory 2048 --hostname newname

# Power control
proxmoxctl lxc start 300
proxmoxctl lxc stop 300

# Delete
proxmoxctl lxc delete 300
proxmoxctl lxc delete 300 --force
```

**Flags:** `--node`, `--vmid`, `--hostname`, `--template`, `--memory`, `--cores`, `--disk`, `--password`, `--force`

### snapshot

Manage snapshots for both KVM VMs and LXC containers. Use `--type lxc` to target containers (default: `qemu`).

```bash
# List snapshots
proxmoxctl snapshot list 100
proxmoxctl snapshot list 101 --type lxc

# Show config stored inside a snapshot
proxmoxctl snapshot show 100 --name before-upgrade

# Create a snapshot
proxmoxctl snapshot create 100 --name before-upgrade --desc "Pre v2 migration"

# Create with RAM state (running VMs only)
proxmoxctl snapshot create 100 --name with-ram --vmstate

# Roll back to a snapshot
proxmoxctl snapshot rollback 100 --name before-upgrade

# Delete a snapshot
proxmoxctl snapshot delete 100 --name before-upgrade --force
```

**Flags:** `--node`, `--name` (required), `--type` (qemu|lxc), `--desc`, `--vmstate`, `--force`

### status

Show health and resource usage for a node or the entire cluster.

```bash
# Cluster-wide overview: quorum, node list, CPU/mem/uptime per node
proxmoxctl status cluster

# Deep dive into a single node (CPU, memory, swap, root FS with usage bars)
proxmoxctl status node
proxmoxctl status node --node pve2

# All cluster resources grouped by type (VMs, LXC, storage, nodes)
proxmoxctl status resources

# Filter resources by type
proxmoxctl status resources --type vm
proxmoxctl status resources --type storage

# Recent task log
proxmoxctl status tasks
proxmoxctl status tasks --limit 50
proxmoxctl status tasks --errors        # show only failed tasks
proxmoxctl status tasks --node pve2
```

**Flags:** `--node`, `--type` (vm|lxc|storage|node), `--limit`, `--errors`

### storage

List storage pools and inspect their contents.

```bash
# List all storage pools (cluster config view)
proxmoxctl storage list

# List with live usage stats for a specific node
proxmoxctl storage list --node pve1

# Only show active/enabled storage
proxmoxctl storage list --active

# Filter to storage that supports a content type
proxmoxctl storage list --content backup

# Show full details of one storage pool
proxmoxctl storage show local
proxmoxctl storage show local --node pve1   # includes usage stats

# List contents of a storage pool
proxmoxctl storage content local
proxmoxctl storage content local --type iso
proxmoxctl storage content local --type backup
proxmoxctl storage content local --type vztmpl
proxmoxctl storage content local --vmid 100
```

**Flags:** `--node`, `--active`, `--content`, `--type`, `--vmid`

### user

Manage Proxmox users. User IDs are always in `USER@REALM` format (e.g. `alice@pam`, `bob@pve`).

```bash
# List all users
proxmoxctl user list

# List only enabled users
proxmoxctl user list --enabled

# Show a user
proxmoxctl user show alice@pam

# Create a user
proxmoxctl user create alice@pam \
  --password secret \
  --firstname Alice \
  --lastname Smith \
  --email alice@example.com \
  --groups devs,ops

# Modify a user
proxmoxctl user modify alice@pam \
  --email newemail@example.com \
  --groups devs,ops,admins \
  --enabled true

# Disable an account
proxmoxctl user modify alice@pam --enabled false

# Change password (prompts securely if --password not passed)
proxmoxctl user passwd alice@pam
proxmoxctl user passwd alice@pam --password newpassword

# Show which groups a user belongs to
proxmoxctl user groups alice@pam

# Delete a user
proxmoxctl user delete alice@pam
proxmoxctl user delete alice@pam --force
```

**Flags:** `--password`, `--firstname`, `--lastname`, `--email`, `--comment`, `--groups`, `--enabled`, `--expire`, `--force`

### version

Show the tool version, and check for updates.

```bash
proxmoxctl version
```

### vm — KVM Virtual Machines

```bash
# List all VMs
proxmoxctl vm list
proxmoxctl vm list -o json

# Show detailed status
proxmoxctl vm status 100

# Create a VM
proxmoxctl vm create --vmid 200 --name myvm --memory 4096 --cores 2 \
  --disk local-lvm:32 --iso local:iso/debian-12.iso

# Modify a VM
proxmoxctl vm modify 200 --name newname --memory 8192 --cores 4

# Power control
proxmoxctl vm start 200
proxmoxctl vm stop 200

# Delete (prompts for confirmation)
proxmoxctl vm delete 200
proxmoxctl vm delete 200 --force
```

**Flags:** `--node`, `--vmid`, `--name`, `--memory`, `--cores`, `--disk`, `--iso`, `--force`

## API Token Setup in Proxmox

1. Go to **Datacenter → Permissions → API Tokens**
2. Add a token for your user (e.g. `user@pam!proxmoxctl`)
3. Assign permissions as needed — recommended roles:
   - `PVEAdmin` on `/` for full access (**USE CAREFULLY!!**)
   - Or scope roles per path for least-privilege setups (**Preferred Method**)
4. Copy the token secret — it is only shown once

Token format for this tool: `user@pam!proxmoxctl=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

## Notes

- **Environment variables** override config file values. Prefix any config key with `PROXMOX_` (e.g. `PROXMOX_API_TOKEN`).
- **JSON output** (`-o json`) is available on every read command and is suitable for piping into `jq` or other tools.
- **Self-Signed Certificates** (`--insecure`) disables TLS certificate verification. Lack of certification verification can lead to man-in-the-middle attacks and is not recommended to do this outside of a lab environment. This can be set in the configuration file with: `tls_insecure: true` or be provided with each command via the flag.
