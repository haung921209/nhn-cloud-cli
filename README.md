# NHN Cloud CLI

NHN Cloud 서비스를 위한 강력하고 새롭게 리팩토링된 Command Line Interface (CLI) 도구입니다.
터미널 환경에서 NHN Cloud의 다양한 리소스(Compute, Network, Database, Storage 등)를 쉽고 빠르게 관리할 수 있도록 설계되었습니다.

> **현재 상태**: ✅ **검증 완료** (2026년 1월)

## 🚀 시작하기 (Getting Started)

### 1. 설치 (Installation)
소스 코드를 빌드하여 실행 파일을 생성합니다.

```bash
go build -o nhncloud .
```

### 2. 환경 설정 (Configuration)
CLI를 사용하기 위해서는 인증 정보 설정이 필수적입니다.
**[설정 가이드 (Configuration Guide)](docs/CONFIGURATION.md)** 문서를 참고하여 AppKey, Tenant ID, 환경 변수 등을 올바르게 설정해주세요.

설정은 다음 순서로 우선순위가 적용됩니다:
1. **CLI 플래그** (예: `--region`, `--appkey`)
2. **환경 변수** (예: `NHN_CLOUD_APPKEY`)
3. **설정 파일** (`~/.nhncloud/credentials`)

---

## 📚 서비스 가이드 (Service Guides)

각 서비스별 상세 사용법과 명령어 예제는 아래 문서에서 확인할 수 있습니다.

| 카테고리 | 서비스 | 문서 링크 | 상태 |
|----------|---------|-------------------|--------|
| **Compute** | Nova | 📖 **[컴퓨트 가이드 (Compute)](docs/services/compute.md)** | 검증 완료 |
| **Network** | VPC | 📖 **[네트워크 가이드 (Network)](docs/services/network.md)** | 검증 완료 |
| **Container** | **NKS** (K8s) | 📖 **[NKS 가이드 (Golden Config)](docs/services/nks.md)** | **검증 완료** |
| | NCR/NCS | 📖 **[컨테이너 레지스트리/서비스 가이드](docs/services/container.md)** | 검증 완료 |
| **Database** | RDS | 📖 **[RDS 가이드 (MySQL/Maria/PG)](docs/services/rds.md)** | 검증 완료 |
| **Storage** | Object/NAS | 📖 **[스토리지 가이드 (Storage)](docs/services/storage.md)** | 검증 완료 |

---

## ✅ 검증된 기능 (Verified Features)

현재 버전에서 다음과 같은 주요 기능들이 실제 환경에서 테스트되고 검증되었습니다.

- **NKS (Kubernetes)**: Golden Configuration을 적용한 실제 클러스터 생성 및 관리 기능 검증.
- **RDS (Database)**: MySQL, MariaDB, PostgreSQL의 전체 CRUD 작업 및 고가용성(HA) 기능 검증.
- **Storage**: 실제 300GB NAS 볼륨 생성 및 Object Storage 파일 업로드/다운로드 테스트 완료.

## 📄 라이선스 (License)
MIT License
