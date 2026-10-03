package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/php-any/origami/data"

	"github.com/php-any/origami/std/laravel/serve"
)

func TestRuntimeHTTPLifecycle(t *testing.T) {
	vm, _ := buildVM()
	signals := make(chan string, 4)
	if ctl := vm.RegisterFunction("runtimeShutdownSignal", func(s string) { signals <- s }); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if ctl := vm.RegisterFunction("runtimePanicSignal", func() { panic("lifecycle fixture panic") }); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	result, ctl := vm.LoadAndRun("tests/origami/runtime_lifecycle.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	app, ok := result.(*data.ClassValue)
	if !ok {
		t.Fatalf("fixture returned %T", result)
	}

	handler, err := serve.NewHTTPHandler(vm, app)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	client := &http.Client{Timeout: 35 * time.Second}
	check := func(method, path, want string, status int) {
		t.Helper()
		req, err := http.NewRequest(method, server.URL+path, nil)
		if err != nil {
			t.Error(err)
			return
		}
		response, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if response.StatusCode != status || string(body) != want {
			t.Errorf("%s %s: status=%d body=%q want=%d %q", method, path, response.StatusCode, body, status, want)
		}
	}
	check("GET", "/__runtime/order", "global:one>group>route>body<route<group<global|group-term|route-term|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/short", "short:one|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/custom-send", "custom-send|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/custom-throw", "partial-custom|shutdown:one", 200)
	check("GET", "/__runtime/stream", "first|second|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/chunks", "ab|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/exit", "exit-body|shutdown:one", 200)
	check("GET", "/__runtime/empty", "", 204)
	check("GET", "/__runtime/not-modified", "", 304)
	check("HEAD", "/__runtime/order", "", 200)
	check("HEAD", "/__runtime/stream", "", 200)
	check("GET", "/__runtime/throw", "handled-exception|global-term:one|shutdown:one", 500)
	check("GET", "/__runtime/stream-throw", "partial-stream|shutdown:one", 200)
	{
		response, err := client.Get(server.URL + "/__runtime/panic")
		if err != nil {
			t.Fatal(err)
		}
		panicBody, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != 500 || !strings.Contains(string(panicBody), "<title>Laravel</title>") || !strings.HasSuffix(string(panicBody), "|global-term:one|shutdown:one") {
			t.Fatalf("official Kernel panic handling: status=%d bytes=%d", response.StatusCode, len(panicBody))
		}
	}
	check("GET", "/__runtime/scoped", "1|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/scoped", "1|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/events", "1|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/events", "1|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/livewire-call", "magic-action|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/livewire-call", "magic-action|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/graph", "1|global-term:one|shutdown:one", 200)
	check("GET", "/__runtime/graph", "1|global-term:one|shutdown:one", 200)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			check("GET", "/__runtime/livewire-call", "magic-action|global-term:one|shutdown:one", 200)
			id := fmt.Sprint(i)
			check("GET", "/__runtime/graph?id="+id, "1|global-term:"+id+"|shutdown:"+id, 200)
			check("GET", "/__runtime/order?id="+id, "global:"+id+">group>route>body<route<group<global|group-term|route-term|global-term:"+id+"|shutdown:"+id, 200)
		}(i)
	}
	wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	request, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/__runtime/cancel", nil)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		response, _ := client.Do(request)
		if response != nil {
			response.Body.Close()
		}
	}()
	select {
	case s := <-signals:
		if s != "started" {
			t.Fatalf("cancel start=%q", s)
		}
	case <-time.After(5 * time.Second):
		cancel()
		t.Fatal("cancel route did not start")
	}
	cancel()
	select {
	case s := <-signals:
		if s != "shutdown" {
			t.Fatalf("cancel shutdown=%q", s)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("shutdown did not run after client cancellation")
	}
	<-done
	response, err := client.Get(server.URL + "/__runtime/file")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil || !strings.HasPrefix(string(body), "<?php") {
		t.Fatalf("file response: %v %q", err, body)
	}
	for _, method := range []string{"GET", "HEAD"} {
		req, _ := http.NewRequest(method, server.URL+"/__runtime/file", nil)
		req.Header.Set("Range", "bytes=0-3")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if method == "GET" && (resp.StatusCode != 206 || string(body) != "<?ph") {
			t.Fatalf("range: %d %q", resp.StatusCode, body)
		}
		if method == "HEAD" && (len(body) != 0 || resp.ContentLength <= 0) {
			t.Fatalf("HEAD file: len=%d content-length=%d", len(body), resp.ContentLength)
		}
	}
}
