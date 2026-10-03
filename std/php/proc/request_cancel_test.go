package proc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php/core"
)

type requestTestContext struct {
	data.Context
	request context.Context
}

func (c *requestTestContext) GoContext() context.Context { return c.request }

func TestProcessCancellationChild(t *testing.T) {
	marker := os.Args[len(os.Args)-1]
	if marker == "origami-tree-parent" {
		child := exec.Command(os.Args[0], "-test.run=^TestProcessCancellationChild$", "--", "origami-tree-descendant")
		child.Stdout, child.Stderr = os.Stdout, os.Stderr
		if err := child.Start(); err != nil {
			t.Fatal(err)
		}
		fmt.Println(child.Process.Pid)
		time.Sleep(30 * time.Second)
		return
	}
	if marker != "origami-cancel-child" && marker != "origami-tree-descendant" {
		return
	}
	time.Sleep(30 * time.Second)
}

func TestCancelProcOpenClosesDescendantPipes(t *testing.T) {
	function := NewProcOpenFunction()
	base := runtime.NewContext(runtime.NewVM(parser.NewParser())).CreateContext(function.GetVariables())
	base.SetIndexZVal(0, data.NewZVal(data.NewArrayValue([]data.Value{
		data.NewStringValue(os.Args[0]), data.NewStringValue("-test.run=^TestProcessCancellationChild$"), data.NewStringValue("--"), data.NewStringValue("origami-tree-parent"),
	})))
	base.SetIndexZVal(1, data.NewZVal(data.NewNullValue()))
	base.SetIndexZVal(2, data.NewZVal(data.NewNullValue()))
	request, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	result, ctl := function.Call(&requestTestContext{Context: base, request: request})
	if ctl != nil {
		t.Fatal(ctl)
	}
	resource, ok := result.(*core.ResourceValue)
	if !ok {
		t.Fatalf("proc_open = %T", result)
	}
	process := resource.GetResource().(*ProcessInfo)
	defer process.Close()
	pipes, _ := base.GetIndexValue(2)
	stdout, _ := pipes.(*data.ArrayValue).FindSlotByIntKey(1)
	reader := stdout.ReadValue().(*core.ResourceValue).GetResource().(io.ReadCloser)
	defer reader.Close()
	buffered := bufio.NewReader(reader)
	line, err := buffered.ReadString('\n')
	pid, parseErr := strconv.Atoi(strings.TrimSpace(line))
	if err != nil || parseErr != nil {
		t.Fatalf("descendant startup: %q %v", line, err)
	}
	waitExit := watchDescendant(t, pid)
	readDone := make(chan struct{})
	go func() { _, _ = io.ReadAll(buffered); close(readDone) }()
	cancel()
	if !waitExit() {
		t.Fatal("descendant survived request cancellation")
	}
	select {
	case <-process.done:
	case <-time.After(3 * time.Second):
		t.Fatal("parent remained alive")
	}
	select {
	case <-readDone:
	case <-time.After(time.Second):
		t.Fatal("descendant kept inherited pipe open after cancellation")
	}
}

func TestCancelProcOpenReapsChildAndPipes(t *testing.T) {
	function := NewProcOpenFunction()
	base := runtime.NewContext(runtime.NewVM(parser.NewParser())).CreateContext(function.GetVariables())
	base.SetIndexZVal(0, data.NewZVal(data.NewArrayValue([]data.Value{
		data.NewStringValue(os.Args[0]), data.NewStringValue("-test.run=^TestProcessCancellationChild$"), data.NewStringValue("--"), data.NewStringValue("origami-cancel-child"),
	})))
	base.SetIndexZVal(1, data.NewZVal(data.NewNullValue()))
	base.SetIndexZVal(2, data.NewZVal(data.NewArrayValue(nil)))
	request, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := &requestTestContext{Context: base, request: request}
	result, control := function.Call(ctx)
	if control != nil {
		t.Fatal(control)
	}
	resource, ok := result.(*core.ResourceValue)
	if !ok {
		t.Fatalf("proc_open = %T", result)
	}
	process := resource.GetResource().(*ProcessInfo)
	defer process.Cmd.Cancel()
	defer process.Close()
	statusFunction := NewProcGetStatusFunction()
	statusContext := base.CreateContext(statusFunction.GetVariables())
	statusContext.SetIndexZVal(0, data.NewZVal(resource))
	statusValue, ctl := statusFunction.Call(statusContext)
	status, ok := statusValue.(*data.ArrayValue)
	if ctl != nil || !ok {
		t.Fatalf("proc_get_status returned %T, %v", statusValue, ctl)
	}
	running, exists := status.LookupZValByStringKey("running")
	if !exists || !running.ReadValue().(*data.BoolValue).Value || !process.GetRunning() {
		t.Fatal("proc_get_status marked a running child as exited")
	}
	pipes, _ := base.GetIndexValue(2)
	stdout, _ := pipes.(*data.ArrayValue).FindSlotByIntKey(1)
	reader := stdout.ReadValue().(*core.ResourceValue).GetResource().(io.ReadCloser)
	defer reader.Close()
	readDone := make(chan struct{})
	go func() { _, _ = io.ReadAll(reader); close(readDone) }()
	cancel()
	select {
	case <-process.done:
	case <-time.After(3 * time.Second):
		t.Fatal("canceled subprocess was not reaped")
	}
	select {
	case <-readDone:
	case <-time.After(time.Second):
		t.Fatal("canceled process pipe remained blocked")
	}
	if process.GetRunning() {
		t.Fatal("reaped process still reported running")
	}
}
