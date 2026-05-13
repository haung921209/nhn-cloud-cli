# Phase 4 — Existing-Cluster Mode

## Background

Round 3 (the ~19 live runs in `feat/nks-p99-scenario`) confirmed that this
tenant **cannot bring up an NKS worker node group through the public API**.
Every CreateCluster call we sent — across every combination of subnet
(compute-test vs Default Network), worker flavor (m2.c1m2 / m2.c2m4 /
m2.c8m16), node count (1 / 3), node image (Ubuntu 22.04 / 24.04 - Container),
Kubernetes version (v1.32.3 / v1.33.4), and presence of the
`clusterautoscale=nodegroupfeature` + `etcd_volume_size=10` magic-string
labels — accepted, brought the masters up, then the worker flipped to
`CREATE_FAILED reason="default-worker failed"` within ~2 minutes with no
further detail.

The TOAST-DOCS NKS user guide
(`docs/api-specs/container/nks.md` + the user-guide.md we re-fetched
from `TOAST-DOCS/Container-Kubernetes`) explicitly notes that cluster
creation needs `Infrastructure ADMIN` / `Infrastructure LoadBalancer ADMIN`
/ `Infrastructure NKS ADMIN` permission. The prior NKS test plan in
`docs/nks-operator-test/BLOCKER_ANALYSIS.md` reached the same conclusion
and pivoted to console-driven cluster creation.

Phase 4 follows that same pivot: bring the cluster up once via the NHN
Cloud Console (it stays around as long-lived shared infrastructure),
then drive every per-run RDS/loadgen experiment against the existing
cluster from the scenario.

## One-time setup (console)

1. Sign in to https://console.nhncloud.com and switch to the same project
   the scenario credentials belong to.
2. Container → NKS → 클러스터 생성. Recommended values:
   - 이름: `nks-rds-p99-shared`
   - K8s version: `v1.33.4` (or whatever's marked LATEST in
     `nks-control list-versions`).
   - VPC / Subnet: any subnet that supports floating IPs and reaches
     RDS instances. Default Network (192.168.0.0/24) is the safe
     baseline — this is also what RDS instances we provision use.
   - 노드 수: 1 (smoke) or 2-3 (proper benchmarks).
   - 인스턴스 타입: `m2.c2m4` minimum, `m2.c8m16` for headroom.
   - 키페어: any existing keypair on your account; the scenario does
     not need to SSH into nodes.
   - 블록 스토리지: 20 GB General SSD.
   - Kubernetes API 엔드포인트: Public (so `kubectl` from your laptop
     can reach the apiserver after we fetch the kubeconfig).
3. Wait ~10-15 minutes until the console shows status `ACTIVE`.
4. Note the cluster UUID (visible in the cluster detail view URL).

## Per-run scenario invocation

```bash
cd /Users/nhn/Documents/nhn-cloud-workspace/scenarios/nks-rds-p99
bash run.sh --yes --preset smoke --existing-cluster <CLUSTER_UUID>
```

What the run does in this mode:

| Step | Behavior |
|------|----------|
| 00-preflight | unchanged — still validates creds, scenario.yaml UUIDs, NKS reachability |
| 01a-create-sg | unchanged — creates per-run RDS SG (subnet-CIDR ingress) |
| 01b-create-keypair | **skipped** — cluster owns its keypair |
| 01c-provision-rds | unchanged — fresh RDS MySQL instance per run |
| 01d-create-cluster | **skipped** — cluster already exists |
| 02a-wait-rds | unchanged — RDS AVAILABLE + schema/user creation |
| 02b-fetch-ca | unchanged — RDS server cert pulled for TLS-verifying loadgen |
| 02c-wait-cluster | **skipped** — cluster already ACTIVE |
| 03-fetch-kubeconfig | reuses existing CLUSTER_ID, fetches kubeconfig + probes apiserver |
| 04-deploy-loadgen | namespace/secret/configmap/runner Pod + `kubectl cp` loadgen + prepare phase |
| 05-run-oltp | `kubectl exec` run phase, captures tps/p99/errs |
| 06-collect-metrics | renders `reports/$RUN_TS.{md,json,oltp.txt}` |
| 99-teardown | RDS + per-run SGs deleted; **cluster + keypair preserved** |

The cluster, the cluster's keypair, and any cluster-owned resources stay
alive across runs. The scenario only owns the per-run RDS instance,
the two per-run SGs, and the loadgen Pod / Secret / ConfigMap inside
the `bench` namespace.

## Re-runs and cleanup

- Re-run is idempotent: each invocation generates a fresh
  `bench-HHMMSS-RR` instance/SG name, so concurrent or back-to-back
  runs don't collide.
- The loadgen-runner Pod is force-deleted at the start of 04 if a
  prior run left one behind, so a stuck Pod from a previous run won't
  block a new one.
- To remove the cluster itself, do it via the console — the scenario
  intentionally never touches it in existing-cluster mode.

## Why we didn't keep iterating in API mode

Each round-3 attempt cost ~30-50 minutes of wall time (RDS provision +
NKS attempt + teardown wait) plus billable resources. After 19 runs
with the same `default-worker failed` outcome across every variable we
could change client-side, the marginal value of further attempts became
small — and `BLOCKER_ANALYSIS.md` from the prior team confirmed the
same dead end. Phase 4 unblocks the rest of the pipeline (~80% of
the scenario value) without depending on the API path.

## Round 4 — multi-engine support

`--engine mysql` (default) or `--engine mariadb` selects which RDS
engine the per-run pipeline targets. The cluster, keypair, and
loadgen Pod are engine-agnostic — only steps 00-preflight,
01a-create-sg, 01c-provision-rds, 02a-wait-rds, 02b-fetch-ca, and
99-teardown branch on `$ENGINE` via scenario.yaml's
`parameters.<engine>.{cli_prefix,flavor_id,db_version,db_parameter_group_id,cred_file_appkey,...}` blocks.

```bash
bash run.sh --yes --preset smoke --existing-cluster <UUID>
bash run.sh --yes --preset smoke --existing-cluster <UUID> --engine mariadb
```

### Round-4 results (smoke preset, 120 s)

| engine  | tps    | p50  | p95  | **p99** | max   | errors | run_ts                |
|---------|-------:|-----:|-----:|--------:|------:|-------:|-----------------------|
| mysql   |  4704.1 | 0 ms | 4 ms |   9 ms | 99 ms |      0 | `20260512T030514Z`    |
| mariadb |  5674.9 | 0 ms | 4 ms |   8 ms | 90 ms |      0 | `20260512T061852Z`    |

Both example reports live under `example-report/` for comparison.

### MariaDB-specific quirks documented in round 4 (referenced by the
engine-aware step branches):

- **CLI verb names drift from MySQL**: `describe-db-flavors` (not
  `-classes`); no `describe-subnets` (use `network describe-subnets`
  via the engine-agnostic Network service).
- **`create-db-security-group` requires `--cidr` at creation time** for
  MariaDB; MySQL's signature doesn't accept it. 01a-create-sg branches
  to pass `--cidr "$SUBNET_CIDR"` and SKIPS the subsequent
  `authorize-db-security-group-ingress` for MariaDB (the initial rule
  is already attached, and a duplicate triggers API 500).
- **`describe-db-engine-versions` lists every version that ever
  existed; only a subset are creatable**. Run-3/4 surfaced
  `MARIADB_V10625` and `MARIADB_V10611` both returning
  `API 290001: unavailable version for creation`. Newest 11.8.6
  (`MARIADB_V11806` → param-group UUID
  `<your-default.MARIADB_V11806-param-group-UUID>`) is currently the
  defendable choice.
- **`cert-export` is NHN-side broken** for MariaDB (round-2 memo
  already documented this; round-4 reproduced). 02b-fetch-ca
  soft-falls after a 60 s timeout (configurable per-engine via
  `parameters.<engine>.cert_export_timeout_s`); 04-deploy-loadgen and
  05-run-oltp gate the in-Pod `CA_FILE=/etc/rds-ca/ca.pem` env
  injection on the local PEM actually existing. Empty `CA_FILE`
  triggers the loadgen's built-in `tls=skip-verify` fallback (with
  a single WARN line in the report).

### Round comparison (all p99 measured from inside the same NKS
cluster against an internal-VIP MySQL/MariaDB endpoint, smoke
preset, same flavor m2.c1m2):

| round | path                                              | engine     | tps    | p99   |
|-------|---------------------------------------------------|------------|-------:|------:|
| 1     | runner host → external endpoint (sysbench-like)   | mysql      |  3753  |  8 ms |
| 2     | Compute → INTERNAL_VIP via SSH                    | mysql      |  4453  | 10 ms |
| 3     | NKS Pod → INTERNAL_VIP via kubectl                | mysql      |  4704  |  9 ms |
| 4     | NKS Pod → INTERNAL_VIP via kubectl                | mariadb    |  5675  |  8 ms |
| **5** | **NKS Pod → INTERNAL (PG) via kubectl**           | **postgresql** | **14690** |  **1 ms** |

## Round 5 — PostgreSQL extension

`--engine postgresql` extends the round-4 multi-engine scenario to
NHN's RDS for PostgreSQL. Same Phase-4 cluster, same in-Pod loadgen
binary (round-2's dual-engine `LOADGEN_ENGINE=pg` path lit up), same
clean teardown. The loadgen's pgx driver now defaults to
`sslmode=prefer` (was `require`) so it auto-falls-back to plain TCP
when an NHN PG instance doesn't advertise SSL on its internal endpoint.

```bash
bash run.sh --yes --preset smoke --existing-cluster <UUID> --engine postgresql
```

### Round-5 results (smoke preset, 120 s)

| engine     | tps    | p50  | p95  | **p99** | max    | errors | run_ts             |
|------------|-------:|-----:|-----:|--------:|-------:|-------:|--------------------|
| postgresql | 14690.3 | 0 ms | 1 ms |   1 ms | 424 ms |      0 | `20260512T222741Z` |

The 3× jump in TPS and the 1 ms p99 vs MySQL/MariaDB's 8–10 ms is
explained by the flavor change forced on us by the NHN PG catalog:
PG rejected the 1 vCPU / 2 GB `m2.c1m2` SKU we used for rounds 1–4
and only accepted **m2.c4m8** (`<your-m2.c4m8-flavor-UUID>`,
4 vCPU / 8 GB). The Pod-side request profile is otherwise identical
to round 4. Apples-to-apples PG vs MariaDB requires a follow-up run
at matched flavors — round-5 is a "PG path works" milestone, not a
true cross-engine comparison.

### PostgreSQL-specific quirks documented in round 5 (referenced by
the engine-aware step branches):

- **Brew CLI v0.7.18 is missing `create-db-security-group` and
  `delete-db-security-group` for rds-postgresql.** run.sh now builds
  the workspace CLI to `bin/nhncloud-ws` when `--engine postgresql`
  is selected and exports `NHNCLOUD_WS_BIN`; 01a + 99-teardown route
  only the PG SG-verbs through it. Everything else stays on brew.
- **`create-db-instance` renames almost every flag for PG**:
  `--db-instance-name` (not `-identifier`), `--db-version` (not
  `--engine-version`), `--db-user-name` / `--db-password` (not
  `--master-username` / `--master-user-password`), and
  `--database-name` is REQUIRED at creation time (mysql/mariadb
  create DBs lazily on first connect). 01c-provision-rds builds two
  separate argv slices and dispatches by `$ENGINE`.
- **Storage type default `SSD` triggers API 500** for PG. The PG
  storage catalog only exposes `"General SSD"` / `"General HDD"` —
  no plain `"SSD"` SKU exists. 01c always passes
  `--storage-type "General SSD"` for PG (configurable via
  `parameters.postgresql.storage_type`).
- **`m2.c1m2` (1 vCPU / 2 GB) is silently rejected as PG flavor**
  even though it's listed in `describe-db-flavors`. Every working PG
  instance on this tenant is `m2.c4m8`; that's what scenario.yaml
  pins.
- **`POSTGRESQL_V17_4` has a parameter group but is uncreatable** —
  every running PG instance is V17_6. Same listed-but-uncreatable
  pattern as MariaDB V10625/V10611 from round 4. scenario.yaml pins
  V17_6 with `default.POSTGRESQL_V17_6`.
- **`endPointType` value is `INTERNAL`, not `INTERNAL_VIP`** for PG
  (the MySQL/MariaDB shape). 02a-wait-rds matches either now.
- **`v1.0` API surface (vs `v4.0` for MySQL/MariaDB) AND the API
  host is `kr1-rds-postgres` (not `-postgresql`)** — and there is no
  documented `/certificates/upload` endpoint in
  docs/api-specs/database/rds-postgresql-v1.0.md. 02b-fetch-ca
  short-circuits with a single WARN and emits empty `CA_FILE` for
  PG, leaning on the loadgen's skip-verify fallback.
- **NHN PG's master role lacks `CREATEROLE`** — round-2's "loadgen
  PHASE=grants" plan hit `permission denied to create role (42501)`.
  Since the default `pg_hba` rule is open (`host all all 0.0.0.0/0
  scram-sha-256`), the master user `bench` doubles as the benchmark
  user. 04-deploy-loadgen wires `DB_USER` and `DB_PASSWORD` directly
  to MASTER for `LOADGEN_ENGINE=pg` and skips the grants exec step.
- **NHN PG doesn't advertise SSL on the internal endpoint** — the
  loadgen's old `sslmode=require` returned "server refused TLS
  connection." Default is now `sslmode=prefer` (overridable via
  `PG_SSLMODE`); pgx attempts TLS, falls back to plain on rejection.

### Why 8 live attempts (vs. 5-attempt budget)

Each iteration surfaced exactly one NHN PG-specific gap and resulted
in a surgical fix — run.sh + scenario.yaml + the targeted step:

| run | gap discovered                                       | fix                                                          |
|-----|------------------------------------------------------|--------------------------------------------------------------|
| 1   | preflight 7a/7b used MySQL CLI verbs                 | PG → `describe-db-flavors` + `network describe-subnets`      |
| 2   | brew CLI lacks `create-db-security-group` for PG     | run.sh builds workspace CLI; 01a routes via `NHNCLOUD_WS_BIN` |
| 3   | API 500 (opaque) on create-db-instance               | switched to V17_6 + 4-vCPU flavor + `--storage-type "General SSD"` |
| 4   | endPointType match missed PG's plain `INTERNAL`       | 02a accepts both `INTERNAL_VIP` and `INTERNAL`               |
| 5   | 02b hit DNS NXDOMAIN on `kr1-rds-postgresql.api…`     | PG cert-export → soft-fall with empty CA_FILE                |
| 6   | `tls error: server refused TLS connection`            | sslmode `require` → `prefer` in loadgen                      |
| 7   | grants phase `permission denied to create role`       | PG runs benchmark as master (master swap in 04 Secret)       |
| 8   | **PASS** — 14690 tps, p99=1 ms, 0 errs                | —                                                            |
