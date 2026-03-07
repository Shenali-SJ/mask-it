# Manual testing guide (developer)

Step-by-step way to test MASKIT locally as the developer.

---

## 1. Build the binary

From the project root:

```bash
cd "/Users/shenali/Projects/The MaskIT Project"
go build -o maskit .
```

You should get a `maskit` binary in the current directory. (Using `-o maskit` avoids a name clash with a `maskit` directory if you had one.)

---

## 2. Run the unit tests

```bash
go test ./... -v
```

All tests should pass. This checks the `Mask()` logic (emails, http/https/ftp/ftps URLs, placeholders, whitespace).

---

## 3. Test the CLI with stdin

The CLI reads from **stdin** and prints masked output to **stdout**.

**Plain text (no masking):**

```bash
echo "hello world" | ./maskit
```

Expected: `hello world`

**Single email:**

```bash
echo "contact me at alice@company.com" | ./maskit
```

Expected: `contact me at {{EMAIL:1}}`

**Single URL (https):**

```bash
echo "see https://api.company.com" | ./maskit
```

Expected: `see {{URL:1}}`

**FTP URL:**

```bash
echo "fetch from ftp://files.company.com/pub" | ./maskit
```

Expected: `fetch from {{URL:1}}`

**Example domains not masked (docs/samples):**

```bash
echo "contact user@example.com or see https://example.org" | ./maskit
```

Expected: `contact user@example.com or see https://example.org` (unchanged)

**Mixed (emails + URLs):**

```bash
echo "mail alice@company.com or visit https://api.company.com and ftp://internal.server.local" | ./maskit
```

Expected: `mail {{EMAIL:1}} or visit {{URL:1}} and {{URL:2}}`

**Whitespace preserved:**

```bash
echo "  foo@bar.com  " | ./maskit
```

Expected: `  {{EMAIL:1}}  ` (spaces and newline unchanged)

---

## 4. Test with a file (simulate paste from file)

Put sample text in a file, then pipe it in:

```bash
echo "My email is dev@company.com. Docs: https://docs.company.com" > /tmp/sample.txt
./maskit < /tmp/sample.txt
```

Expected: `My email is {{EMAIL:1}}. Docs: {{URL:1}}`

---

## 5. Quick checklist

| Scenario              | Command (concept)              | What to check                    |
|-----------------------|--------------------------------|----------------------------------|
| No secrets            | `echo "hello" \| ./maskit`     | Output equals input              |
| One email             | `echo "a@b.co" \| ./maskit`    | `{{EMAIL:1}}`                    |
| Two emails            | `echo "a@b.co and c@d.co" \| ./maskit` | `{{EMAIL:1}}` and `{{EMAIL:2}}` |
| HTTP/HTTPS URL        | `echo "https://x.com" \| ./maskit`     | `{{URL:1}}`                      |
| FTP/FTPS URL          | `echo "ftp://host/path" \| ./maskit`   | `{{URL:1}}`                      |
| Email inside URL      | `echo "https://u@h.com/path" \| ./maskit` | One `{{URL:1}}`, no separate email |
| Example domain       | `echo "user@example.com" \| ./maskit`    | Output unchanged (example.com not masked) |
| Multi-line / spaces   | `echo "  a@b.co  " \| ./maskit`| Spaces and newlines unchanged    |

---

## 6. Safe-hosts allowlist (optional)

URLs whose host is in an allowlist file are **not** masked, so you can keep reference URLs (e.g. docs) visible for the LLM.

- **Default path:** `~/.config/maskit/safe-hosts.json`
- **Override:** set `MASKIT_SAFE_HOSTS_FILE` to the path of your JSON file (absolute path recommended).
- **Format:** JSON array of host strings, e.g. `["docs.python.org", "developer.mozilla.org"]`. Subdomains match (e.g. `python.org` allows `docs.python.org`).

Example using the example file:

```bash
export MASKIT_SAFE_HOSTS_FILE="$(pwd)/safe-hosts.example.json"
echo "see https://docs.python.org/3/ and https://api.private.com" | ./maskit
```

Expected: `see https://docs.python.org/3/ and {{URL:1}}` (docs URL unchanged, api URL masked).

---

## 7. If something fails

- Run `go test ./... -v` and fix any failing test first.
- Ensure you’re in the project root and using the binary you just built (`./maskit`).
- On Windows, use `.\maskit` and adjust `echo`/quoting as needed for your shell.
