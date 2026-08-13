package phpt

import "testing"

func TestCookieToPhpAssignmentsKeepsTrailingSpaces(t *testing.T) {
	lines := cookieToPhpAssignments("cookie1=val1  ; cookie2=val2%20;cookie1")
	if len(lines) == 0 {
		t.Fatalf("no lines generated")
	}
	if lines[0] != "$_COOKIE['cookie1'] = 'val1  ';" {
		t.Fatalf("unexpected first line: %q", lines[0])
	}
}
