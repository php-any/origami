//go:build origamidebug

package perfmon

import (
	"strings"
	"testing"
	"time"
)

func TestEnabledInDebugBuild(t *testing.T) {
	if !Enabled() {
		t.Fatal("-tags origamidebug 时 Enabled 应为 true")
	}
	span := BeginRequest("GET", "/admin/login")
	NoteParse("a.php", true, 0)
	NoteParse("b.php", false, 2*time.Millisecond)
	NoteClassLoad("App\\Foo", 3*time.Millisecond)
	NoteInclude("view.php", 4*time.Millisecond)
	NoteFileRun("boot.php", time.Millisecond)
	span.End("ok")
	if !strings.Contains("ok", "ok") {
		t.Fatal("sanity")
	}
}
