package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	netannotation "github.com/php-any/origami/std/net/annotation"
	netdata "github.com/php-any/origami/std/net/data"
)

type reloadingHandler struct {
	mu sync.RWMutex
	h  http.Handler
}

func (r *reloadingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	h := r.h
	r.mu.RUnlock()
	if h == nil {
		http.Error(w, "服务正在重载，请稍后重试", http.StatusServiceUnavailable)
		return
	}
	h.ServeHTTP(w, req)
}

func (r *reloadingHandler) set(h http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.h = h
}

type devServer struct {
	root    string
	port    int
	version atomic.Int64
	handler *reloadingHandler
}

func runDevServer(port int) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	ds := &devServer{
		root:    root,
		port:    port,
		handler: &reloadingHandler{},
	}

	if err := ds.reload(); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	watchDirs := []string{
		filepath.Join(root, "src"),
		filepath.Join(root, "lib"),
		filepath.Join(root, "cli"),
		filepath.Join(root, "config"),
		filepath.Join(root, "public"),
	}
	for _, dir := range watchDirs {
		if err := watchDirRecursive(watcher, dir); err != nil {
			log.Printf("警告: 无法监听 %s: %v", dir, err)
		}
	}
	for _, f := range []string{"index.php", "autoload.php"} {
		if err := watcher.Add(filepath.Join(root, f)); err != nil {
			log.Printf("警告: 无法监听 %s: %v", f, err)
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	srv := &http.Server{Addr: addr, Handler: ds.handler}

	go func() {
		var debounce *time.Timer
		var debounceC <-chan time.Time

		scheduleReload := func() {
			if debounce != nil {
				debounce.Stop()
			}
			debounce = time.NewTimer(300 * time.Millisecond)
			debounceC = debounce.C
		}

		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !isReloadableFile(ev.Name) {
					continue
				}
				if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
					continue
				}
				scheduleReload()
			case <-debounceC:
				debounceC = nil
				log.Println("[dev] 检测到文件变更，正在热重载...")
				if err := ds.reload(); err != nil {
					log.Printf("[dev] 热重载失败: %v", err)
				} else {
					log.Printf("[dev] 热重载完成 (version=%d)", ds.version.Load())
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("[dev] watcher 错误: %v", err)
			}
		}
	}()

	log.Printf("[dev] Agents Web 开发服务: http://localhost:%d", port)
	log.Println("[dev] 修改 src/、lib/、cli/、config/ 或 public/ 下的文件将自动热重载")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 服务错误: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func (ds *devServer) reload() error {
	resetHotReloadState()

	vm, _ := buildVM()
	version := ds.version.Add(1)

	os.Setenv("ORIGAMI_DEV", "1")
	os.Setenv("ORIGAMI_DEV_VERSION", strconv.FormatInt(version, 10))
	os.Setenv("ORIGAMI_DEV_PORT", strconv.Itoa(ds.port))

	if _, ctl := vm.LoadAndRun(filepath.Join(ds.root, "index.php")); ctl != nil {
		return fmt.Errorf("加载 index.php 失败: %v", ctl)
	}

	handler, err := extractHTTPHandler(vm)
	if err != nil {
		return err
	}

	ds.handler.set(handler)
	return nil
}

func extractHTTPHandler(vm *runtime.VM) (http.Handler, error) {
	val := vm.GlobalValue("__origami_dev_server")
	if val == nil {
		return nil, fmt.Errorf("未找到 __origami_dev_server，请确认 index.php 在 ORIGAMI_DEV 模式下运行")
	}
	var classVal *data.ClassValue
	switch v := val.(type) {
	case *data.ClassValue:
		classVal = v
	default:
		return nil, fmt.Errorf("__origami_dev_server 类型错误: %T", val)
	}

	if src, ok := classVal.Class.(data.GetSource); ok {
		if mux, ok := src.GetSource().(http.Handler); ok {
			return mux, nil
		}
	}

	return nil, fmt.Errorf("无法从 Server 对象提取 http.Handler")
}

func watchDirRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}

func isReloadableFile(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".php", ".zy", ".html", ".css", ".js":
		return true
	default:
		return false
	}
}

func resetHotReloadState() {
	runtime.ClearAutoLoad()
	node.ClearIncludeCache()
	netdata.ClearHTTPRoutes()
	netannotation.ResetHTTPApplicationState()
}
