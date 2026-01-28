# Object Storage (OBS) 활용 사례 (Use Cases)

`nhncloud object-storage` (별칭 `obs`) 명령어를 사용하면 CLI에서 직접 객체(Object)와 컨테이너(Container)를 관리할 수 있습니다. 특히 Static Large Objects (SLO)를 통한 대용량 파일 자동 분할 업로드를 지원합니다.

## 1. 컨테이너 및 객체 조회 (List)

`ls` 명령어를 사용하여 컨테이너 목록이나 컨테이너 내부의 객체 목록을 조회할 수 있습니다.

### 모든 컨테이너 조회
```bash
nhncloud obs ls
```

### 객체 조회 (디렉토리 뷰)
기본적으로 `ls`는 **디렉토리 뷰**를 제공합니다 (폴더별로 그룹화하여 표시).

```bash
nhncloud obs ls obs://my-container/
# Output:
#                            PRE images/
#                            PRE data/
# file.txt      1024    2023-10-01T12:00:00
```

**재귀적 조회 (Recursive)**: 하위 폴더까지 플랫(Flat)하게 모두 조회하려면 `--recursive` (`-r`) 플래그를 사용합니다.

```bash
nhncloud obs ls -r obs://my-container/
# Output:
# images/logo.png
# data/log.txt
# file.txt
```

---

## 2. 파일 작업 (cp)

`cp` 명령어를 사용하여 로컬 머신과 Object Storage 간에 파일을 복사합니다. `--recursive` (`-r`) 플래그를 사용하면 디렉토리 전체를 복사할 수 있습니다.

### 업로드 (Upload)
```bash
# 단일 파일 업로드
nhncloud obs cp ./image.png obs://my-container/images/image.png

# 디렉토리 전체 업로드 (Recursive)
nhncloud obs cp -r ./local-dir obs://my-container/remote-dir

# 파일명 암묵적 지정 (디렉토리로 업로드)
nhncloud obs cp ./document.pdf obs://my-container/docs/
```

> [!NOTE]
> **대용량 파일 지원 (SLO)**
> 5GB보다 큰 파일은 **자동으로** 분할되어 Static Large Objects (SLO) 방식으로 업로드됩니다. 분할된 세그먼트는 `_segments`가 붙은 별도의 쉐도우 컨테이너(예: `my-container_segments`)에 저장됩니다.
> `--segment-size` 플래그로 분할 크기(기본값 1GB)를 조정할 수 있습니다.

```bash
# 10GB 파일 업로드 (자동으로 1GB 단위 분할)
nhncloud obs cp ./large-backup.tar.gz obs://backups/
```

### 다운로드 (Download)
```bash
# 단일 파일 다운로드
nhncloud obs cp obs://my-container/file.txt .

# 디렉토리 전체 다운로드 (Recursive)
nhncloud obs cp -r obs://my-container/remote-dir ./local-dir
```

### 컨테이너 간 복사 (Copy)
```bash
# 단일 파일 복사
nhncloud obs cp obs://src/file.txt obs://dst/file.txt

# 디렉토리 전체 복사 (Recursive)
nhncloud obs cp -r obs://src/folder obs://dst/backup-folder
```

---

## 3. 설정 (Configuration)

Object Storage 인증은 보통 전역 `tenant-id`를 따릅니다. 하지만 Object Storage 서비스가 다른 테넌트에 있는 경우(일부 조직 구성에서 발생), OBS 전용 Tenant ID를 설정할 수 있습니다.

```bash
nhncloud configure
# ...
# Object Storage Tenant ID [current-id]: <new-obs-tenant-id>
```

또는 `~/.nhncloud/credentials` 파일을 직접 수정하여 설정할 수 있습니다:
```ini
[default]
...
obs_tenant_id = <your-obs-tenant-id>
```
