//go:build origamidebug

package perfmon

import (
	"fmt"
	"net/http"
	httppprof "net/http/pprof"
	"os"
	"runtime"
	goruntimepprof "runtime/pprof"
	"runtime/trace"
	"sync"
)

var startOnce sync.Once
var stopOnce sync.Once
var stopFn func()

// Start 按环境变量打开 CPU/堆/trace 采样和可选 HTTP pprof。返回的 stop 幂等，os.Exit 前必须显式调用。
func Start() func() {
	startOnce.Do(func() {
		var stops []func()
		fmt.Fprintln(os.Stderr, "origami-perf: origamidebug 已启用（生产构建不含此代码）")

		if path := os.Getenv("ORIGAMI_CPUPROFILE"); path != "" {
			f, err := os.Create(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "origami-perf: ORIGAMI_CPUPROFILE 无法创建 %s: %v\n", path, err)
			} else if err := goruntimepprof.StartCPUProfile(f); err != nil {
				fmt.Fprintf(os.Stderr, "origami-perf: StartCPUProfile: %v\n", err)
				_ = f.Close()
			} else {
				fmt.Fprintf(os.Stderr, "origami-perf: CPU profile -> %s\n", path)
				stops = append(stops, func() {
					goruntimepprof.StopCPUProfile()
					_ = f.Close()
				})
			}
		}

		if path := os.Getenv("ORIGAMI_TRACE"); path != "" {
			f, err := os.Create(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "origami-perf: ORIGAMI_TRACE 无法创建 %s: %v\n", path, err)
			} else if err := trace.Start(f); err != nil {
				fmt.Fprintf(os.Stderr, "origami-perf: trace.Start: %v\n", err)
				_ = f.Close()
			} else {
				fmt.Fprintf(os.Stderr, "origami-perf: execution trace -> %s\n", path)
				stops = append(stops, func() {
					trace.Stop()
					_ = f.Close()
				})
			}
		}

		if path := os.Getenv("ORIGAMI_MEMPROFILE"); path != "" {
			stops = append(stops, func() {
				f, err := os.Create(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "origami-perf: ORIGAMI_MEMPROFILE 无法创建 %s: %v\n", path, err)
					return
				}
				runtime.GC()
				if err := goruntimepprof.WriteHeapProfile(f); err != nil {
					fmt.Fprintf(os.Stderr, "origami-perf: WriteHeapProfile: %v\n", err)
				}
				_ = f.Close()
				fmt.Fprintf(os.Stderr, "origami-perf: heap profile -> %s\n", path)
			})
		}

		if addr := os.Getenv("ORIGAMI_PPROF_ADDR"); addr != "" {
			mux := http.NewServeMux()
			mux.HandleFunc("/debug/pprof/", httppprof.Index)
			mux.HandleFunc("/debug/pprof/cmdline", httppprof.Cmdline)
			mux.HandleFunc("/debug/pprof/profile", httppprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", httppprof.Symbol)
			mux.HandleFunc("/debug/pprof/trace", httppprof.Trace)
			go func() {
				fmt.Fprintf(os.Stderr, "origami-perf: pprof HTTP %s/debug/pprof/\n", addr)
				if err := http.ListenAndServe(addr, mux); err != nil {
					fmt.Fprintf(os.Stderr, "origami-perf: pprof ListenAndServe: %v\n", err)
				}
			}()
		}

		stopFn = func() {
			stopOnce.Do(func() {
				for i := len(stops) - 1; i >= 0; i-- {
					stops[i]()
				}
			})
		}
	})
	if stopFn == nil {
		return func() {}
	}
	return stopFn
}
