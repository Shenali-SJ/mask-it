package mask

import "testing"

func setSafeHostsForTest(list []string) {
	safeHostsMu.Lock()
	defer safeHostsMu.Unlock()
	safeHostsOverride = list
	safeHostsLoaded = false
}

func TestMask(t *testing.T) {
	// Use empty allowlist so tests do not depend on user's safe-hosts file.
	setSafeHostsForTest([]string{})
	defer setSafeHostsForTest(nil)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain text unchanged", "hello world", "hello world"},
		{"single email", "contact me at a@b.co", "contact me at {{EMAIL:1}}"},
		{"two emails", "a@b.co and c@d.co", "{{EMAIL:1}} and {{EMAIL:2}}"},
		{"same email twice gets same placeholder", "see i have an email shenali@ifs.com or shenali@ifs.com", "see i have an email {{EMAIL:1}} or {{EMAIL:1}}"},
		{"same URL twice gets same placeholder", "link https://x.com and https://x.com", "link {{URL:1}} and {{URL:1}}"},
		{"single IPv4", "connect to 10.1.2.3", "connect to {{IPV4:1}}"},
		{"same IPv4 twice gets same placeholder", "primary 10.1.2.3, backup 10.1.2.3", "primary {{IPV4:1}}, backup {{IPV4:1}}"},
		{"different IPv4s get different placeholders", "a 10.1.2.3 and b 192.168.0.1", "a {{IPV4:1}} and b {{IPV4:2}}"},
		{"single URL", "see https://api.foo.com", "see {{URL:1}}"},
		{"http URL", "visit http://foo.org/path", "visit {{URL:1}}"},
		{"IPv4 inside URL not double-masked", "visit http://10.1.2.3/path", "visit {{URL:1}}"},
		{"ftp URL", "fetch from ftp://files.company.com/pub", "fetch from {{URL:1}}"},
		{"ftps URL", "secure ftps://secure.company.com/data", "secure {{URL:1}}"},
		{"email and URL", "mail a@b.co or https://x.com", "mail {{EMAIL:1}} or {{URL:1}}"},
		{"preserve whitespace", "  a@b.co  \n\t", "  {{EMAIL:1}}  \n\t"},
		{"IPv4 loopback not masked", "ping 127.0.0.1 then 0.0.0.0", "ping 127.0.0.1 then 0.0.0.0"},
		{"email inside URL not matched as email", "https://user@host.com/path", "{{URL:1}}"},
		{"localhost URL not masked", "see http://localhost:3030", "see http://localhost:3030"},
		{"localhost without port not masked", "see http://localhost", "see http://localhost"},
		{"127.0.0.1 not masked", "see http://127.0.0.1:8080", "see http://127.0.0.1:8080"},
		{"IPv6 loopback not masked", "see http://[::1]:8080", "see http://[::1]:8080"},
		{"0.0.0.0 not masked", "see http://0.0.0.0:80", "see http://0.0.0.0:80"},
		{"localhost unchanged but real URL masked", "http://localhost:3030 and https://api.foo.com", "http://localhost:3030 and {{URL:1}}"},
		{"localhost.example.com not masked (example domain)", "see http://localhost.example.com", "see http://localhost.example.com"},
		{"non-example host with localhost in name still masked", "see http://localhost.othercompany.com", "see {{URL:1}}"},
		{"example.com URL not masked", "see https://example.com", "see https://example.com"},
		{"example.org URL not masked", "see https://example.org/docs", "see https://example.org/docs"},
		{"www.example.net URL not masked", "visit https://www.example.net", "visit https://www.example.net"},
		{"email at example.com not masked", "contact user@example.com", "contact user@example.com"},
		{"email at example.org not masked", "reply to dev@example.org", "reply to dev@example.org"},
		{"email at docs.example.com not masked", "team@docs.example.com", "team@docs.example.com"},
		// Connection strings
		{"single postgres connection string", "postgres://user:pass@db.internal:5432/app", "{{CONN:1}}"},
		{"same connection string twice gets same placeholder", "a postgres://user:pass@db:5432/app and postgres://user:pass@db:5432/app", "a {{CONN:1}} and {{CONN:1}}"},
		{"different connection strings get different placeholders", "postgres://a:a@db1/app postgres://b:b@db2/app", "{{CONN:1}} {{CONN:2}}"},
		{"connection string and normal URL", "postgres://user:pass@db:5432/app and https://x.com", "{{CONN:1}} and {{URL:1}}"},
		{"connection string with email-like query is single placeholder", "postgres://user:pass@db.internal/app?email=user@example.com", "{{CONN:1}}"},
		{"empty string", "", ""},
		{"no matches", "no secrets here", "no secrets here"},
		// Input already contains placeholders: new sensitive data numbered from max+1
		{"existing EMAIL:1 and real email", "see i have an email {{EMAIL:1}} and shenali@ifs.com", "see i have an email {{EMAIL:1}} and {{EMAIL:2}}"},
		{"existing URL:1 and real URL", "link {{URL:1}} and https://real.com", "link {{URL:1}} and {{URL:2}}"},
		{"existing EMAIL:1 and EMAIL:2 and one more real email", "{{EMAIL:1}} {{EMAIL:2}} new@real.com", "{{EMAIL:1}} {{EMAIL:2}} {{EMAIL:3}}"},
		{"existing placeholders and new email and URL", "{{EMAIL:1}} and admin@corp.com see {{URL:1}} or https://api.corp.com", "{{EMAIL:1}} and {{EMAIL:2}} see {{URL:1}} or {{URL:2}}"},
		{"existing CONN:1 and new connection string", "see {{CONN:1}} and postgres://user:pass@db2/app2", "see {{CONN:1}} and {{CONN:2}}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mask(tt.input)
			if got != tt.want {
				t.Errorf("Mask(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMask_safeHostsAllowlist(t *testing.T) {
	// Allowlist makes URLs to these hosts (and subdomains) not masked.
	setSafeHostsForTest([]string{"docs.python.org", "developer.mozilla.org"})
	defer setSafeHostsForTest(nil)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"URL with allowlisted host not masked", "see https://docs.python.org/3/", "see https://docs.python.org/3/"},
		{"URL with allowlisted subdomain not masked", "see https://developer.mozilla.org/en-US/docs", "see https://developer.mozilla.org/en-US/docs"},
		{"URL with non-allowlisted host still masked", "see https://api.secret.com", "see {{URL:1}}"},
		{"allowlisted and non-allowlisted", "docs https://docs.python.org/3/ and api https://api.private.com", "docs https://docs.python.org/3/ and api {{URL:1}}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mask(tt.input)
			if got != tt.want {
				t.Errorf("Mask(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
