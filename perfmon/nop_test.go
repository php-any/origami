//go:build !origamidebug

package perfmon

import "testing"

func TestDisabledInReleaseBuild(t *testing.T) {
	if Enabled() {
		t.Fatal("无 origamidebug tag 时 Enabled 应为 false")
	}
	stop := Start()
	stop()
	span := BeginRequest("GET", "/")
	span.End("ok")
	NoteParse("x.php", false, 0)
	NoteClassLoad("Foo", 0)
	NoteInclude("x.php", 0)
	NoteFileRun("x.php", 0)
}
