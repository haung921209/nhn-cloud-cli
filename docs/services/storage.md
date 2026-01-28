# Storage Service Guide

This guide covers **Object Storage (OBS)** and **Network Attached Storage (NAS)**.

## 1. Object Storage (OBS)

Compatible with OpenStack Swift/S3. The CLI uses Swift protocol.

### Upload Object
Upload a local file to a container (bucket).
```bash
nhncloud object-storage put-object \
  --container my-container \
  --object images/logo.png \
  --file ./logo.png
```

### Download Object
```bash
nhncloud object-storage get-object \
  --container my-container \
  --object images/logo.png \
  --output-file ./downloaded-logo.png
```

### List Objects
```bash
nhncloud object-storage list-objects --container my-container
```

### Advanced Use Cases (Use Cases)

#### 1. 대용량 파일 업로드 (Multipart Upload)
5GB 이상의 대용량 파일은 자동으로 분할되어 업로드(SLO) 됩니다. `--segment-size` 옵션으로 분할 크기를 조정할 수 있습니다.

```bash
# 20MB 단위로 분할하여 업로드 (대용량 파일 권장)
nhncloud obs cp ./large_backup.tar.gz obs://my-backup-container/large_backup.tar.gz \
  --segment-size 20971520
```

#### 2. 디렉토리 전체 업로드/다운로드 (Recursive)
`-r` 또는 `--recursive` 옵션을 사용하여 폴더 구조를 유지한 채로 업로드하거나 다운로드할 수 있습니다.

```bash
# 로컬 디렉토리 전체를 Object Storage로 업로드
nhncloud obs cp ./my-data-folder obs://my-container/data-folder --recursive

# Object Storage의 특정 경로 하위를 로컬로 전체 다운로드
nhncloud obs cp obs://my-container/data-folder ./downloaded-data --recursive
```


---

## 2. NAS (Network Storage)

Managed NFS volumes for Compute instances.

### Create Volume
Volume size must be at least **300GB** (or 500GB depending on type).
```bash
nhncloud nas create-volume --name data-vol --size 300
```

### Delete Volume
```bash
nhncloud nas delete-volume --volume-id <volume-id>
```

### Snapshot Management
Create a point-in-time copy of your volume.
```bash
nhncloud nas create-nas-volume-snapshot --volume-id <volume-id> --name my-snap
```
