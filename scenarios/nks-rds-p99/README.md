# nks-rds-p99

End-to-end NHN Cloud benchmark that provisions:

- a per-run **RDS** instance (MySQL / MariaDB / PostgreSQL — selectable via
  `--engine`), and
- an in-cluster **NKS Pod** that drives an OLTP load against the RDS
  `INTERNAL` endpoint,

then captures TPS / p50 / p95 / p99 / max-latency from inside the cluster
and tears down everything except the cluster itself.

The cluster is **owned outside the scenario** (Phase-4 pattern — see
`PHASE_4.md`). Create it once in the NHN console and pass its UUID with
`--existing-cluster <UUID>`. The scenario never creates or deletes a
cluster.

## Quick start

```bash
# 0. Create an NKS cluster from the console once. Note its UUID.

# 1. Fill in the tenant-specific UUIDs in scenario.yaml (every `FILL_ME_IN`
#    field — flavor_id, subnet_id, parameter group, external network, etc.).
#    See the inline comments in scenario.yaml for the matching describe-* verbs.

# 2. Run the scenario.
bash scenarios/nks-rds-p99/run.sh \
  --yes \
  --preset smoke \
  --engine mysql \
  --existing-cluster <CLUSTER-UUID>
```

`--yes` is required (acknowledges that billable RDS resources will be
created). `--keep` skips teardown if you want to inspect leftover state.

## Engines

| `--engine`   | RDS service     | Default port | Notes                                                                            |
|--------------|-----------------|-------------:|----------------------------------------------------------------------------------|
| `mysql`      | `rds-mysql`     |         3306 | Reference shape; loadgen uses `go-sql-driver/mysql`.                            |
| `mariadb`    | `rds-mariadb`   |         3306 | Wire-compatible with mysql driver. `--cidr` required at SG-create.              |
| `postgresql` | `rds-postgresql`|         5432 | Uses `pgx/v5`. SSL `prefer` by default. Master role is the benchmark identity.  |

## Prerequisites

- `nhncloud` CLI on PATH ([install guide](../../README.md))
- `kubectl`
- `go` (for the per-run cross-compile of the loadgen binary and the
  `nks-control` helper)
- `yq`, `jq` for shell parsing
- `~/.nhncloud/credentials` populated with `access_key_id`,
  `secret_access_key`, the per-engine app key (`rds_app_key` /
  `rds_mariadb_app_key` / `rds_postgresql_app_key`), `username`,
  `api_password`, `tenant_id`, `obs_tenant_id` (for cert-export to OBS)

## Presets

| preset   | sysbench threads | duration | tables | rows per table |
|----------|-----------------:|---------:|-------:|---------------:|
| `smoke`  | 8                |     120s |      4 |         10 000 |
| `proper` | 32               |     600s |      8 |        100 000 |

## What you get back

Each run writes to `reports/<RUN_TS>.*` next to this README:

- `<RUN_TS>.md`  — one-page summary table (TPS / p99 / errors)
- `<RUN_TS>.json` — machine-readable form of the same numbers
- `<RUN_TS>.oltp.txt` — full loadgen stdout (per-second log + final tally)
- `<RUN_TS>.create-db-instance.raw` — raw API response for the create
  call (kept on success and on failure for diagnostics)

`example-report/` ships completed reports from a development run on each
engine for quick comparison.

## Engine-specific gotchas (cliff notes)

Detailed retrospective lives in `PHASE_4.md`. The short version:

- **MariaDB**: `describe-db-flavors` (not `-classes`); `--cidr` is required
  at SG create time AND the initial ingress is auto-attached (the explicit
  `authorize-db-security-group-ingress` step is skipped). Many listed
  versions aren't creatable (`describe-db-engine-versions` is over-broad);
  pick a `default.MARIADB_V…` that actually appears in
  `describe-db-parameter-groups`. Cert-export is currently broken
  NHN-side — the scenario soft-falls to skip-verify TLS.
- **PostgreSQL**: flag names diverge (`--db-instance-name`, `--db-version`,
  `--db-user-name`, `--db-password`, `--database-name` REQUIRED at create);
  `--storage-type "General SSD"` is mandatory (the CLI's default `"SSD"`
  returns API 500); the smallest `m2.c1m2` flavor silently 500s — pick
  ≥ m2.c4m8; `endPointType` is plain `INTERNAL` (not `INTERNAL_VIP`); the
  master role lacks `CREATEROLE` so the benchmark runs as master; NHN PG's
  internal endpoint doesn't advertise SSL → loadgen defaults to
  `sslmode=prefer` and falls back to plain.

## Layout

```
nks-rds-p99/
├── run.sh                  # orchestrator
├── scenario.yaml           # tenant + preset configuration
├── steps/                  # one shell script per phase
├── loadgen/                # Go binary (cross-compiled to linux/amd64)
├── helpers/nks-control/    # Go binary that wraps the NKS SDK for cluster
│                           # lifecycle ops (the CLI doesn't expose all of
│                           # them today)
├── example-report/         # reference runs per engine
└── PHASE_4.md              # design + retrospective (rounds 3-5)
```

## Cleanup contract

By default the scenario tears down everything it created on this run:
RDS instance, RDS security group, compute security group, runner Pod,
keypair (if `01b-create-keypair.sh` ran). The cluster is owned outside
the scenario and is **never** deleted.

`--keep` skips teardown so you can inspect Pod logs or the RDS instance
afterward; the next clean-state run requires `bash steps/99-teardown.sh`
first to clear leftover resources.
