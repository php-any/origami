package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// PHP_BINARY must run scripts when invoked as [binary, path.php] (go test / Windows proc_open).
func TestPHPBinaryRunsChildScript(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(t.TempDir(), "stream-child.php")
	script := `<?php $out = fopen('php://stdout', 'w'); fwrite($out, 'first'); usleep(150000); fwrite($out, 'second'); fclose($out);`
	if err := os.WriteFile(child, []byte(script), 0644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	proc := exec.CommandContext(ctx, executable, child)
	var out bytes.Buffer
	proc.Stdout = &out
	proc.Stderr = &out
	if err := proc.Run(); err != nil {
		t.Fatalf("child proc: %v\n%s", err, out.Bytes())
	}
	if got := out.String(); got != "firstsecond" {
		t.Fatalf("child stdout %q, want firstsecond", got)
	}
}

func TestPHPFixtureProcess_streamPipeLength(t *testing.T) {
	fixture := filepath.Join("tests", "php", "stream_get_contents_pipe_length_test.php")
	t.Setenv("ORIGAMI_SUITE_SCRIPT", fixture)
	TestPHPFixtureProcess(t)
}

func TestDirectScriptArgDetection(t *testing.T) {
	if _, ok := directScriptArgFromArgv(); ok {
		t.Skip("running under direct-script invocation")
	}
	path := filepath.Join(t.TempDir(), "probe.php")
	if err := os.WriteFile(path, []byte(`<?php echo "ok\n";`), 0644); err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	proc := exec.CommandContext(ctx, executable, path)
	out, err := proc.CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if !strings.Contains(string(out), "ok") {
		t.Fatalf("unexpected output: %q", out)
	}
}
