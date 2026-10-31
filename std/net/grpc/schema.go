package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bufbuild/protocompile"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/net/internal/script"
	"github.com/php-any/origami/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

type schema struct {
	files      *protoregistry.Files
	extensions *protoregistry.Types
	types      *dynamicpb.Types
	services   []protoreflect.ServiceDescriptor
}

func compile(ctx context.Context, files, paths []string) (*schema, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one proto file is required")
	}
	compiler := protocompile.Compiler{Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{ImportPaths: paths})}
	result, e := compiler.Compile(ctx, files...)
	if e != nil {
		return nil, e
	}
	s := &schema{files: new(protoregistry.Files), extensions: new(protoregistry.Types)}
	seen := map[string]bool{}
	var add func(protoreflect.FileDescriptor) error
	add = func(f protoreflect.FileDescriptor) error {
		if seen[f.Path()] {
			return nil
		}
		seen[f.Path()] = true
		for i := 0; i < f.Imports().Len(); i++ {
			if e := add(f.Imports().Get(i).FileDescriptor); e != nil {
				return e
			}
		}
		return s.files.RegisterFile(f)
	}
	for _, f := range result {
		if e := add(f); e != nil {
			return nil, e
		}
		for i := 0; i < f.Services().Len(); i++ {
			s.services = append(s.services, f.Services().Get(i))
		}
	}
	var register func(protoreflect.ExtensionDescriptors, protoreflect.MessageDescriptors)
	register = func(ext protoreflect.ExtensionDescriptors, msg protoreflect.MessageDescriptors) {
		for i := 0; i < ext.Len(); i++ {
			s.extensions.RegisterExtension(dynamicpb.NewExtensionType(ext.Get(i)))
		}
		for i := 0; i < msg.Len(); i++ {
			m := msg.Get(i)
			register(m.Extensions(), m.Messages())
		}
	}
	s.files.RangeFiles(func(f protoreflect.FileDescriptor) bool { register(f.Extensions(), f.Messages()); return true })
	s.types = dynamicpb.NewTypes(s.files)
	return s, nil
}
func newSchema() data.ClassStmt {
	s := &schema{}
	c := &script.Class{Name: "Net\\Grpc\\Schema", Factory: newSchema, Source: s}
	c.Methods = []data.Method{
		&script.Method{Name: "__construct", Args: []string{"files", "importPaths"}, Run: func(ctx data.Context) (data.GetValue, error) {
			files, e := utils.ConvertFromIndex[[]string](ctx, 0)
			if e != nil {
				return nil, e
			}
			paths, e := utils.ConvertFromIndex[[]string](ctx, 1)
			if e != nil {
				return nil, e
			}
			compiled, e := compile(ctx.GoContext(), files, paths)
			if e == nil {
				*s = *compiled
			}
			return nil, e
		}},
		&script.Method{Name: "methodsJson", Run: func(data.Context) (data.GetValue, error) {
			methods := []map[string]any{}
			for _, svc := range s.services {
				for i := 0; i < svc.Methods().Len(); i++ {
					m := svc.Methods().Get(i)
					methods = append(methods, map[string]any{"service": string(svc.FullName()), "method": string(m.Name()), "path": "/" + string(svc.FullName()) + "/" + string(m.Name()), "input": string(m.Input().FullName()), "output": string(m.Output().FullName()), "clientStreaming": m.IsStreamingClient(), "serverStreaming": m.IsStreamingServer()})
				}
			}
			b, e := json.Marshal(methods)
			return data.NewStringValue(string(b)), e
		}},
	}
	return c
}
