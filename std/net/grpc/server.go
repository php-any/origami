// Package grpc exposes general unary gRPC transport and dynamic protobuf schemas.
// Application routing, validation and status policy remain in script callbacks.
package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/net/internal/script"
	"github.com/php-any/origami/utils"
	grpcsrc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	rv1 "google.golang.org/grpc/reflection/grpc_reflection_v1"
	ralpha "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func Load(vm data.VM) {
	vm.AddClass(newSchema())
	vm.AddClass(newServer())
	vm.AddClass(callClass(nil))
}

var defaults = map[string]int{"maxConnections": 128, "maxConcurrentStreams": 32, "maxCalls": 128, "maxRecvBytes": 1048576, "maxSendBytes": 4194304, "maxHeaderBytes": 16384, "requestTimeoutMs": 10000, "connectionTimeoutMs": 5000, "idleTimeoutMs": 30000, "maxConnectionAgeMs": 300000}

func newServer() data.ClassStmt {
	address := ""
	var sch *schema
	opts := defaults
	reflectEnabled := false
	c := &script.Class{Name: "Net\\Grpc\\Server", Factory: newServer}
	c.Methods = []data.Method{
		&script.Method{Name: "__construct", Args: []string{"address", "schema"}, Run: func(ctx data.Context) (data.GetValue, error) {
			var e error
			address, e = utils.ConvertFromIndex[string](ctx, 0)
			if e != nil {
				return nil, e
			}
			sch, e = utils.ConvertFromIndex[*schema](ctx, 1)
			return nil, e
		}},
		&script.Method{Name: "configure", Args: []string{"options"}, Run: func(ctx data.Context) (data.GetValue, error) {
			v, e := script.Options(ctx, 0, opts)
			if e == nil {
				opts = v
			}
			return nil, e
		}},
		&script.Method{Name: "reflection", Args: []string{"enabled"}, Run: func(ctx data.Context) (data.GetValue, error) {
			var e error
			reflectEnabled, e = utils.ConvertFromIndex[bool](ctx, 0)
			return nil, e
		}},
		&script.Method{Name: "serve", Args: []string{"callback"}, Run: func(ctx data.Context) (data.GetValue, error) {
			fn, e := script.Callback(ctx, 0)
			if e != nil {
				return nil, e
			}
			if sch == nil || sch.files == nil {
				return nil, fmt.Errorf("schema not initialized")
			}
			srv, e := makeServer(sch, opts, reflectEnabled, func(gc context.Context, call *rpcCall) error {
				return script.Invoke(ctx.GetVM(), fn, gc, func(c data.Context) []data.Value { return []data.Value{data.NewProxyValue(callClass(call), c)} })
			})
			if e != nil {
				return nil, e
			}
			ln, e := net.Listen("tcp", address)
			if e != nil {
				srv.Stop()
				return nil, e
			}
			defer ln.Close()
			defer srv.Stop()
			stop := context.AfterFunc(ctx.GoContext(), srv.Stop)
			defer stop()
			return nil, srv.Serve(script.LimitListener(ln, opts["maxConnections"]))
		}},
	}
	return c
}

type rpcCall struct {
	method   string
	input    string
	output   string
	err      error
	answered bool
	ctx      context.Context
}

func callClass(call *rpcCall) data.ClassStmt {
	c := &script.Class{Name: "Net\\Grpc\\Call", Source: call}
	valid := func() error {
		if call == nil {
			return fmt.Errorf("uninitialized RPC")
		}
		if call.answered {
			return fmt.Errorf("RPC already completed")
		}
		return call.ctx.Err()
	}
	c.Methods = []data.Method{
		&script.Method{Name: "method", Run: func(data.Context) (data.GetValue, error) {
			if call == nil {
				return nil, fmt.Errorf("uninitialized RPC")
			}
			return data.NewStringValue(call.method), nil
		}},
		&script.Method{Name: "requestJson", Run: func(data.Context) (data.GetValue, error) {
			if call == nil {
				return nil, fmt.Errorf("uninitialized RPC")
			}
			return data.NewStringValue(call.input), nil
		}},
		&script.Method{Name: "respondJson", Args: []string{"json"}, Run: func(ctx data.Context) (data.GetValue, error) {
			if e := valid(); e != nil {
				return nil, e
			}
			v, e := utils.ConvertFromIndex[string](ctx, 0)
			if e == nil {
				call.output = v
				call.answered = true
			}
			return nil, e
		}},
		&script.Method{Name: "fail", Args: []string{"code", "message"}, Run: func(ctx data.Context) (data.GetValue, error) {
			if e := valid(); e != nil {
				return nil, e
			}
			code, e := utils.ConvertFromIndex[int](ctx, 0)
			if e != nil {
				return nil, e
			}
			if code < 1 || code > 16 {
				return nil, fmt.Errorf("gRPC error code must be 1..16")
			}
			message, e := utils.ConvertFromIndex[string](ctx, 1)
			if e == nil {
				call.err = status.Error(codes.Code(code), message)
				call.answered = true
			}
			return nil, e
		}},
	}
	return c
}
func makeServer(s *schema, o map[string]int, reflectionEnabled bool, callback func(context.Context, *rpcCall) error) (*grpcsrc.Server, error) {
	seen := map[string]bool{}
	for _, svc := range s.services {
		name := string(svc.FullName())
		if seen[name] || (reflectionEnabled && (name == "grpc.reflection.v1.ServerReflection" || name == "grpc.reflection.v1alpha.ServerReflection")) {
			return nil, fmt.Errorf("duplicate gRPC service: %s", name)
		}
		seen[name] = true
	}
	slots := make(chan struct{}, o["maxCalls"])
	// Cancellation releases transport admission, but cannot forcibly stop a
	// non-cooperative native callback. Keep its execution slot until it returns.
	executions := make(chan struct{}, o["maxCalls"])
	admission := func(ctx context.Context, _ *tap.Info) (context.Context, error) {
		select {
		case slots <- struct{}{}:
		default:
			return nil, status.Error(codes.ResourceExhausted, "concurrent call limit")
		}
		ctx, cancel := context.WithTimeout(ctx, time.Duration(o["requestTimeoutMs"])*time.Millisecond)
		var once sync.Once
		release := func() { once.Do(func() { <-slots; cancel() }) }
		context.AfterFunc(ctx, release)
		return context.WithValue(ctx, callSlotKey{}, release), nil
	}
	srv := grpcsrc.NewServer(grpcsrc.InTapHandle(admission), grpcsrc.MaxRecvMsgSize(o["maxRecvBytes"]), grpcsrc.MaxSendMsgSize(o["maxSendBytes"]), grpcsrc.MaxHeaderListSize(uint32(o["maxHeaderBytes"])), grpcsrc.MaxConcurrentStreams(uint32(o["maxConcurrentStreams"])), grpcsrc.ConnectionTimeout(time.Duration(o["connectionTimeoutMs"])*time.Millisecond), grpcsrc.StatsHandler(callStats{}), grpcsrc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: time.Duration(o["idleTimeoutMs"]) * time.Millisecond, MaxConnectionAge: time.Duration(o["maxConnectionAgeMs"]) * time.Millisecond, MaxConnectionAgeGrace: time.Duration(o["requestTimeoutMs"]) * time.Millisecond}))
	for _, service := range s.services {
		desc := grpcsrc.ServiceDesc{ServiceName: string(service.FullName()), HandlerType: (*interface{})(nil), Metadata: service.ParentFile().Path()}
		for i := 0; i < service.Methods().Len(); i++ {
			m := service.Methods().Get(i)
			if m.IsStreamingClient() || m.IsStreamingServer() {
				srv.Stop()
				return nil, fmt.Errorf("streaming application RPC is unsupported: %s", m.FullName())
			}
			path := "/" + string(service.FullName()) + "/" + string(m.Name())
			desc.Methods = append(desc.Methods, grpcsrc.MethodDesc{MethodName: string(m.Name()), Handler: func(_ any, ctx context.Context, decode func(any) error, interceptor grpcsrc.UnaryServerInterceptor) (any, error) {
				select {
				case executions <- struct{}{}:
					defer func() { <-executions }()
				default:
					return nil, status.Error(codes.ResourceExhausted, "active callback limit")
				}
				input := dynamicpb.NewMessage(m.Input())
				if e := decode(input); e != nil {
					return nil, e
				}
				handler := func(ctx context.Context, _ any) (any, error) { return dispatch(ctx, s, m, path, input, o, callback) }
				if interceptor == nil {
					return handler(ctx, input)
				}
				return interceptor(ctx, input, &grpcsrc.UnaryServerInfo{FullMethod: path}, handler)
			}})
		}
		srv.RegisterService(&desc, struct{}{})
	}
	if reflectionEnabled {
		ro := reflection.ServerOptions{Services: srv, DescriptorResolver: s.files, ExtensionResolver: s.extensions}
		rv1.RegisterServerReflectionServer(srv, reflection.NewServerV1(ro))
		ralpha.RegisterServerReflectionServer(srv, reflection.NewServer(ro))
	}
	return srv, nil
}
func dispatch(ctx context.Context, s *schema, m protoreflect.MethodDescriptor, path string, input *dynamicpb.Message, o map[string]int, callback func(context.Context, *rpcCall) error) (any, error) {
	b, e := (protojson.MarshalOptions{Resolver: s.types, UseProtoNames: true}).Marshal(input)
	if e != nil {
		return nil, status.Error(codes.Internal, "request serialization failed")
	}
	call := &rpcCall{method: path, input: string(b), ctx: ctx}
	e = callback(ctx, call)
	if ctx.Err() != nil {
		return nil, status.FromContextError(ctx.Err()).Err()
	}
	if e != nil {
		return nil, status.Error(codes.Internal, "script callback failed")
	}
	if call.err != nil {
		return nil, call.err
	}
	if !call.answered {
		return nil, status.Error(codes.Internal, "callback did not respond")
	}
	if int64(len(call.output)) > int64(o["maxSendBytes"])*6 {
		return nil, status.Error(codes.ResourceExhausted, "response JSON too large")
	}
	output := dynamicpb.NewMessage(m.Output())
	if e = (protojson.UnmarshalOptions{Resolver: s.types}).Unmarshal([]byte(call.output), output); e != nil {
		return nil, status.Error(codes.Internal, "invalid response JSON")
	}
	if proto.Size(output) > o["maxSendBytes"] {
		return nil, status.Error(codes.ResourceExhausted, "response too large")
	}
	return output, nil
}

// Release admission only after the entire RPC (including response transmission).
// The context cancellation hook also covers streams rejected before dispatch.
type callSlotKey struct{}
type callStats struct{}

func (callStats) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context { return ctx }
func (callStats) HandleRPC(ctx context.Context, event stats.RPCStats) {
	if _, ok := event.(*stats.End); ok {
		if release, ok := ctx.Value(callSlotKey{}).(func()); ok {
			release()
		}
	}
}
func (callStats) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context { return ctx }
func (callStats) HandleConn(context.Context, stats.ConnStats)                       {}
