package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/php-any/origami/cmd"
)

// Each script is a separate CLI request. Sharing one VM changes declarations,
// handlers, argv and globals, and cannot test an uncaught-exception handler.
func TestPHPFixtureProcess(t *testing.T) {
	file := os.Getenv("ORIGAMI_SUITE_SCRIPT")
	if file == "" {
		return
	}
	os.Args = []string{os.Args[0], file}
	if err := cmd.RunScriptFile(file); err != nil {
		t.Fatal(err)
	}
}

func TestRunTests(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	files, err := filepath.Glob(filepath.Join("tests", "*", "*.php"))
	if err != nil {
		t.Fatal(err)
	}
	// Include fixtures execute under their owning regression's bound scope.
	fixtures := map[string]string{
		"bound_start_component_inner.php":  "bound_start_component_test.php",
		"include_bound_this_html_view.php": "include_bound_this_html_test.php",
		"include_bound_this_view.php":      "include_bound_this_test.php",
		"include_throw_target.php":         "include_throw_propagates_test.php",
	}
	unix := map[string]bool{"proc_open_array_cmd_test.php": true, "proc_open_reuse_pipes_test.php": true}
	debug := filepath.Join("examples", "laravel13", "storage", "origami-debug", "root-suite")
	if err := os.MkdirAll(debug, 0755); err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		t.Run(filepath.ToSlash(file), func(t *testing.T) {
			name := filepath.Base(file)
			if owner, fixture := fixtures[name]; fixture {
				t.Skip("include fixture covered by " + owner)
			}
			if runtime.GOOS == "windows" && unix[name] {
				t.Skip("requires Unix /bin/echo and stty; run suite on Linux")
			}
			request, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			defer cancel()
			process := exec.CommandContext(request, executable, "-test.run=^TestPHPFixtureProcess$")
			process.Env = append(os.Environ(), "ORIGAMI_SUITE_SCRIPT="+file)
			output, err := process.CombinedOutput()
			logName := strings.ReplaceAll(filepath.ToSlash(file), "/", "_") + ".log"
			if writeErr := os.WriteFile(filepath.Join(debug, logName), output, 0644); writeErr != nil {
				t.Error(writeErr)
			}
			// This regression intentionally ends in an uncaught Error, proving
			// that the following assignment and Log::fatal never execute.
			expectedUncaught := name == "throw_aborts_following_statements_test.php"
			if request.Err() != nil {
				t.Fatalf("script exceeded 90s: %s", output)
			}
			if strings.Contains(string(output), "[FATAL]") {
				t.Fatalf("script assertion failed: %s", output)
			}
			if expectedUncaught {
				if err == nil || !strings.Contains(string(output), "has") {
					t.Fatalf("expected uncaught null method Error: %v\n%s", err, output)
				}
			} else if err != nil {
				t.Fatalf("script failed: %v\n%s", err, output)
			}
		})
	}
}
