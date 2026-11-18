# API catalog status

Which catalog operations were run against the live API. Last run: 2026-11-18, one
test project, region kr1, read-only (`GET`) operations only. Create, update
and delete operations (957 of 1632) were not run.

How the run works: every `GET` operation that needs no input is called; then
`GET` operations whose path IDs can be taken from the list operation one path
level up. Downloads are skipped. Operations with the same name as a
hand-written command are not counted (57); the hand-written command
runs instead.

| Column | Meaning |
| --- | --- |
| Operations | all catalog operations of the service, including ones a hand-written command replaces |
| GET | catalog `GET` operations (hand-written ones excluded) |
| OK | answered successfully |
| Failed | the API returned an error; see the note |
| Not run | needed an ID the project does not have, a value only a person can give (search terms, dates), or a secret key the project does not have |

Totals: 618 `GET` operations, 176 OK, 36 failed, 406 not run.

A service marked "guide only" follows the API guide but has not been called.
Please open an issue when a command and the API disagree.

| Service | Operations | GET | OK | Failed | Not run | State | Note |
| --- | --- | --- | --- | --- | --- | --- | --- |
| apigw | 60 | 18 | 0 | 1 | 17 | failed | metrics call returned Permission denied; other reads need IDs the test project does not have |
| appguard | 1 | 1 | 0 | 0 | 1 | guide only |  |
| autocomplete | 7 | 2 | 0 | 0 | 2 | guide only |  |
| block-storage | 12 | 6 | 5 | 0 | 1 | live read checked |  |
| captcha | 5 | 3 | 0 | 0 | 3 | guide only |  |
| cdn | 16 | 3 | 0 | 3 | 0 | failed | server rejected the project's appkey (service not enabled in the test project) |
| certmanager | 4 | 2 | 0 | 0 | 2 | guide only |  |
| cloud-foundry | 14 | 1 | 0 | 0 | 1 | guide only |  |
| cloud-functions | 15 | 6 | 0 | 2 | 4 | failed | server rejected the project's appkey (service not enabled in the test project) |
| cloud-search | 8 | 3 | 0 | 0 | 3 | guide only |  |
| cloudtrail | 1 | 0 | 0 | 0 | 0 | guide only |  |
| colocation-gateway | 2 | 2 | 1 | 0 | 1 | live read checked |  |
| compute | 34 | 14 | 11 | 0 | 3 | live read checked |  |
| corporation-search | 5 | 4 | 0 | 0 | 4 | guide only |  |
| deploy | 9 | 7 | 0 | 1 | 6 | failed | server rejected the project's appkey (service not enabled in the test project) |
| dnsplus | 24 | 2 | 2 | 0 | 0 | live read checked |  |
| easyqueue | 8 | 5 | 0 | 0 | 5 | guide only |  |
| email | 35 | 17 | 0 | 0 | 17 | guide only |  |
| face-recognition | 17 | 4 | 0 | 0 | 4 | guide only |  |
| file-crafter | 2 | 0 | 0 | 0 | 0 | guide only |  |
| flow-log | 7 | 0 | 0 | 0 | 0 | guide only |  |
| gamebase | 65 | 28 | 0 | 0 | 28 | guide only |  |
| iam | 75 | 36 | 5 | 0 | 31 | live read checked |  |
| iap | 3 | 0 | 0 | 0 | 0 | guide only |  |
| image | 14 | 3 | 0 | 0 | 3 | guide only |  |
| image-manager | 16 | 7 | 0 | 2 | 5 | failed | server rejected the project's appkey (service not enabled in the test project) |
| internet-gateway | 4 | 0 | 0 | 0 | 0 | guide only |  |
| kakao | 91 | 41 | 0 | 0 | 41 | guide only |  |
| key-manager | 40 | 7 | 0 | 0 | 7 | guide only |  |
| launching | 1 | 1 | 0 | 1 | 0 | failed | server rejected the project's appkey (service not enabled in the test project) |
| leaderboard | 15 | 7 | 0 | 1 | 6 | failed | server rejected the project's appkey (service not enabled in the test project) |
| loadbalancer | 65 | 25 | 10 | 0 | 15 | live read checked |  |
| log-crash-search | 7 | 1 | 0 | 0 | 1 | guide only |  |
| mirroring | 15 | 6 | 3 | 0 | 3 | live read checked |  |
| nas | 20 | 4 | 0 | 0 | 4 | guide only |  |
| nat-gateway | 5 | 0 | 0 | 0 | 0 | guide only |  |
| ncr | 52 | 23 | 5 | 0 | 18 | live read checked | get-registry accepts the registry name; the ID from list-registries returns 'Could not find the registry' |
| ncs | 25 | 9 | 1 | 0 | 8 | live read checked |  |
| network | 42 | 16 | 15 | 0 | 1 | live read checked |  |
| network-acl | 14 | 2 | 1 | 0 | 1 | live read checked |  |
| nks | 37 | 16 | 6 | 0 | 10 | live read checked |  |
| notification-hub | 95 | 38 | 0 | 14 | 24 | failed | list calls return HTTP 500 (internal server error); probably the service is not enabled in the test project |
| object-storage | 16 | 3 | 0 | 0 | 3 | guide only |  |
| ocr | 10 | 1 | 0 | 0 | 1 | guide only |  |
| pipeline | 1 | 0 | 0 | 0 | 0 | guide only |  |
| private-ca | 22 | 11 | 0 | 1 | 10 | failed | server rejected the project's appkey (service not enabled in the test project) |
| private-dns | 10 | 4 | 1 | 0 | 3 | live read checked |  |
| push | 35 | 18 | 7 | 0 | 11 | live read checked |  |
| rcs-bizmessage | 13 | 8 | 0 | 0 | 8 | guide only |  |
| rds-mariadb | 97 | 39 | 31 | 1 | 7 | live read checked |  |
| rds-mysql | 97 | 37 | 31 | 1 | 5 | live read checked |  |
| rds-postgresql | 105 | 37 | 29 | 0 | 8 | live read checked | get-db-instance-restoration-info did not answer within the CLI's 60 s timeout |
| resource-watcher | 6 | 1 | 0 | 0 | 1 | guide only |  |
| role | 71 | 15 | 0 | 0 | 15 | guide only |  |
| rtcs | 5 | 3 | 0 | 1 | 2 | failed | the guide's host api-gw.cloud.toast.com has no address in public DNS |
| s3-credential | 3 | 0 | 0 | 0 | 0 | guide only |  |
| security-advisor | 3 | 3 | 0 | 2 | 1 | failed | the guide's host advisor.api.nhncloudservice.com does not resolve in public DNS |
| security-monitoring | 8 | 6 | 0 | 5 | 1 | failed | server rejected the project's appkey (service not enabled in the test project) |
| server-security-check | 2 | 2 | 0 | 0 | 2 | guide only |  |
| service-gateway | 19 | 5 | 3 | 0 | 2 | live read checked |  |
| service-monitoring | 6 | 1 | 0 | 0 | 1 | guide only |  |
| shorturl | 3 | 2 | 0 | 0 | 2 | guide only |  |
| sms | 60 | 32 | 0 | 0 | 32 | guide only |  |
| speech-to-text | 4 | 2 | 0 | 0 | 2 | guide only |  |
| text-to-speech | 1 | 0 | 0 | 0 | 0 | guide only |  |
| transit-hub | 48 | 18 | 9 | 0 | 9 | live read checked |  |
