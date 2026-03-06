# MASKIT — Test Scenarios

---

## Scenarios

Each scenario is a **situation** (who, what, outcome). When you add or change behavior, add or edit a scenario and keep it in sync with `mask_test.go` and manual checks.

### 1.1 No sensitive content


| Scenario                | Given         | Expected      | Test in `TestMask`                   | Manual check                    |
| ----------------------- | ------------- | ------------- | ------------------------------------ | ------------------------------- |
| User pastes plain text  | `hello world` | `hello world` | `plain text unchanged`, `no matches` | `echo "hello world" | ./maskit` |
| User pastes empty input | ``            | ``            | `empty string`                       | `echo -n "" | ./maskit`         |


### 1.2 Single kind of sensitive content


| Scenario                     | Given                                    | Expected                    | Test in `TestMask` | Manual check                                  |
| ---------------------------- | ---------------------------------------- | --------------------------- | ------------------ | --------------------------------------------- |
| User pastes one email        | `contact me at a@b.co`                   | `contact me at {{EMAIL:1}}` | `single email`     | `echo "contact alice@example.com" | ./maskit` |
| User pastes one HTTPS URL    | `see https://api.foo.com`                | `see {{URL:1}}`             | `single URL`       |                                               |
| User pastes one HTTP URL     | `visit http://foo.org/path`              | `visit {{URL:1}}`           | `http URL`         |                                               |
| User pastes one FTP URL      | `fetch from ftp://files.company.com/pub` | `fetch from {{URL:1}}`      | `ftp URL`          |                                               |
| User pastes one FTPS URL     | `secure ftps://secure.company.com/data`  | `secure {{URL:1}}`          | `ftps URL`         |                                               |
| User pastes one IPv4 address | `connect to 10.1.2.3`                    | `connect to {{IPV4:1}}`     | `single IPv4`      |                                               |


### 1.3 Multiple items and mixed content


| Scenario                                 | Given                                    | Expected                                            | Test in `TestMask`                           | Manual check                                               |
| ---------------------------------------- | ---------------------------------------- | --------------------------------------------------- | -------------------------------------------- | ---------------------------------------------------------- |
| User pastes two emails                   | `a@b.co and c@d.co`                      | `{{EMAIL:1}} and {{EMAIL:2}}`                       | `two emails`                                 |                                                            |
| User pastes same email twice             | `see shenali@ifs.com or shenali@ifs.com` | `see {{EMAIL:1}} or {{EMAIL:1}}` (same placeholder) | `same email twice gets same placeholder`     | `echo "see shenali@ifs.com or shenali@ifs.com" | ./maskit` |
| User pastes same URL twice               | `link https://x.com and https://x.com`   | `link {{URL:1}} and {{URL:1}}` (same placeholder)   | `same URL twice gets same placeholder`       |                                                            |
| User pastes same IPv4 twice              | `primary 10.1.2.3, backup 10.1.2.3`      | `primary {{IPV4:1}}, backup {{IPV4:1}}`             | `same IPv4 twice gets same placeholder`      |                                                            |
| User pastes two different IPv4 addresses | `a 10.1.2.3 and b 192.168.0.1`           | `a {{IPV4:1}} and b {{IPV4:2}}`                     | `different IPv4s get different placeholders` |                                                            |
| User pastes email and URL                | `mail a@b.co or https://x.com`           | `mail {{EMAIL:1}} or {{URL:1}}`                     | `email and URL`                              | `echo "mail a@b.co or https://x.com" | ./maskit`           |


### 1.4 Formatting and placeholders


| Scenario                                     | Given                        | Expected                        | Test in `TestMask`                      | Manual check                             |
| -------------------------------------------- | ---------------------------- | ------------------------------- | --------------------------------------- | ---------------------------------------- |
| User pastes text with surrounding whitespace | `a@b.co \n\t`                | `{{EMAIL:1}} \n\t`              | `preserve whitespace`                   | `echo " a@b.co " | ./maskit`             |
| Email appears inside a URL                   | `https://user@host.com/path` | `{{URL:1}}` (no separate email) | `email inside URL not matched as email` | `echo "https://u@h.com/path" | ./maskit` |


### 1.5 Localhost and loopback (safe URLs — not masked)


| Scenario                                              | Given                                           | Expected                                        | Test in `TestMask`                                     | Manual check                                                      |
| ----------------------------------------------------- | ----------------------------------------------- | ----------------------------------------------- | ------------------------------------------------------ | ----------------------------------------------------------------- |
| User pastes localhost with port                       | `see http://localhost:3030`                     | `see http://localhost:3030`                     | `localhost URL not masked`                             | `echo "see http://localhost:3030" | ./maskit`                     |
| User pastes localhost without port                    | `see http://localhost`                          | `see http://localhost`                          | `localhost without port not masked`                    |                                                                   |
| User pastes 127.0.0.1                                 | `see http://127.0.0.1:8080`                     | `see http://127.0.0.1:8080`                     | `127.0.0.1 not masked`                                 |                                                                   |
| User pastes IPv6 loopback                             | `see http://[::1]:8080`                         | `see http://[::1]:8080`                         | `IPv6 loopback not masked`                             |                                                                   |
| User pastes 0.0.0.0                                   | `see http://0.0.0.0:80`                         | `see http://0.0.0.0:80`                         | `0.0.0.0 not masked`                                   |                                                                   |
| User pastes IPv4 loopback without URL                 | `ping 127.0.0.1 then 0.0.0.0`                   | `ping 127.0.0.1 then 0.0.0.0`                   | `IPv4 loopback not masked`                             |                                                                   |
| User pastes localhost + real URL                      | `http://localhost:3030 and https://api.foo.com` | `http://localhost:3030 and {{URL:1}}`           | `localhost unchanged but real URL masked`              | `echo "http://localhost:3030 and https://api.foo.com" | ./maskit` |
| User pastes localhost.example.com (example domain)    | `see http://localhost.example.com`              | `see http://localhost.example.com` (not masked) | `localhost.example.com not masked (example domain)`    |                                                                   |
| User pastes non-example host with "localhost" in name | `see http://localhost.othercompany.com`         | `see {{URL:1}}` (masked)                        | `non-example host with localhost in name still masked` |                                                                   |


### 1.6 Example domains (safe — not masked)

RFC 2606 reserved example.com, example.org, example.net (and subdomains) for docs/samples; these are not masked to reduce noise.


| Scenario                                                     | Given                           | Expected                        | Test in `TestMask`                     | Manual check                                 |
| ------------------------------------------------------------ | ------------------------------- | ------------------------------- | -------------------------------------- | -------------------------------------------- |
| User pastes URL to example.com                               | `see https://example.com`       | `see https://example.com`       | `example.com URL not masked`           | `echo "see https://example.com" | ./maskit`  |
| User pastes URL to example.org with path                     | `see https://example.org/docs`  | `see https://example.org/docs`  | `example.org URL not masked`           |                                              |
| User pastes URL to [www.example.net](http://www.example.net) | `visit https://www.example.net` | `visit https://www.example.net` | `www.example.net URL not masked`       |                                              |
| User pastes email at example.com                             | `contact user@example.com`      | `contact user@example.com`      | `email at example.com not masked`      | `echo "contact user@example.com" | ./maskit` |
| User pastes email at example.org                             | `reply to dev@example.org`      | `reply to dev@example.org`      | `email at example.org not masked`      |                                              |
| User pastes email at subdomain of example.com                | `team@docs.example.com`         | `team@docs.example.com`         | `email at docs.example.com not masked` |                                              |


### 1.8 Input already contains placeholders (re-mask / paste)

When input already has `{{EMAIL:n}}` or `{{URL:n}}`, new sensitive data is numbered from max existing +1 so identities stay distinct.


| Scenario                                                 | Given                                                                  | Expected                                                 | Test in `TestMask`                                     | Manual check                                                            |
| -------------------------------------------------------- | ---------------------------------------------------------------------- | -------------------------------------------------------- | ------------------------------------------------------ | ----------------------------------------------------------------------- |
| Existing {{EMAIL:1}} and one real email                  | `see i have an email {{EMAIL:1}} and shenali@ifs.com`                  | `see i have an email {{EMAIL:1}} and {{EMAIL:2}}`        | `existing EMAIL:1 and real email`                      | `echo "see i have an email {{EMAIL:1}} and shenali@ifs.com" | ./maskit` |
| Existing {{URL:1}} and one real URL                      | `link {{URL:1}} and https://real.com`                                  | `link {{URL:1}} and {{URL:2}}`                           | `existing URL:1 and real URL`                          |                                                                         |
| Existing {{EMAIL:1}} {{EMAIL:2}} and one more real email | `{{EMAIL:1}} {{EMAIL:2}} new@real.com`                                 | `{{EMAIL:1}} {{EMAIL:2}} {{EMAIL:3}}`                    | `existing EMAIL:1 and EMAIL:2 and one more real email` |                                                                         |
| Existing placeholders and new email and URL              | `{{EMAIL:1}} and admin@corp.com see {{URL:1}} or https://api.corp.com` | `{{EMAIL:1}} and {{EMAIL:2}} see {{URL:1}} or {{URL:2}}` | `existing placeholders and new email and URL`          |                                                                         |


### 1.9 CLI usage


| Scenario                            | Given                      | Expected                 | Test in `TestMask` | Manual check                                                  |
| ----------------------------------- | -------------------------- | ------------------------ | ------------------ | ------------------------------------------------------------- |
| User runs maskit with stdin         | `echo "a@b.co" | ./maskit` | `{{EMAIL:1}}`            | (Mask behavior)    | Run from repo root after `go build -o maskit .`               |
| User runs maskit with file as input | `./maskit < file.txt`      | Masked content on stdout | —                  | `echo "secret@x.com" > /tmp/in.txt && ./maskit < /tmp/in.txt` |
| Binary builds from repo root        | `go build -o maskit .`     | Success, `./maskit` runs | —                  | `go build -o maskit .`                                        |


---

## Test plan changelog


| Date      | Change                                                                                                                                                                                                                                                                                             |
| --------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| (initial) | Test plan created. Restructured as scenario testing: no sensitive content, single/multiple/mixed content, formatting, localhost/loopback safe URLs, CLI. All 19 `TestMask` cases mapped to scenarios.                                                                                              |
| (update)  | Safe example domains (RFC 2606): example.com, example.org, example.net and subdomains—URLs and emails at these hosts are not masked. Added §1.6 Example domains; updated single-URL and localhost scenarios to use non-example URLs; added scenario for non-example host with "localhost" in name. |
| (update)  | Input already contains placeholders: new sensitive data is numbered from max existing +1 (re-mask / paste). Added §1.8 and four scenarios; renumbered CLI to §1.9.                                                                                                                                 |
| (update)  | Same email or same URL repeated: deduplicate by value—repeated content gets the same placeholder number (e.g. two [shenali@ifs.com](mailto:shenali@ifs.com) → {{EMAIL:1}} and {{EMAIL:1}}). Added scenarios in §1.3.                                                                               |
| (update)  | IPv4 addresses: added detection and masking with {{IPV4:n}}, including dedupe (same IPv4 → same placeholder), multiple IPv4s, IPv4 inside URLs (URL wins), and loopback IPv4 left unmasked. Added scenarios in §1.2, §1.3, and §1.5.                                                               |


---

