package http_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/php"
)

func TestControllerRequestParseLifecycle(t *testing.T) {
	dir := t.TempDir()
	files := []string{"bootstrap.php", "routes.php", "startup.php", "src/Controller/GreetingController.php", "src/Service/GreetingService.php", "src/View/GreetingView.php", "views/greeting.php"}
	for _, file := range files {
		content, err := os.ReadFile(filepath.Join("../../../examples/http-hot-reload", file))
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, file)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, content, 0600); err != nil {
			t.Fatal(err)
		}
	}
	entry := filepath.Join(dir, "test.php")
	if err := os.WriteFile(entry, []byte(`<?php require __DIR__.'/bootstrap.php'; $server = new Net\Http\Server('127.0.0.1', port: 18087); require __DIR__.'/routes.php'; return $server;`), 0600); err != nil {
		t.Fatal(err)
	}
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	httplib.Load(base)
	value, ctl := base.LoadAndRun(entry)
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	server := httptest.NewServer(value.(*data.ClassValue).GetSource().(http.Handler))
	defer server.Close()
	client := &http.Client{Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	check := func(path, want string, contains bool) {
		t.Helper()
		response, err := client.Get(server.URL + path)
		if err != nil {
			t.Error(err)
			return
		}
		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			t.Error(err)
			return
		}
		matches := string(body) == want
		if contains {
			matches = strings.Contains(string(body), want)
		}
		if response.StatusCode != 200 || !matches {
			t.Errorf("%s status=%d body=%q want=%q", path, response.StatusCode, body, want)
		}
	}
	check("/api/greeting?name=one", `{"controller":"controller-v1","message":"service-v1: Hello, one"}`, false)
	check("/", "<h1>view-v1</h1>", true)
	for _, file := range []string{files[3], files[4], files[6], "startup.php"} {
		path := filepath.Join(dir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(strings.ReplaceAll(string(content), "v1", "v2")), 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, info.ModTime(), info.ModTime()); err != nil {
			t.Fatal(err)
		}
	}
	var group sync.WaitGroup
	for i := 0; i < 12; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			check(fmt.Sprintf("/api/greeting?name=%d", i), fmt.Sprintf(`{"controller":"controller-v2","message":"service-v2: Hello, %d"}`, i), false)
		}(i)
	}
	group.Wait()
	check("/", "<h1>view-v2</h1>", true)
	check("/startup", "v1", false)
	if _, ok := base.GetClass("HotReload\\Controller\\GreetingController"); ok {
		t.Fatal("controller leaked to Base")
	}
	if _, ok := base.GetParser().GetClassPathManager().FindClassFile("HotReload\\Controller\\GreetingController"); ok {
		t.Fatal("request class path leaked to Base")
	}
}
