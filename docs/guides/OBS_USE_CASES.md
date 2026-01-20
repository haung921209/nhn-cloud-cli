# Object Storage (OBS) Use Cases

The `nhncloud object-storage` (alias `obs`) command allows you to manage objects and containers directly from the CLI. It supports efficient file transfers, including automatic handling of large files via Static Large Objects (SLO).

## 1. List Containers and Objects

Use the `ls` command to view your containers or lists objects within them.

### List All Containers
```bash
nhncloud obs ls
```
**Output:**
```
my-container    1024 bytes    5 objects
backup-data     500 GB        120 objects
```

### List Objects in a Container
```bash
nhncloud obs ls obs://my-container
```
**Output:**
```
image.jpg       2048    2023-10-01T12:00:00
data/log.txt    500     2023-10-02T10:00:00
```

---

## 2. File Operations (cp)

The `cp` command copies files between your local machine and Object Storage, or between two locations in Object Storage. It uses the `obs://<container>/<key>` URI scheme.

### Upload Files
Upload a local file to a container.

```bash
# Simple upload
nhncloud obs cp ./image.png obs://my-container/images/image.png

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

### Download Files
Download a specific object to your local machine.

```bash
# Download to current directory
nhncloud obs cp obs://my-container/images/image.png .

# Download with specific filename
nhncloud obs cp obs://my-container/config.json ./local-config.json
```

### Copy Between Containers
Copy objects directly between containers (server-side copy).

```bash
nhncloud obs cp obs://source-container/file.txt obs://dest-container/file.txt
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
