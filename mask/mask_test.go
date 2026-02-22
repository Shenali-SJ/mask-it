package mask

import "testing"

func TestMask(t *testing.T) {
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
		{"single URL", "see https://api.foo.com", "see {{URL:1}}"},
		{"http URL", "visit http://foo.org/path", "visit {{URL:1}}"},
		{"ftp URL", "fetch from ftp://files.company.com/pub", "fetch from {{URL:1}}"},
		{"ftps URL", "secure ftps://secure.company.com/data", "secure {{URL:1}}"},
		{"email and URL", "mail a@b.co or https://x.com", "mail {{EMAIL:1}} or {{URL:1}}"},
		{"preserve whitespace", "  a@b.co  \n\t", "  {{EMAIL:1}}  \n\t"},
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
		{"empty string", "", ""},
		{"no matches", "no secrets here", "no secrets here"},
		// Input already contains placeholders: new sensitive data numbered from max+1
		{"existing EMAIL:1 and real email", "see i have an email {{EMAIL:1}} and shenali@ifs.com", "see i have an email {{EMAIL:1}} and {{EMAIL:2}}"},
		{"existing URL:1 and real URL", "link {{URL:1}} and https://real.com", "link {{URL:1}} and {{URL:2}}"},
		{"existing EMAIL:1 and EMAIL:2 and one more real email", "{{EMAIL:1}} {{EMAIL:2}} new@real.com", "{{EMAIL:1}} {{EMAIL:2}} {{EMAIL:3}}"},
		{"existing placeholders and new email and URL", "{{EMAIL:1}} and admin@corp.com see {{URL:1}} or https://api.corp.com", "{{EMAIL:1}} and {{EMAIL:2}} see {{URL:1}} or {{URL:2}}"},
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
