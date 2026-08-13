package php

import (
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewVersionCompareFunction 创建 version_compare
func NewVersionCompareFunction() data.FuncStmt {
	return &VersionCompareFunction{}
}

type VersionCompareFunction struct{}

func (f *VersionCompareFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v1 := ""
	v2 := ""
	if s, ok := ctx.GetIndexValue(0); ok && s != nil {
		if as, ok := s.(data.AsString); ok {
			v1 = as.AsString()
		}
	}
	if s, ok := ctx.GetIndexValue(1); ok && s != nil {
		if as, ok := s.(data.AsString); ok {
			v2 = as.AsString()
		}
	}
	cmp := comparePHPVersions(v1, v2)

	op := ""
	if s, ok := ctx.GetIndexValue(2); ok && s != nil {
		if as, ok := s.(data.AsString); ok {
			op = as.AsString()
		}
	}
	if op == "" {
		return data.NewIntValue(cmp), nil
	}
	return data.NewBoolValue(versionCompareOp(cmp, op)), nil
}

func versionCompareOp(cmp int, op string) bool {
	switch op {
	case "<", "lt":
		return cmp < 0
	case "<=", "le":
		return cmp <= 0
	case ">", "gt":
		return cmp > 0
	case ">=", "ge":
		return cmp >= 0
	case "==", "=", "eq":
		return cmp == 0
	case "!=", "<>", "ne":
		return cmp != 0
	default:
		return false
	}
}

func comparePHPVersions(a, b string) int {
	pa := splitPHPVersion(a)
	pb := splitPHPVersion(b)
	n := len(pa)
	if len(pb) > n {
		n = len(pb)
	}
	for i := 0; i < n; i++ {
		ai, asi := 0, ""
		bi, bsi := 0, ""
		if i < len(pa) {
			ai, asi = pa[i].num, pa[i].special
		}
		if i < len(pb) {
			bi, bsi = pb[i].num, pb[i].special
		}
		aw := phpVersionWeight(ai, asi)
		bw := phpVersionWeight(bi, bsi)
		if aw < bw {
			return -1
		}
		if aw > bw {
			return 1
		}
	}
	return 0
}

type phpVerPart struct {
	num     int
	special string
}

func splitPHPVersion(v string) []phpVerPart {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	parts := strings.FieldsFunc(v, func(r rune) bool {
		return r == '.' || r == '-' || r == '_' || r == '+'
	})
	out := make([]phpVerPart, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if n, err := strconv.Atoi(p); err == nil {
			out = append(out, phpVerPart{num: n})
			continue
		}
		// 7.3.0-dev → special "dev"
		out = append(out, phpVerPart{special: strings.ToLower(p)})
	}
	return out
}

func phpVersionWeight(num int, special string) int {
	// PHP: 任何特殊后缀低于同级正式版；dev < alpha < beta < RC < # < pl/stable
	if special == "" {
		return num*100 + 50
	}
	base := num * 100
	switch {
	case strings.HasPrefix(special, "dev"):
		return base + 0
	case strings.HasPrefix(special, "alpha") || special == "a":
		return base + 10
	case strings.HasPrefix(special, "beta") || special == "b":
		return base + 20
	case strings.HasPrefix(special, "rc"):
		return base + 30
	case special == "#" || strings.HasPrefix(special, "pl") || special == "p" || special == "stable":
		return base + 60
	default:
		return base + 40
	}
}

func (f *VersionCompareFunction) GetName() string { return "version_compare" }
func (f *VersionCompareFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "version1", 0, nil, data.String{}),
		node.NewParameter(nil, "version2", 1, nil, data.String{}),
		node.NewParameter(nil, "operator", 2, data.NewNullValue(), data.NewNullableType(data.String{})),
	}
}
func (f *VersionCompareFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "version1", 0, data.String{}),
		node.NewVariable(nil, "version2", 1, data.String{}),
		node.NewVariable(nil, "operator", 2, data.NewNullableType(data.String{})),
	}
}
