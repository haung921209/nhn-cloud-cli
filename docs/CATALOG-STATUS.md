# API catalog status

Which catalog operations were run against the live API. Last run: 2026-11-18, one
test project, region kr1. There were two runs.

Read run: every `GET` operation that needs no input is called; then
`GET` operations whose path IDs can be taken from the list operation one path
level up. Downloads are skipped. Operations with the same name as a
hand-written command are not counted (57); the hand-written command
runs instead.

Write run: scenarios that create a resource, read it, change it and delete
it, with the service's product enabled. 187 of the 957 non-`GET`
operations succeeded. It did not run operations that a hand-written command
replaces, operations that need an instance, cluster or load balancer, and
operations whose input can only be set up in the console. Sending messages
was not tested.

| Column | Meaning |
| --- | --- |
| Operations | all catalog operations of the service, including ones a hand-written command replaces |
| GET | catalog `GET` operations (hand-written ones excluded) |
| OK | answered successfully |
| Failed | the API returned an error; see the note |
| Not run | needed an ID the project does not have, a value only a person can give (search terms, dates), or a secret key the project does not have |
| Write OK | non-`GET` operations that succeeded in the write run |
| State | read: a read succeeded; write: a write succeeded |

Read totals: 618 `GET` operations, 176 OK, 36 failed, 406 not run.

A service marked "guide only" follows the API guide but has not been called.
Please open an issue when a command and the API disagree.

| Service | Operations | GET | OK | Failed | Not run | Write OK | State | Note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| apigw | 60 | 18 | 0 | 1 | 17 | 0 | failed | metrics call returned Permission denied; other reads need IDs the test project does not have |
| appguard | 1 | 1 | 0 | 0 | 1 | 0 | guide only |  |
| autocomplete | 7 | 2 | 0 | 0 | 2 | 0 | guide only |  |
| block-storage | 12 | 6 | 5 | 0 | 1 | 0 | read |  |
| captcha | 5 | 3 | 0 | 0 | 3 | 2 | write |  |
| cdn | 16 | 3 | 0 | 3 | 0 | 4 | write | read run: appkey rejected (product not enabled then); the write run with the product enabled passed |
| certmanager | 4 | 2 | 0 | 0 | 2 | 0 | guide only |  |
| cloud-foundry | 14 | 1 | 0 | 0 | 1 | 0 | guide only |  |
| cloud-functions | 15 | 6 | 0 | 2 | 4 | 0 | failed | server rejected the project's appkey (service not enabled in the test project) |
| cloud-search | 8 | 3 | 0 | 0 | 3 | 0 | guide only |  |
| cloudtrail | 1 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| colocation-gateway | 2 | 2 | 1 | 0 | 1 | 0 | read |  |
| compute | 34 | 14 | 11 | 0 | 3 | 4 | read; write |  |
| corporation-search | 5 | 4 | 0 | 0 | 4 | 0 | guide only | create/update/delete run: the guide's recent-request paths return 404 |
| deploy | 9 | 7 | 0 | 1 | 6 | 0 | failed | server rejected the project's appkey (service not enabled in the test project); create/update/delete run: the API rejected the product's appkey |
| dnsplus | 24 | 2 | 2 | 0 | 0 | 0 | read |  |
| easyqueue | 8 | 5 | 0 | 0 | 5 | 3 | write |  |
| email | 35 | 17 | 0 | 0 | 17 | 10 | write |  |
| face-recognition | 17 | 4 | 0 | 0 | 4 | 2 | write |  |
| file-crafter | 2 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| flow-log | 7 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| gamebase | 65 | 28 | 0 | 0 | 28 | 0 | guide only |  |
| iam | 75 | 36 | 5 | 0 | 31 | 0 | read |  |
| iap | 3 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| image | 14 | 3 | 0 | 0 | 3 | 0 | guide only |  |
| image-manager | 16 | 7 | 0 | 2 | 5 | 5 | write | read run: appkey rejected (product not enabled then); the write run with the product enabled passed |
| internet-gateway | 4 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| kakao | 91 | 41 | 0 | 0 | 41 | 2 | write |  |
| key-manager | 40 | 7 | 0 | 0 | 7 | 22 | write |  |
| launching | 1 | 1 | 0 | 1 | 0 | 0 | failed | server rejected the project's appkey (service not enabled in the test project) |
| leaderboard | 15 | 7 | 0 | 1 | 6 | 0 | failed | server rejected the project's appkey (service not enabled in the test project) |
| loadbalancer | 65 | 25 | 10 | 0 | 15 | 0 | read |  |
| log-crash-search | 7 | 1 | 0 | 0 | 1 | 0 | guide only |  |
| mirroring | 15 | 6 | 3 | 0 | 3 | 6 | read; write |  |
| nas | 20 | 4 | 0 | 0 | 4 | 0 | guide only |  |
| nat-gateway | 5 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| ncr | 52 | 23 | 5 | 0 | 18 | 0 | read | get-registry accepts the registry name; the ID from list-registries returns 'Could not find the registry' |
| ncs | 25 | 9 | 1 | 0 | 8 | 4 | read; write |  |
| network | 42 | 16 | 15 | 0 | 1 | 4 | read; write |  |
| network-acl | 14 | 2 | 1 | 0 | 1 | 0 | read |  |
| nks | 37 | 16 | 6 | 0 | 10 | 0 | read |  |
| notification-hub | 95 | 38 | 0 | 14 | 24 | 0 | failed | list calls return HTTP 500 (internal server error); probably the service is not enabled in the test project; create/update/delete run: 404 'no identity-verification record' for the project |
| object-storage | 16 | 3 | 0 | 0 | 3 | 0 | guide only |  |
| ocr | 10 | 1 | 0 | 0 | 1 | 0 | guide only |  |
| pipeline | 1 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| private-ca | 22 | 11 | 0 | 1 | 10 | 10 | write | read run: appkey rejected (product not enabled then); the write run with the product enabled passed |
| private-dns | 10 | 4 | 1 | 0 | 3 | 6 | read; write |  |
| push | 35 | 18 | 7 | 0 | 11 | 13 | read; write |  |
| rcs-bizmessage | 13 | 8 | 0 | 0 | 8 | 0 | guide only | create/update/delete run: upload-media returned resultCode 500 (the project has no RCS brand) |
| rds-mariadb | 97 | 39 | 31 | 1 | 7 | 9 | read; write |  |
| rds-mysql | 97 | 37 | 31 | 1 | 5 | 9 | read; write |  |
| rds-postgresql | 105 | 37 | 29 | 0 | 8 | 6 | read; write | get-db-instance-restoration-info did not answer within the CLI's 60 s timeout |
| resource-watcher | 6 | 1 | 0 | 0 | 1 | 0 | guide only |  |
| role | 71 | 15 | 0 | 0 | 15 | 53 | write |  |
| rtcs | 5 | 3 | 0 | 1 | 2 | 0 | failed | the guide's host api-gw.cloud.toast.com has no address in public DNS |
| s3-credential | 3 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| security-advisor | 3 | 3 | 0 | 2 | 1 | 0 | failed | the guide's host advisor.api.nhncloudservice.com does not resolve; the catalog uses security-advisor.api.nhncloudservice.com |
| security-monitoring | 8 | 6 | 0 | 5 | 1 | 0 | failed | server rejected the project's appkey (service not enabled in the test project) |
| server-security-check | 2 | 2 | 0 | 0 | 2 | 0 | guide only |  |
| service-gateway | 19 | 5 | 3 | 0 | 2 | 0 | read |  |
| service-monitoring | 6 | 1 | 0 | 0 | 1 | 0 | guide only | create/update/delete run: the API rejected the product's appkey |
| shorturl | 3 | 2 | 0 | 0 | 2 | 1 | write |  |
| sms | 60 | 32 | 0 | 0 | 32 | 11 | write |  |
| speech-to-text | 4 | 2 | 0 | 0 | 2 | 1 | write |  |
| text-to-speech | 1 | 0 | 0 | 0 | 0 | 0 | guide only |  |
| transit-hub | 48 | 18 | 9 | 0 | 9 | 0 | read |  |
