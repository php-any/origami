package php

import (
	"net"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// InetPtonFunction 实现 inet_pton：IP 文本 → 二进制，失败返回 false。
type InetPtonFunction struct{}

func NewInetPtonFunction() data.FuncStmt { return &InetPtonFunction{} }

func (f *InetPtonFunction) GetName() string { return "inet_pton" }

func (f *InetPtonFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ip_address", 0, nil, data.NewBaseType("string")),
	}
}

func (f *InetPtonFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ip_address", 0, data.NewBaseType("string")),
	}
}

func (f *InetPtonFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	s := ""
	if v != nil {
		s = v.AsString()
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return data.NewBoolValue(false), nil
	}
	if v4 := ip.To4(); v4 != nil {
		return data.NewStringValue(string(v4)), nil
	}
	v16 := ip.To16()
	if v16 == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(string(v16)), nil
}

// InetNtopFunction 实现 inet_ntop：二进制 → IP 文本，失败返回 false。
type InetNtopFunction struct{}

func NewInetNtopFunction() data.FuncStmt { return &InetNtopFunction{} }

func (f *InetNtopFunction) GetName() string { return "inet_ntop" }

func (f *InetNtopFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "in_addr", 0, nil, data.NewBaseType("string")),
	}
}

func (f *InetNtopFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "in_addr", 0, data.NewBaseType("string")),
	}
}

func (f *InetNtopFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	if _, isBool := v.(*data.BoolValue); isBool {
		return data.NewBoolValue(false), nil
	}
	packed := []byte(v.AsString())
	if len(packed) != net.IPv4len && len(packed) != net.IPv6len {
		return data.NewBoolValue(false), nil
	}
	s := net.IP(packed).String()
	if s == "<nil>" || s == "" {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(s), nil
}
