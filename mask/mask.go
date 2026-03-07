package mask

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// urlSchemes lists protocol schemes matched as URLs. Add entries here to
// support more protocols; the URL regex is built from this list.
var urlSchemes = []string{"https?", "ftp", "ftps"}

func buildURLPattern() string {
	return "(" + strings.Join(urlSchemes, "|") + ")://[^\\s]+"
}

// urlRx matches URLs with any of urlSchemes (e.g. http, https, ftp, ftps).
var urlRx = regexp.MustCompile(buildURLPattern())

// connRx matches DB-style connection strings such as
// postgres://user:pass@host:5432/db or mysql://user:pass@host/db.
// It looks for a known scheme followed by :// and either an @ (user:pass@)
// or "password=" within the authority/query.
var connRx = regexp.MustCompile(`(?i)\b(postgres|postgresql|mysql|mariadb|sqlserver|mssql|mongodb|redis)://[^\s]*(?:@|password=)[^\s]*`)

// ipv4Rx matches IPv4 addresses like 192.168.0.1. It is intentionally simple
// and may match some invalid octets; this is acceptable for MVP masking.
var ipv4Rx = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

// emailRx matches typical email addresses.
var emailRx = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

// safeHostsPath returns the path to the safe-hosts allowlist file.
// Uses MASKIT_SAFE_HOSTS_FILE if set, otherwise ~/.config/maskit/safe-hosts.json.
func safeHostsPath() string {
	if p := os.Getenv("MASKIT_SAFE_HOSTS_FILE"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "maskit", "safe-hosts.json")
}

// loadSafeHostsFile reads the allowlist from the configured path.
// File format: JSON array of host strings, e.g. ["docs.python.org", "developer.mozilla.org"].
// Returns nil on any error or missing file (no extra safe hosts).
func loadSafeHostsFile() []string {
	path := safeHostsPath()
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return nil
	}
	return list
}

var (
	safeHostsMu       sync.Mutex
	safeHostsCached   []string
	safeHostsLoaded   bool
	safeHostsOverride []string
)

// getSafeHostsAllowlist returns the list of hosts that should not be masked as URLs.
// If safeHostsOverride is set, that list is returned; otherwise the file at
// safeHostsPath() is read once and cached.
func getSafeHostsAllowlist() []string {
	safeHostsMu.Lock()
	defer safeHostsMu.Unlock()
	if safeHostsOverride != nil {
		return safeHostsOverride
	}
	if !safeHostsLoaded {
		safeHostsCached = loadSafeHostsFile()
		safeHostsLoaded = true
	}
	return safeHostsCached
}

// isHostInAllowlist returns true if host equals an entry or is a subdomain of one
// (e.g. entry "python.org" matches "docs.python.org").
func isHostInAllowlist(host string, list []string) bool {
	for _, entry := range list {
		if host == entry || (entry != "" && strings.HasSuffix(host, "."+entry)) {
			return true
		}
	}
	return false
}

// isSafeURL reports whether the URL is a localhost/loopback address, an
// example-domain URL (RFC 2606), or a host in the safe-hosts allowlist.
func isSafeURL(url string) bool {
	i := strings.Index(url, "://")
	if i < 0 {
		return false
	}
	hostPart := url[i+3:]
	host := extractHostFromURLHostPart(hostPart)
	return isLoopbackHost(hostPart) || isExampleDomain(host) || isHostInAllowlist(host, getSafeHostsAllowlist())
}

// hostFromURLHostPart extracts the host from "host:port/path" or "host/path".
func extractHostFromURLHostPart(hostPart string) string {
	segment := hostPart
	if i := strings.Index(hostPart, "/"); i >= 0 {
		segment = hostPart[:i]
	}
	if j := strings.Index(segment, ":"); j >= 0 && !strings.HasPrefix(segment, "[") {
		segment = segment[:j]
	}
	return segment
}

// isExampleDomain reports whether the host is an RFC 2606 reserved example domain.
// RFC 2606 reserves only these second-level names for documentation: example.com,
// example.net, example.org (and their subdomains, e.g. foo.example.com). The RFC
// also reserves TLDs .test, .example, .invalid, .localhost—we do not treat those
// here. "foo" and "bar" are not in the RFC; they are conventional placeholders,
// but domains like foo.com can be real, so we do not treat them as safe.
func isExampleDomain(host string) bool {
	return host == "example.com" || host == "example.org" || host == "example.net" ||
		strings.HasSuffix(host, ".example.com") || strings.HasSuffix(host, ".example.org") || strings.HasSuffix(host, ".example.net")
}

// isSafeEmail reports whether the email is at an example domain (should not be masked).
func isSafeEmail(email string) bool {
	at := strings.LastIndex(email, "@")
	if at < 0 {
		return false
	}
	return isExampleDomain(email[at+1:])
}

func isLoopbackHost(hostPart string) bool {
	if hostPart == "localhost" || strings.HasPrefix(hostPart, "localhost:") || strings.HasPrefix(hostPart, "localhost/") {
		return true
	}
	if hostPart == "127.0.0.1" || strings.HasPrefix(hostPart, "127.0.0.1:") || strings.HasPrefix(hostPart, "127.0.0.1/") {
		return true
	}
	if hostPart == "[::1]" || strings.HasPrefix(hostPart, "[::1]:") || strings.HasPrefix(hostPart, "[::1]/") {
		return true
	}
	if hostPart == "0.0.0.0" || strings.HasPrefix(hostPart, "0.0.0.0:") || strings.HasPrefix(hostPart, "0.0.0.0/") {
		return true
	}
	return false
}

// existingPlaceholderRx finds {{EMAIL:n}}, {{URL:n}}, {{IPV4:n}}, or
// {{CONN:n}} in already-masked input.
var existingPlaceholderRx = regexp.MustCompile(`\{\{(EMAIL|URL|IPV4|CONN):(\d+)\}\}`)

// maxExistingPlaceholderNumbers returns the highest n already present in input
// for each type, so new placeholders can be numbered from max+1 (avoids
// reusing the same counter when input already contains placeholders, e.g.
// re-mask or paste).
func maxExistingPlaceholderNumbers(input string) map[string]int {
	max := map[string]int{
		"EMAIL": 0,
		"URL":   0,
		"IPV4":  0,
		"CONN":  0,
	}
	for _, sub := range existingPlaceholderRx.FindAllStringSubmatch(input, -1) {
		kind, numStr := sub[1], sub[2]
		n, _ := strconv.Atoi(numStr)
		if n > max[kind] {
			max[kind] = n
		}
	}
	return max
}

func insideAnySpan(start, end int, spans [][2]int) bool {
	for _, span := range spans {
		if start >= span[0] && end <= span[1] {
			return true
		}
	}
	return false
}

func placeholderFor(kind string, num int) string {
	switch kind {
	case "URL":
		return "{{URL:" + strconv.Itoa(num) + "}}"
	case "EMAIL":
		return "{{EMAIL:" + strconv.Itoa(num) + "}}"
	case "IPV4":
		return "{{IPV4:" + strconv.Itoa(num) + "}}"
	case "CONN":
		return "{{CONN:" + strconv.Itoa(num) + "}}"
	default:
		return ""
	}
}

// Mask returns input with emails, URLs (http, https, ftp, ftps), and IPv4
// addresses replaced by deterministic placeholders {{EMAIL:n}}, {{URL:n}},
// and {{IPV4:n}}. If the input already contains such placeholders, new ones
// are numbered from max existing +1 so identities stay distinct. Original
// formatting and whitespace are preserved.
func Mask(input string) string {
	max := maxExistingPlaceholderNumbers(input)
	maxEmail, maxURL, maxIPV4, maxConn := max["EMAIL"], max["URL"], max["IPV4"], max["CONN"]

	type match struct {
		start, end int
		kind       string
		num        int
	}
	var matches []match

	// Find connection strings first so we do not separately treat their inner
	// parts (emails, URLs, IPs) as independent secrets.
	connContentToNum := make(map[string]int)
	connNum := maxConn
	connSpans := make([][2]int, 0)
	for _, loc := range connRx.FindAllStringIndex(input, -1) {
		connStr := input[loc[0]:loc[1]]
		num, seen := connContentToNum[connStr]
		if !seen {
			connNum++
			num = connNum
			connContentToNum[connStr] = num
		}
		matches = append(matches, match{loc[0], loc[1], "CONN", num})
		connSpans = append(connSpans, [2]int{loc[0], loc[1]})
	}

	// Find URLs so we don't treat content inside URLs as emails.
	// Skip localhost/loopback URLs so they are not masked.
	// Same URL content gets the same placeholder number (deduplicate by value).
	urlContentToNum := make(map[string]int)
	urlNum := maxURL
	for _, loc := range urlRx.FindAllStringIndex(input, -1) {
		urlStr := input[loc[0]:loc[1]]
		if insideAnySpan(loc[0], loc[1], connSpans) {
			continue
		}
		if isSafeURL(urlStr) {
			continue
		}
		num, seen := urlContentToNum[urlStr]
		if !seen {
			urlNum++
			num = urlNum
			urlContentToNum[urlStr] = num
		}
		matches = append(matches, match{loc[0], loc[1], "URL", num})
	}

	// Find emails that are not inside a URL/CONN span.
	// Same email content gets the same placeholder number (deduplicate by value).
	urlSpans := make([][2]int, 0, len(matches))
	for _, m := range matches {
		if m.kind == "URL" || m.kind == "CONN" {
			urlSpans = append(urlSpans, [2]int{m.start, m.end})
		}
	}
	emailContentToNum := make(map[string]int)
	emailNum := maxEmail
	for _, loc := range emailRx.FindAllStringIndex(input, -1) {
		if insideAnySpan(loc[0], loc[1], urlSpans) {
			continue
		}
		emailStr := input[loc[0]:loc[1]]
		if isSafeEmail(emailStr) {
			continue
		}
		num, seen := emailContentToNum[emailStr]
		if !seen {
			emailNum++
			num = emailNum
			emailContentToNum[emailStr] = num
		}
		matches = append(matches, match{loc[0], loc[1], "EMAIL", num})
	}

	// Find IPv4 addresses that are not inside a URL span.
	// Same IPv4 content gets the same placeholder number (deduplicate by value).
	ipv4ContentToNum := make(map[string]int)
	ipv4Num := maxIPV4
	for _, loc := range ipv4Rx.FindAllStringIndex(input, -1) {
		if insideAnySpan(loc[0], loc[1], urlSpans) {
			continue
		}
		ipStr := input[loc[0]:loc[1]]
		// Treat loopback/localhost-style IPv4 as safe (not masked).
		if isLoopbackHost(ipStr) {
			continue
		}
		num, seen := ipv4ContentToNum[ipStr]
		if !seen {
			ipv4Num++
			num = ipv4Num
			ipv4ContentToNum[ipStr] = num
		}
		matches = append(matches, match{loc[0], loc[1], "IPV4", num})
	}

	// Replace from end to start so indices stay valid.
	sort.Slice(matches, func(i, j int) bool { return matches[i].start > matches[j].start })
	result := input
	for _, m := range matches {
		placeholder := placeholderFor(m.kind, m.num)
		result = result[:m.start] + placeholder + result[m.end:]
	}
	return result
}
