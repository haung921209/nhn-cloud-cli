# NHN Cloud CLI

A powerful, refactored Command Line Interface for NHN Cloud.

> **Status**: ✅ **Verified** (January 2026)

## Getting Started

### 1. Installation
```bash
go build -o nhncloud .
```

### 2. Configuration (Entry Point)
Before running commands, you **MUST** configure your credentials.
👉 **[Read Configuration Guide](docs/CONFIGURATION.md)**
*(Covers AppKeys, Tenant IDs, and Environment Variables)*

---

## Service Guides

Detailed usage instructions and command examples for each service:

| Category | Service | Documentation Link | Status |
|----------|---------|-------------------|--------|
| **Compute** | Nova | 📖 **[Compute Guide](docs/services/compute.md)** | Verified |
| **Network** | VPC | 📖 **[Network Guide](docs/services/network.md)** | Verified |
| **Container** | **NKS** (K8s) | 📖 **[NKS Guide (Golden Config)](docs/services/nks.md)** | **Verified** |
| | NCR/NCS | 📖 **[Container Guide](docs/services/container.md)** | Verified |
| **Database** | RDS | 📖 **[RDS Guide (MySQL/Maria/PG)](docs/services/rds.md)** | Verified |
| **Storage** | Object/NAS | 📖 **[Storage Guide](docs/services/storage.md)** | Verified |

---

## Verified Features
- **NKS**: Real Cluster Creation verified with Golden Configuration.
- **RDS**: Full CRUD for MySQL, MariaDB, PostgreSQL (including HA).
- **Storage**: Real 300GB NAS Volume and Object Storage Uploads verified.

## License
MIT
