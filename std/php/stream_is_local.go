package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/core"
	"strings"
)

type StreamIsLocalFunction struct{ IsFileFunction }

func (*StreamIsLocalFunction) GetName() string            { return "stream_is_local" }
func (*StreamIsLocalFunction) GetParams() []data.GetValue { return filepermsFunctionGetParams }
func (*StreamIsLocalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if resource, ok := v.(*core.ResourceValue); ok {
		if remote, ok := resource.GetResource().(interface{ IsRemote() bool }); ok {
			return data.NewBoolValue(!remote.IsRemote()), nil
		}
		return data.NewBoolValue(resource.IsPHPResourceOpen()), nil
	}
	path := strings.ToLower(v.AsString())
	for _, prefix := range []string{"http://", "https://", "ftp://", "ftps://", "data:"} {
		if strings.HasPrefix(path, prefix) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}
