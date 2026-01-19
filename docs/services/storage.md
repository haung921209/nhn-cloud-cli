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
