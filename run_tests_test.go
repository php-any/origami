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
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/mongodb"
	"github.com/php-any/origami/std/php/redis"
)

var suiteOptInExtensions = map[string]func(data.VM){
	"mongodb_bson_driver_test.php":        mongodb.Load,
	"redis_driver_contract_test.php":      redis.Load,
	"redis_driver_integration_test.php":   redis.Load,
	"redis_phpredis_features_test.php":    redis.Load,
}

var redisIntegrationScripts = map[string]bool{
	"redis_driver_integration_test.php": true,
	"redis_phpredis_features_test.php":  true,
}

func directScriptArgFromArgv() (string, bool) {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		if cmd.IsDirectScriptArg(arg) {
			return arg, true
		}
		break
	}
	return "", false
}

func runtimeLoaderForSuiteScript(scriptPath string) func(data.VM) {
	base := filepath.Base(scriptPath)
	optIn := suiteOptInExtensions[base]
	return func(vm data.VM) {
		loadDefaultRuntime(vm)
		if optIn != nil {
			optIn(vm)
		}
	}
}

func TestMain(m *testing.M) {
	// proc_open([PHP_BINARY, script.php]) uses the test binary as PHP_BINARY; run the script directly.
	if script, ok := directScriptArgFromArgv(); ok {
		if err := cmd.RunScriptFile(script); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// Each script is a separate CLI request. Sharing one VM changes declarations,
// handlers, argv and globals, and cannot test an uncaught-exception handler.
func TestPHPFixtureProcess(t *testing.T) {
	file := os.Getenv("ORIGAMI_SUITE_SCRIPT")
	if file == "" {
		return
	}
	base := filepath.Base(file)
	if redisIntegrationScripts[base] && os.Getenv("ORIGAMI_REDIS_TEST_PORT") == "" {
		t.Skip("set ORIGAMI_REDIS_TEST_PORT for real Redis integration")
	}
	cmd.SetRuntimeLoader(runtimeLoaderForSuiteScript(file))
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
			if redisIntegrationScripts[name] && os.Getenv("ORIGAMI_REDIS_TEST_PORT") == "" {
				t.Skip("set ORIGAMI_REDIS_TEST_PORT for real Redis integration")
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
