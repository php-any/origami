//go:build origamidebug

package perfmon

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func Enabled() bool { return true }

func Now() time.Time { return time.Now() }

func Since(t time.Time) time.Duration {
	if t.IsZero() {
		return 0
	}
	return time.Since(t)
}

type slowOp struct {
	kind string
	name string
	d    time.Duration
}

type reqAcc struct {
	parseHit   atomic.Int64
	parseMiss  atomic.Int64
	classLoad  atomic.Int64
	includes   atomic.Int64
	fileRuns   atomic.Int64
	parseDur   atomic.Int64
	classDur   atomic.Int64
	includeDur atomic.Int64
	fileDur    atomic.Int64

	mu    sync.Mutex
	slows []slowOp
}

func (a *reqAcc) addSlow(kind, name string, d time.Duration) {
	if d < time.Millisecond {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.slows = append(a.slows, slowOp{kind: kind, name: name, d: d})
	if len(a.slows) > 64 {
		sort.Slice(a.slows, func(i, j int) bool { return a.slows[i].d > a.slows[j].d })
		a.slows = a.slows[:32]
	}
}

type Span struct {
	method, path string
	start        time.Time
	acc          *reqAcc
	ended        atomic.Bool
}

var (
	spanMu sync.Mutex
	active *reqAcc
)

func currentAcc() *reqAcc {
	spanMu.Lock()
	a := active
	spanMu.Unlock()
	return a
}

func BeginRequest(method, path string) *Span {
	acc := &reqAcc{}
	spanMu.Lock()
	active = acc
	spanMu.Unlock()
	return &Span{method: method, path: path, start: time.Now(), acc: acc}
}

func (s *Span) End(result string) {
	if s == nil || !s.ended.CompareAndSwap(false, true) {
		return
	}
	spanMu.Lock()
	if active == s.acc {
		active = nil
	}
	spanMu.Unlock()

	dur := time.Since(s.start)
	acc := s.acc
	if acc == nil {
		fmt.Fprintf(os.Stderr, "origami-perf %s %s result=%s dur=%s\n", s.method, s.path, result, dur.Round(time.Millisecond))
		return
	}

	fmt.Fprintf(os.Stderr,
		"origami-perf %s %s result=%s dur=%s parse=hit:%d miss:%d/%s class_load=%d/%s include=%d/%s file_run=%d/%s\n",
		s.method, s.path, result, dur.Round(time.Millisecond),
		acc.parseHit.Load(), acc.parseMiss.Load(), time.Duration(acc.parseDur.Load()).Round(time.Millisecond),
		acc.classLoad.Load(), time.Duration(acc.classDur.Load()).Round(time.Millisecond),
		acc.includes.Load(), time.Duration(acc.includeDur.Load()).Round(time.Millisecond),
		acc.fileRuns.Load(), time.Duration(acc.fileDur.Load()).Round(time.Millisecond),
	)

	acc.mu.Lock()
	slows := append([]slowOp(nil), acc.slows...)
	acc.mu.Unlock()
	sort.Slice(slows, func(i, j int) bool { return slows[i].d > slows[j].d })
	n := 12
	if len(slows) < n {
		n = len(slows)
	}
	for i := 0; i < n; i++ {
		op := slows[i]
		fmt.Fprintf(os.Stderr, "  slow %-7s %8s  %s\n", op.kind, op.d.Round(time.Millisecond), shortName(op.name))
	}
}

func shortName(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	if i := strings.LastIndex(name, "/vendor/"); i >= 0 {
		return name[i+1:]
	}
	if i := strings.LastIndex(name, "/storage/framework/views/"); i >= 0 {
		return "views/" + filepath.Base(name)
	}
	if len(name) > 120 {
		return "..." + name[len(name)-117:]
	}
	return name
}

func NoteParse(file string, cached bool, d time.Duration) {
	acc := currentAcc()
	if acc == nil {
		return
	}
	if cached {
		acc.parseHit.Add(1)
		return
	}
	acc.parseMiss.Add(1)
	acc.parseDur.Add(int64(d))
	acc.addSlow("parse", file, d)
}

func NoteClassLoad(name string, d time.Duration) {
	acc := currentAcc()
	if acc == nil {
		return
	}
	acc.classLoad.Add(1)
	acc.classDur.Add(int64(d))
	acc.addSlow("class", name, d)
}

func NoteInclude(file string, d time.Duration) {
	acc := currentAcc()
	if acc == nil {
		return
	}
	acc.includes.Add(1)
	acc.includeDur.Add(int64(d))
	acc.addSlow("include", file, d)
}

func NoteFileRun(file string, d time.Duration) {
	acc := currentAcc()
	if acc == nil {
		return
	}
	acc.fileRuns.Add(1)
	acc.fileDur.Add(int64(d))
	acc.addSlow("run", file, d)
}
