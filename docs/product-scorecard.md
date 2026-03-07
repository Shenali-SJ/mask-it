# MASKIT — Product scorecard & reality check

Track complexity, productizability, and other signals so you know when this is **market-ready** vs. **still a side project**. Update this file whenever you implement a feature or a meaningful part of one.

---
 
## Scorecard (current)

| Dimension | Score (1–5) | Notes |
|-----------|-------------|--------|
| **Productize / monetize** | 1 | No paid use case yet; single dev tool, no distribution. |
| **Complexity** | 1 | Small codebase, regex-based; easy to understand and copy. |
| **Market readiness** | 1 | MVP logic only; no packaging, docs, or install story. |
| **Differentiation** | 2 | Local-only + paste-before-AI is a clear angle; many “mask secrets” tools exist. |
| **Scope completeness (MVP)** | 2 | Core flow works (stdin → mask → stdout); email + URL only; file arg + more detectors pending. |
| **Defensibility / moat** | 1 | Easy to replicate in a weekend; no IP, network, or ecosystem lock-in. |
| **User friction** | 3 | Zero config, pipe-only; but manual paste/pipe is friction for “paste into AI” workflow. |

**Overall “can I market this?”** — **Not yet.** You have a working proof-of-concept. Need packaging, clear positioning, and at least one “why pay” hook (e.g. IDE integration, team policy, or compliance story).

---

## Scoring guide (1–5)

Use this when updating scores so ratings stay consistent.

| Dimension | 1 | 2 | 3 | 4 | 5 |
|-----------|---|---|---|---|---|
| **Productize / monetize** | No path to payment; feature, not product | Could be OSS + donate; no clear paid tier | Clear free/paid or B2B angle; no sales yet | Early revenue or committed pilots | Recurring revenue or clear monetization playbook |
| **Complexity** | Trivial (script-level) | Small, single-purpose (e.g. current MVP) | Multiple components, some design decisions | Rich behavior, config, or integrations | System product; many moving parts |
| **Market readiness** | Barely runs locally | Runs; no install story or docs | Installable + README; no support story | Docs, install, basic support or community | Launch-ready: positioning, support, updates |
| **Differentiation** | “Yet another X” | Clear twist (e.g. local-only, AI-paste) | Distinct workflow or audience | Recognizable category + unique angle | Category-defining or strong brand |
| **Scope completeness (MVP)** | Idea only | Core flow works; major gaps | MVP feature-complete per spec | MVP + one “v2” differentiator | Beyond MVP; roadmap in motion |
| **Defensibility / moat** | Copyable in days | Copyable in weeks | Needs some depth (e.g. detectors, quality) | Integrations or data/network effects | Strong IP, ecosystem, or switching cost |
| **User friction** | High (compile, hack config) | Medium (CLI works; manual steps) | Low (one command or install) | Very low (IDE/automation) | Friction-free in target workflow |

---

## Feature log (update on each change)

Log each feature or meaningful change: what shipped, and how it moved the scores (if at all).

| Date | What shipped | Productize | Complexity | Market readiness | Differentiation | Scope (MVP) | Defensibility | User friction |
|------|----------------|------------|------------|------------------|-----------------|-------------|---------------|----------------|
| (initial) | CLI stdin→stdout; `Mask()` with email + URL detection; deterministic placeholders; table-driven tests; preserve whitespace | 1 | 1 | 1 | 2 | 2 | 1 | 3 |
| (update) | Safe example domains (RFC 2606): example.com, example.org, example.net and subdomains—URLs and emails at these hosts are not masked; reduces noise in docs/samples | 1 | 1 | 1 | 2 | 2 | 1 | 3 |
| (update) | Re-mask / paste: when input already contains {{EMAIL:n}} or {{URL:n}}, new sensitive data is numbered from max existing +1 so identities stay distinct | 1 | 1 | 1 | 2 | 2 | 1 | 3 |
| (update) | Same email/URL repeated: deduplicate by value—repeated content gets the same placeholder (e.g. shenali@ifs.com twice → {{EMAIL:1}} and {{EMAIL:1}}) | 1 | 1 | 1 | 2 | 2 | 1 | 3 |
| (update) | IPv4 detection: mask IPv4 addresses as {{IPV4:n}} (dedupe by value, loopback IPv4 left unmasked, IPv4 inside URLs treated as part of URL) | 1 | 1 | 1 | 2 | 2 | 1 | 3 |
| (update) | DB connection strings: mask common connection-string URLs (postgres/mysql/mariadb/sqlserver/mssql/mongodb/redis) as {{CONN:n}} (dedupe by value, no separate masking of inner emails/URLs/IPs) | 1 | 1 | 1 | 3 | 3 | 1 | 3 |
| (update) | Safe-hosts allowlist: optional JSON file lists hosts whose URLs are not masked (reference URLs for LLM context); default path ~/.config/maskit/safe-hosts.json, override via MASKIT_SAFE_HOSTS_FILE | 1 | 1 | 1 | 3 | 3 | 1 | 3 |

---

## Reality-check questions

Ask these periodically. Answer in 1–2 sentences and add the date.

- **Who would pay, and for what?** (e.g. “Dev teams who need to paste logs into AI without leaking PII.”)
- **What’s the smallest “product” that could take money?** (e.g. “CLI + one IDE plugin + pay‑what‑you‑want.”)
- **What would make someone choose this over a 10-line script?** (e.g. “Better detectors, zero config, one binary.”)
- **What’s the next single thing that would most improve “can I market this?”** (e.g. “Single-binary release + README with install and one use case.”)

---

## How to use this file

1. **After implementing a feature:** add a row to the **Feature log** and adjust the **Scorecard** if any dimension changed.
2. **Before adding scope:** check **Scope completeness** and **Complexity** so you don’t overbuild before productizing.
3. **When considering “launch”:** ensure **Market readiness** and **Productize / monetize** are at least 3, and **Differentiation** is clear.
4. **Monthly (or per milestone):** answer the **Reality-check questions** and keep the “Overall” line honest.

---

*Last scorecard update: safe-hosts allowlist (optional JSON file) for reference URLs.* 
