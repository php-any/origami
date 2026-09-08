package php

import (
	"net"
	"strings"

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

// Ip2LongFunction 实现 ip2long：IPv4 点分地址 → 无符号 32 位整数。
// 对齐 64 位 PHP 8：成功返回 0..4294967295 的 int，非法地址返回 false；不接受 IPv6。
type Ip2LongFunction struct{}

func NewIp2LongFunction() data.FuncStmt { return &Ip2LongFunction{} }

func (f *Ip2LongFunction) GetName() string { return "ip2long" }

func (f *Ip2LongFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ip_address", 0, nil, data.NewBaseType("string")),
	}
}

func (f *Ip2LongFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ip_address", 0, data.NewBaseType("string")),
	}
}

func (f *Ip2LongFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	s := ""
	if v != nil {
		s = v.AsString()
	}
	if s == "" || strings.Contains(s, ":") {
		return data.NewBoolValue(false), nil
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return data.NewBoolValue(false), nil
	}
	v4 := ip.To4()
	if v4 == nil {
		return data.NewBoolValue(false), nil
	}
	n := uint32(v4[0])<<24 | uint32(v4[1])<<16 | uint32(v4[2])<<8 | uint32(v4[3])
	return data.NewIntValue(int(n)), nil
}

// Long2IpFunction 实现 long2ip：32 位整数 → IPv4 点分地址（按无符号 32 位截断）。
type Long2IpFunction struct{}

func NewLong2IpFunction() data.FuncStmt { return &Long2IpFunction{} }

func (f *Long2IpFunction) GetName() string { return "long2ip" }

func (f *Long2IpFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ip", 0, nil, data.NewBaseType("int")),
	}
}

func (f *Long2IpFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ip", 0, data.NewBaseType("int")),
	}
}

func (f *Long2IpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	n := uint32(0)
	if v != nil {
		if ai, ok := v.(data.AsInt); ok {
			iv, err := ai.AsInt()
			if err != nil {
				return data.NewStringValue("0.0.0.0"), nil
			}
			n = uint32(iv)
		}
	}
	return data.NewStringValue(net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n)).String()), nil
}
