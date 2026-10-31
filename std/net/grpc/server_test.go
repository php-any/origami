package grpc

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	grpcsrc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	rv1 "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestDynamicUnaryWireReflectionAndLimits(t *testing.T) {
	dir := t.TempDir()
	e := os.WriteFile(filepath.Join(dir, "echo.proto"), []byte(`syntax="proto3";package generic;import "google/protobuf/wrappers.proto";service Echo{rpc Say(google.protobuf.StringValue)returns(google.protobuf.StringValue);}`), 0600)
	if e != nil {
		t.Fatal(e)
	}
	s, e := compile(context.Background(), []string{"echo.proto"}, []string{dir})
	if e != nil {
		t.Fatal(e)
	}
	o := map[string]int{}
	for k, v := range defaults {
		o[k] = v
	}
	o["maxRecvBytes"] = 128
	o["maxSendBytes"] = 1024
	o["requestTimeoutMs"] = 100
	o["maxCalls"] = 1
	entered := make(chan struct{}, 1)
	release := make(chan struct{})
	srv, e := makeServer(s, o, true, func(ctx context.Context, c *rpcCall) error {
		switch c.input {
		case `"bad"`:
			c.err = status.Error(codes.InvalidArgument, "bad input")
		case `"wait"`:
			entered <- struct{}{}
			select {
			case <-release:
			case <-ctx.Done():
				return ctx.Err()
			}
		case `"large"`:
			c.output = `"` + strings.Repeat("x", 2000) + `"`
		default:
			c.output = c.input
		}
		c.answered = true
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	defer srv.Stop()
	ln := bufconn.Listen(1 << 20)
	defer ln.Close()
	go srv.Serve(ln)
	cc, e := grpcsrc.NewClient("passthrough:///test", grpcsrc.WithTransportCredentials(insecure.NewCredentials()), grpcsrc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return ln.Dial() }))
	if e != nil {
		t.Fatal(e)
	}
	defer cc.Close()
	invoke := func(text string, options ...grpcsrc.CallOption) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply := new(wrapperspb.StringValue)
		e := cc.Invoke(ctx, "/generic.Echo/Say", wrapperspb.String(text), reply, options...)
		if e == nil && reply.Value != text {
			t.Errorf("reply %q != %q", reply.Value, text)
		}
		return e
	}
	if e := invoke("echo", grpcsrc.UseCompressor(gzip.Name)); e != nil {
		t.Fatal(e)
	}
	for _, test := range []struct {
		input string
		code  codes.Code
	}{{"bad", codes.InvalidArgument}, {strings.Repeat("x", 129), codes.ResourceExhausted}, {"large", codes.ResourceExhausted}, {"wait", codes.DeadlineExceeded}} {
		if e := invoke(test.input); status.Code(e) != test.code {
			t.Fatalf("%s: %v", test.input, e)
		}
	}
	<-entered
	done := make(chan error, 1)
	go func() { done <- invoke("wait") }()
	<-entered
	if e := invoke("overflow"); status.Code(e) != codes.ResourceExhausted {
		t.Fatalf("concurrency: %v", e)
	}
	close(release)
	<-done
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, e := rv1.NewServerReflectionClient(cc).ServerReflectionInfo(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if e = stream.Send(&rv1.ServerReflectionRequest{MessageRequest: &rv1.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: "generic.Echo"}}); e != nil {
		t.Fatal(e)
	}
	answer, e := stream.Recv()
	if e != nil || len(answer.GetFileDescriptorResponse().GetFileDescriptorProto()) != 2 {
		t.Fatalf("reflection: %v %v", answer, e)
	}
	// A client that leaves reflection open must not hold a slot indefinitely.
	_, e = stream.Recv()
	if status.Code(e) != codes.DeadlineExceeded {
		t.Fatalf("reflection deadline: %v", e)
	}
	if e = invoke("reclaimed"); e != nil {
		t.Fatalf("reflection slot not reclaimed: %v", e)
	}
}

func TestCanceledCallbackStillOccupiesExecutionSlot(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "echo.proto"), []byte(`syntax="proto3"; package generic; import "google/protobuf/wrappers.proto"; service Echo { rpc Say(google.protobuf.StringValue) returns (google.protobuf.StringValue); }`), 0600); err != nil {
		t.Fatal(err)
	}
	schema, err := compile(context.Background(), []string{"echo.proto"}, []string{dir})
	if err != nil {
		t.Fatal(err)
	}
	options := make(map[string]int)
	for key, value := range defaults {
		options[key] = value
	}
	options["maxCalls"] = 1
	options["requestTimeoutMs"] = 50
	entered := make(chan context.Context, 2)
	release := make(chan struct{})
	defer close(release)
	server, err := makeServer(schema, options, false, func(ctx context.Context, call *rpcCall) error {
		entered <- ctx
		<-release // Simulate a native extension that ignores cancellation.
		call.output, call.answered = call.input, true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	listener := bufconn.Listen(1 << 20)
	defer listener.Close()
	go server.Serve(listener)
	client, err := grpcsrc.NewClient("passthrough:///test", grpcsrc.WithTransportCredentials(insecure.NewCredentials()), grpcsrc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	invoke := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return client.Invoke(ctx, "/generic.Echo/Say", wrapperspb.String("test"), new(wrapperspb.StringValue))
	}
	done := make(chan error, 1)
	go func() { done <- invoke() }()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("callback did not start")
	}
	if err := <-done; status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("deadline: %v", err)
	}
	if err := invoke(); status.Code(err) != codes.ResourceExhausted {
		t.Fatalf("canceled callback escaped execution limit: %v", err)
	}
	select {
	case <-entered:
		t.Fatal("another callback started while the canceled one was still executing")
	default:
	}
}
