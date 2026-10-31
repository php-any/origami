package redis

import (
	"fmt"
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
)

var scanParams = []data.GetValue{
	node.NewParameterReference(nil, "iterator", 0, nil, nil),
	node.NewParameter(nil, "pattern", 1, data.NewNullValue(), nil),
	node.NewParameter(nil, "count", 2, data.NewIntValue(0), nil),
	node.NewParameter(nil, "type", 3, data.NewNullValue(), nil),
}
var scanVariables = []data.Variable{
	node.NewVariable(nil, "iterator", 0, nil), node.NewVariable(nil, "pattern", 1, nil),
	node.NewVariable(nil, "count", 2, nil), node.NewVariable(nil, "type", 3, nil),
}

func scan(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
	if s.mode != atomicMode {
		return bad("SCAN requires atomic mode")
	}
	cursor := 0
	if !null(a[0]) {
		var err error
		cursor, err = integer(a[0])
		if err != nil || cursor < 0 {
			return bad("SCAN iterator must be null or a non-negative integer")
		}
		if cursor == 0 {
			return data.NewBoolValue(false), nil
		}
	}
	count, err := integer(a[2])
	if err != nil || count < 0 {
		return bad("SCAN count must be a non-negative integer")
	}
	for {
		args := []string{"SCAN", strconv.Itoa(cursor)}
		if !null(a[1]) {
			pattern, err := scalar(a[1])
			if err != nil {
				return failure(err)
			}
			if s.scanPrefix {
				pattern = s.prefix + pattern
			}
			args = append(args, "MATCH", pattern)
		}
		if count > 0 {
			args = append(args, "COUNT", strconv.Itoa(count))
		}
		if !null(a[3]) {
			kind, err := scalar(a[3])
			if err != nil {
				return failure(err)
			}
			args = append(args, "TYPE", kind)
		}
		if err := s.ensureConnected(runtime.RequestContext()); err != nil {
			return failure(err)
		}
		v, err := s.exchange(runtime.RequestContext(), args)
		if err != nil {
			return failure(err)
		}
		if _, ok := v.(serverError); ok {
			value, err := s.decode(v, false)
			if err != nil {
				return failure(err)
			}
			return value, nil
		}
		reply, ok := v.([]any)
		if !ok || len(reply) != 2 {
			return bad("invalid SCAN response")
		}
		str, ok := reply[0].(string)
		if !ok {
			return bad("invalid SCAN cursor")
		}
		cursor, err = strconv.Atoi(str)
		if err != nil || cursor < 0 {
			return failure(fmt.Errorf("invalid SCAN cursor %q", str))
		}
		keys, ok := reply[1].([]any)
		if !ok {
			return bad("invalid SCAN keys")
		}
		if ctl := ctx.SetVariableValue(scanVariables[0], data.NewIntValue(cursor)); ctl != nil {
			return nil, ctl
		}
		if len(keys) != 0 || cursor == 0 || !s.scanRetry {
			value, err := s.decode(keys, false)
			if err != nil {
				return failure(err)
			}
			return value, nil
		}
	}
}
