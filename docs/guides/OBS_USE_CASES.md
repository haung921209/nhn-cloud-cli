# Object Storage (OBS) Use Cases

The `nhncloud object-storage` (alias `obs`) command allows you to manage objects and containers directly from the CLI. It supports efficient file transfers, including automatic handling of large files via Static Large Objects (SLO).

## 1. List Containers and Objects

Use the `ls` command to view your containers or lists objects within them.

### List All Containers
```bash
nhncloud obs ls
```

### List Objects
By default, `ls` shows a **directory view** (grouping objects by folder).

```bash
nhncloud obs ls obs://my-container/
# Output:
#                            PRE images/
#                            PRE data/
# file.txt      1024    2023-10-01T12:00:00
```

To list **all objects recursively** (flat view), use the `--recursive` (`-r`) flag:

```bash
nhncloud obs ls -r obs://my-container/
# Output:
# images/logo.png
# data/log.txt
# file.txt
```

---

## 2. File Operations (cp)

The `cp` command copies files between your local machine and Object Storage. Use the `--recursive` (`-r`) flag to copy entire directories.

### Upload Files & Directories
```bash
# Upload a single file
nhncloud obs cp ./image.png obs://my-container/images/image.png

# Upload a directory (Recursive)
nhncloud obs cp -r ./local-dir obs://my-container/remote-dir
```

# Upload with implicit filename
nhncloud obs cp ./document.pdf obs://my-container/docs/
```

> [!NOTE]
> **Large File Support (SLO)**
> Files larger than 5GB are **automatically** split and uploaded as Static Large Objects (SLO). The segments are stored in a shadow container named with the `_segments` suffix (e.g., `my-container_segments`).
> You can control the segment size (default 1GB) with the `--segment-size` flag.

```bash
# Upload a 10GB file (automatically split into 1GB segments)
nhncloud obs cp ./large-backup.tar.gz obs://backups/
```

### Download Files & Directories
```bash
# Download a single file
nhncloud obs cp obs://my-container/file.txt .

# Download a directory (Recursive)
nhncloud obs cp -r obs://my-container/remote-dir ./local-dir
```

### Copy Between Containers
```bash
# Copy a single file
nhncloud obs cp obs://src/file.txt obs://dst/file.txt

# Copy a directory (Recursive)
nhncloud obs cp -r obs://src/folder obs://dst/backup-folder
```

---

## 3. Configuration

Object Storage authentication usually follows your global `tenant-id`. However, if your Object Storage service is in a different tenant (common in some organization setups), you can configure a specific Tenant ID for OBS:

```bash
nhncloud configure
# ...
# Object Storage Tenant ID [current-id]: <new-obs-tenant-id>
```

Or manually verify in `~/.nhncloud/credentials`:
```ini
[default]
...
obs_tenant_id = <your-obs-tenant-id>
```
