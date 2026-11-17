# API catalog

Every NHN Cloud public API the CLI covers without a hand-written command is
described in `internal/catalog/data/<file>.json`. Each operation becomes
`nhncloud <service> <operation>`. Adding an API is adding catalog entries; no Go
code is needed.

After every change run both checks. `catalint` validates the files and every
catalog in the same directory together; the tests also check that catalog
names are canonical command names.

```sh
go run ./internal/catalog/catalint internal/catalog/data/*.json
go test ./cmd -run 'TestEmbeddedCatalogsAreValid|TestCatalogNamesAreCanonical'
```

## Rules

1. **Only what the API guide says.** Every method, path, parameter, field
   name, type, required flag and body encoding must come from the official API
   guide on docs.nhncloud.com. Do not guess field names or shapes; if the guide
   is unclear, leave the field out and note it in the catalog's operation
   summary as "body fields: see API guide" with an object-typed `--body`.
2. **Latest guide version.** Use the newest API version of a service. Put
   the guide URL in `docs` and the version in `version`.
3. **One file per guide**, named `<service>[-<area>].json`, kebab-case.
4. **Service name** (`name`) is the CLI command. Reuse an existing top-level
   command when the service already has one (run `nhncloud --help`), for
   example `network` for VPC, Floating IP and Security Groups,
   `loadbalancer`, `nks`, `rds-mysql`. Otherwise use a short kebab-case name
   of the product (`sms`, `email`, `cdn`, `gamebase`). `api`, `help` and
   `completion` are reserved. Use the canonical command name, not an alias.
   Names and `aliases` must not clash across catalogs, and only one catalog
   per service name sets `aliases`. Aliases are added only to commands the
   catalog creates: they are ignored when the service is a hand-written
   command, and an alias that matches an existing command is dropped.
5. **Operation names** are kebab-case verb-object: `list-<things>`,
   `get-<thing>`, `create-<thing>`, `update-<thing>`, `delete-<thing>`, or a
   specific verb (`attach-gateway`, `send-sms`). They must be unique across all
   catalogs that share the same service name, so include the resource in the
   name. Existing hand-written command names win and are skipped.
6. **No secrets or account data** in catalogs: no real appkeys, IDs, hosts
   other than the public endpoints from the guide, or e-mail addresses.

## Format

```json
{
  "name": "network",
  "title": "Network > Floating IP",
  "docs": "https://docs.nhncloud.com/ko/Network/Floating%20IP/ko/public-api/",
  "version": "v2.0",
  "iaasServiceType": "network",
  "auth": {"headers": {"X-Auth-Token": "{iaasToken}"}},
  "operations": [
    {
      "name": "create-floating-ip",
      "summary": "Create a floating IP",
      "method": "POST",
      "path": "/v2.0/floatingips",
      "params": [
        {"name": "floatingIpId", "in": "path", "type": "string", "required": true, "description": "..."}
      ],
      "body": {
        "wrapper": "floatingip",
        "fields": [
          {"name": "floating_network_id", "type": "string", "required": true, "description": "..."}
        ]
      }
    }
  ]
}
```

### Endpoint (exactly one)

| Key | When |
| --- | --- |
| `iaasServiceType` | IaaS APIs whose endpoint comes from the IaaS token's service catalog (`network`, `compute`, `image`, `volumev2`, `volumev3`, `container-infra`, `object-store`, ...). `path` starts after that endpoint. |
| `regionBaseURLs` | The guide lists a different URL per region: `{"kr1": "https://...", "kr2": "https://..."}`. Region keys are lower case (`kr1`, `kr2`, `kr3`, `jp1`, `us1`, ...). |
| `baseURL` | One URL; `{region}` (the only template allowed) is replaced by the region when the guide's URL contains the region. `{endpoint}` means the endpoint is shown in the console only and must be set by the user. |

Users can override the endpoint with `NHN_CLOUD_<NAME>_ENDPOINT` or
`<name>_endpoint` in the profile (`scheme://host[:port][/path]`, no user info,
query or fragment). `<NAME>` is the service name. When catalogs share a
service name but use different endpoints, give each a distinct
`"endpointName"` (kebab-case); the validator requires it. Example: the Load
Balancer secrets catalog uses `loadbalancer-secrets`, so
`NHN_CLOUD_LOADBALANCER_SECRETS_ENDPOINT` does not move other `loadbalancer`
calls.

### Auth

`auth.headers` maps a header name to a value template. Templates name
credentials in braces:

| Credential | Value |
| --- | --- |
| `iaasToken` | IaaS (Keystone) token |
| `userAccessToken` | User Access Key token, usually `"Bearer {userAccessToken}"` |
| `accessKeyId`, `secretAccessKey` | User Access Key ID and secret, sent as is |
| `appKey` | Service appkey (`--appkey`, `NHN_CLOUD_<SERVICE>_APPKEY`, `<service>_app_key` in the profile, the names the hand-written commands already use such as `NHN_CLOUD_MYSQL_APPKEY` and `rds_app_key`, then `NHN_CLOUD_APPKEY` and `app_key`) |
| `secretKey` | Service secret key (`--secret-key`, `NHN_CLOUD_<SERVICE>_SECRET_KEY`, `<service>_secret_key`) |
| `tenantId` | Project tenant ID |

An appkey or tenant ID in the path is a path param with `"credential"`
instead of a flag: `{"name": "appKey", "in": "path", "required": true, "credential": "appKey"}`.

### Parameters

Paths never end with `/`, even when the guide shows one: servers checked so far
reject the trailing slash.

`in` is `path`, `query` or `header`. Every `{var}` in `path` needs a path
param. The flag is the kebab-case name (`floating_network_id` becomes
`--floating-network-id`); set `"flag"` when it clashes with a global flag
(`region`, `appkey`, `password`, `tenant-id`, `profile`, `output`, `query`,
`debug`, `username`, `help`, `body`, `secret-key`, `data-file`, `header`,
`endpoint-override`) or with another parameter. Mark passwords and keys with
`"secret": true` so `-` reads them from stdin; the mark applies to that flag
only, and only to `string`, `integer` and `number` types. When it is the only
flag reading stdin, `-` takes all of stdin, so multi-line values such as PEM
keys arrive whole. Flag names are kebab-case. A name is declared
once per `in`, and a body field once. Header params cannot be `Content-Type`
(it comes from the body encoding) or a header that `auth` sets. Paths contain
no `%` and no braces outside `{var}`.

Param types: `string` (default), `integer`, `number`, `boolean`, `array`.
Numbers are checked before sending. An `array` query param is sent as a
repeated key (`--id a,b` or `--id a --id b` gives `id=a&id=b`). Path params
must be required, appear in the path, and cannot be empty, `.` or `..`.

Every operation also takes `--header Name:Value` (repeatable) for headers the
catalog does not list. It can set `Content-Type` except on multipart bodies,
whose boundary it would drop.

### Body

| `encoding` | Use |
| --- | --- |
| `json` (default) | JSON object; `wrapper` nests it under one key, and `--body` is the object inside the wrapper |
| `multipart` | `multipart/form-data`; `file` fields upload a file |
| `form` | `application/x-www-form-urlencoded` |
| `raw` | Bytes from `--data-file` (`-` for stdin) with `contentType`; no fields. `contentType` is only for raw bodies |

Field types: `string`, `integer`, `number`, `boolean`, `array` (set
`"items": "string"` for a flag), `object`, `file`. Only top-level fields are
listed; nested objects are `object` and come from `--body`. `required` is
checked at run time against flags and `--body`. `array` fields in `form` and
`multipart` bodies are sent one value per key; `object` fields are json only.

## Output

JSON responses are printed through `--output` and `--query`. A response with
`header.isSuccessful: false` is an error even when the status is 200. An
empty body prints the response headers, with credential-like headers masked
and URL queries dropped. Debug output and errors never show the query string
or credentials in the path. Redirects to another host are refused, so
credential headers never leave the API host.
