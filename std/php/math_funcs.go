package php

import (
	"math"
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ---------- 基础数学函数 ----------

// SqrtFunction 实现 sqrt 函数
type SqrtFunction struct{}

func NewSqrtFunction() data.FuncStmt { return &SqrtFunction{} }

func (f *SqrtFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Sqrt(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Sqrt(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *SqrtFunction) GetName() string { return "sqrt" }
var sqrtFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *SqrtFunction) GetParams() []data.GetValue {
	return sqrtFunctionGetParams
}
var sqrtFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *SqrtFunction) GetVariables() []data.Variable {
	return sqrtFunctionGetVariables
}

// CbrtFunction 实现 cbrt 函数
type CbrtFunction struct{}

func NewCbrtFunction() data.FuncStmt { return &CbrtFunction{} }

func (f *CbrtFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Cbrt(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Cbrt(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *CbrtFunction) GetName() string { return "cbrt" }
var cbrtFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *CbrtFunction) GetParams() []data.GetValue {
	return cbrtFunctionGetParams
}
var cbrtFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *CbrtFunction) GetVariables() []data.Variable {
	return cbrtFunctionGetVariables
}

// ExpFunction 实现 exp 函数
type ExpFunction struct{}

func NewExpFunction() data.FuncStmt { return &ExpFunction{} }

func (f *ExpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Exp(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Exp(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *ExpFunction) GetName() string { return "exp" }
var expFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *ExpFunction) GetParams() []data.GetValue {
	return expFunctionGetParams
}
var expFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *ExpFunction) GetVariables() []data.Variable {
	return expFunctionGetVariables
}

// LogFunction 实现 log 函数 (log(num, base=natural))
type LogFunction struct{}

func NewLogFunction() data.FuncStmt { return &LogFunction{} }

func (f *LogFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	baseV, _ := ctx.GetIndexValue(1)

	if v == nil {
		return data.NewFloatValue(0), nil
	}
	var num, base float64
	if af, ok := v.(data.AsFloat); ok {
		num, _ = af.AsFloat()
	} else if ai, ok := v.(*data.IntValue); ok {
		num = float64(ai.Value)
	}
	if baseV != nil {
		if af, ok := baseV.(data.AsFloat); ok {
			base, _ = af.AsFloat()
		} else if ai, ok := baseV.(*data.IntValue); ok {
			base = float64(ai.Value)
		}
	}
	if base == 0 || base == 1 {
		return data.NewFloatValue(math.Log(num)), nil
	}
	return data.NewFloatValue(math.Log(num) / math.Log(base)), nil
}

func (f *LogFunction) GetName() string { return "log" }
var logFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "num", 0, nil, nil),
	node.NewParameter(nil, "base", 1, node.NewIntLiteral(nil, "0"), nil),
}

func (f *LogFunction) GetParams() []data.GetValue {
	return logFunctionGetParams
}
var logFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "num", 0, nil),
	node.NewVariable(nil, "base", 1, data.NewBaseType("float")),
}

func (f *LogFunction) GetVariables() []data.Variable {
	return logFunctionGetVariables
}

// Log10Function 实现 log10 函数
type Log10Function struct{}

func NewLog10Function() data.FuncStmt { return &Log10Function{} }

func (f *Log10Function) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Log10(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Log10(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *Log10Function) GetName() string { return "log10" }
var log10FunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *Log10Function) GetParams() []data.GetValue {
	return log10FunctionGetParams
}
var log10FunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *Log10Function) GetVariables() []data.Variable {
	return log10FunctionGetVariables
}

// Log1pFunction 实现 log1p 函数
type Log1pFunction struct{}

func NewLog1pFunction() data.FuncStmt { return &Log1pFunction{} }

func (f *Log1pFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Log1p(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Log1p(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *Log1pFunction) GetName() string { return "log1p" }
var log1pFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *Log1pFunction) GetParams() []data.GetValue {
	return log1pFunctionGetParams
}
var log1pFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *Log1pFunction) GetVariables() []data.Variable {
	return log1pFunctionGetVariables
}

// Expm1Function 实现 expm1 函数
type Expm1Function struct{}

func NewExpm1Function() data.FuncStmt { return &Expm1Function{} }

func (f *Expm1Function) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Expm1(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Expm1(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *Expm1Function) GetName() string { return "expm1" }
var expm1FunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *Expm1Function) GetParams() []data.GetValue {
	return expm1FunctionGetParams
}
var expm1FunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *Expm1Function) GetVariables() []data.Variable {
	return expm1FunctionGetVariables
}

// PiFunction 实现 pi 函数
type PiFunction struct{}

func NewPiFunction() data.FuncStmt { return &PiFunction{} }

func (f *PiFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewFloatValue(math.Pi), nil
}

func (f *PiFunction) GetName() string { return "pi" }
var piFunctionGetParams = []data.GetValue{}

func (f *PiFunction) GetParams() []data.GetValue {
	return piFunctionGetParams
}
var piFunctionGetVariables = []data.Variable{}

func (f *PiFunction) GetVariables() []data.Variable {
	return piFunctionGetVariables
}

// ---------- 三角函数 ----------

// SinFunction 实现 sin 函数
type SinFunction struct{}

func NewSinFunction() data.FuncStmt { return &SinFunction{} }

func (f *SinFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Sin(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Sin(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *SinFunction) GetName() string { return "sin" }
var sinFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *SinFunction) GetParams() []data.GetValue {
	return sinFunctionGetParams
}
var sinFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *SinFunction) GetVariables() []data.Variable {
	return sinFunctionGetVariables
}

// CosFunction 实现 cos 函数
type CosFunction struct{}

func NewCosFunction() data.FuncStmt { return &CosFunction{} }

func (f *CosFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Cos(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Cos(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *CosFunction) GetName() string { return "cos" }
var cosFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *CosFunction) GetParams() []data.GetValue {
	return cosFunctionGetParams
}
var cosFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *CosFunction) GetVariables() []data.Variable {
	return cosFunctionGetVariables
}

// TanFunction 实现 tan 函数
type TanFunction struct{}

func NewTanFunction() data.FuncStmt { return &TanFunction{} }

func (f *TanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Tan(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Tan(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *TanFunction) GetName() string { return "tan" }
var tanFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *TanFunction) GetParams() []data.GetValue {
	return tanFunctionGetParams
}
var tanFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *TanFunction) GetVariables() []data.Variable {
	return tanFunctionGetVariables
}

// AcosFunction 实现 acos 函数
type AcosFunction struct{}

func NewAcosFunction() data.FuncStmt { return &AcosFunction{} }

func (f *AcosFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Acos(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Acos(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *AcosFunction) GetName() string { return "acos" }
var acosFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *AcosFunction) GetParams() []data.GetValue {
	return acosFunctionGetParams
}
var acosFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *AcosFunction) GetVariables() []data.Variable {
	return acosFunctionGetVariables
}

// AsinFunction 实现 asin 函数
type AsinFunction struct{}

func NewAsinFunction() data.FuncStmt { return &AsinFunction{} }

func (f *AsinFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Asin(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Asin(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *AsinFunction) GetName() string { return "asin" }
var asinFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *AsinFunction) GetParams() []data.GetValue {
	return asinFunctionGetParams
}
var asinFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *AsinFunction) GetVariables() []data.Variable {
	return asinFunctionGetVariables
}

// AtanFunction 实现 atan 函数
type AtanFunction struct{}

func NewAtanFunction() data.FuncStmt { return &AtanFunction{} }

func (f *AtanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Atan(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Atan(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *AtanFunction) GetName() string { return "atan" }
var atanFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *AtanFunction) GetParams() []data.GetValue {
	return atanFunctionGetParams
}
var atanFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *AtanFunction) GetVariables() []data.Variable {
	return atanFunctionGetVariables
}

// Atan2Function 实现 atan2 函数
type Atan2Function struct{}

func NewAtan2Function() data.FuncStmt { return &Atan2Function{} }

func (f *Atan2Function) Call(ctx data.Context) (data.GetValue, data.Control) {
	yV, _ := ctx.GetIndexValue(0)
	xV, _ := ctx.GetIndexValue(1)
	var y, x float64
	if af, ok := yV.(data.AsFloat); ok {
		y, _ = af.AsFloat()
	} else if ai, ok := yV.(*data.IntValue); ok {
		y = float64(ai.Value)
	}
	if af, ok := xV.(data.AsFloat); ok {
		x, _ = af.AsFloat()
	} else if ai, ok := xV.(*data.IntValue); ok {
		x = float64(ai.Value)
	}
	return data.NewFloatValue(math.Atan2(y, x)), nil
}

func (f *Atan2Function) GetName() string { return "atan2" }
var atan2FunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "y", 0, nil, nil),
	node.NewParameter(nil, "x", 1, nil, nil),
}

func (f *Atan2Function) GetParams() []data.GetValue {
	return atan2FunctionGetParams
}
var atan2FunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "y", 0, nil),
	node.NewVariable(nil, "x", 1, nil),
}

func (f *Atan2Function) GetVariables() []data.Variable {
	return atan2FunctionGetVariables
}

// ---------- 双曲函数 ----------

// CoshFunction 实现 cosh 函数
type CoshFunction struct{}

func NewCoshFunction() data.FuncStmt { return &CoshFunction{} }

func (f *CoshFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Cosh(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Cosh(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *CoshFunction) GetName() string { return "cosh" }
var coshFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *CoshFunction) GetParams() []data.GetValue {
	return coshFunctionGetParams
}
var coshFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *CoshFunction) GetVariables() []data.Variable {
	return coshFunctionGetVariables
}

// SinhFunction 实现 sinh 函数
type SinhFunction struct{}

func NewSinhFunction() data.FuncStmt { return &SinhFunction{} }

func (f *SinhFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Sinh(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Sinh(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *SinhFunction) GetName() string { return "sinh" }
var sinhFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *SinhFunction) GetParams() []data.GetValue {
	return sinhFunctionGetParams
}
var sinhFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *SinhFunction) GetVariables() []data.Variable {
	return sinhFunctionGetVariables
}

// TanhFunction 实现 tanh 函数
type TanhFunction struct{}

func NewTanhFunction() data.FuncStmt { return &TanhFunction{} }

func (f *TanhFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Tanh(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Tanh(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *TanhFunction) GetName() string { return "tanh" }
var tanhFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *TanhFunction) GetParams() []data.GetValue {
	return tanhFunctionGetParams
}
var tanhFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *TanhFunction) GetVariables() []data.Variable {
	return tanhFunctionGetVariables
}

// AcoshFunction 实现 acosh 函数
type AcoshFunction struct{}

func NewAcoshFunction() data.FuncStmt { return &AcoshFunction{} }

func (f *AcoshFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Acosh(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Acosh(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *AcoshFunction) GetName() string { return "acosh" }
var acoshFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *AcoshFunction) GetParams() []data.GetValue {
	return acoshFunctionGetParams
}
var acoshFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *AcoshFunction) GetVariables() []data.Variable {
	return acoshFunctionGetVariables
}

// AsinhFunction 实现 asinh 函数
type AsinhFunction struct{}

func NewAsinhFunction() data.FuncStmt { return &AsinhFunction{} }

func (f *AsinhFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Asinh(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Asinh(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *AsinhFunction) GetName() string { return "asinh" }
var asinhFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *AsinhFunction) GetParams() []data.GetValue {
	return asinhFunctionGetParams
}
var asinhFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *AsinhFunction) GetVariables() []data.Variable {
	return asinhFunctionGetVariables
}

// AtanhFunction 实现 atanh 函数
type AtanhFunction struct{}

func NewAtanhFunction() data.FuncStmt { return &AtanhFunction{} }

func (f *AtanhFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Atanh(fv)), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(math.Atanh(float64(ai.Value))), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *AtanhFunction) GetName() string { return "atanh" }
var atanhFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *AtanhFunction) GetParams() []data.GetValue {
	return atanhFunctionGetParams
}
var atanhFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *AtanhFunction) GetVariables() []data.Variable {
	return atanhFunctionGetVariables
}

// HypotFunction 实现 hypot 函数
type HypotFunction struct{}

func NewHypotFunction() data.FuncStmt { return &HypotFunction{} }

func (f *HypotFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	xV, _ := ctx.GetIndexValue(0)
	yV, _ := ctx.GetIndexValue(1)
	var x, y float64
	if af, ok := xV.(data.AsFloat); ok {
		x, _ = af.AsFloat()
	} else if ai, ok := xV.(*data.IntValue); ok {
		x = float64(ai.Value)
	}
	if af, ok := yV.(data.AsFloat); ok {
		y, _ = af.AsFloat()
	} else if ai, ok := yV.(*data.IntValue); ok {
		y = float64(ai.Value)
	}
	return data.NewFloatValue(math.Hypot(x, y)), nil
}

func (f *HypotFunction) GetName() string { return "hypot" }
var hypotFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "x", 0, nil, nil),
	node.NewParameter(nil, "y", 1, nil, nil),
}

func (f *HypotFunction) GetParams() []data.GetValue {
	return hypotFunctionGetParams
}
var hypotFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "x", 0, nil),
	node.NewVariable(nil, "y", 1, nil),
}

func (f *HypotFunction) GetVariables() []data.Variable {
	return hypotFunctionGetVariables
}

// FmodFunction 实现 fmod 函数
type FmodFunction struct{}

func NewFmodFunction() data.FuncStmt { return &FmodFunction{} }

func (f *FmodFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	xV, _ := ctx.GetIndexValue(0)
	yV, _ := ctx.GetIndexValue(1)
	var x, y float64
	if af, ok := xV.(data.AsFloat); ok {
		x, _ = af.AsFloat()
	} else if ai, ok := xV.(*data.IntValue); ok {
		x = float64(ai.Value)
	}
	if af, ok := yV.(data.AsFloat); ok {
		y, _ = af.AsFloat()
	} else if ai, ok := yV.(*data.IntValue); ok {
		y = float64(ai.Value)
	}
	if y == 0 {
		return data.NewFloatValue(math.NaN()), nil
	}
	return data.NewFloatValue(math.Mod(x, y)), nil
}

func (f *FmodFunction) GetName() string { return "fmod" }
var fmodFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "x", 0, nil, nil),
	node.NewParameter(nil, "y", 1, nil, nil),
}

func (f *FmodFunction) GetParams() []data.GetValue {
	return fmodFunctionGetParams
}
var fmodFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "x", 0, nil),
	node.NewVariable(nil, "y", 1, nil),
}

func (f *FmodFunction) GetVariables() []data.Variable {
	return fmodFunctionGetVariables
}

// ---------- 角度转换函数 ----------

// Deg2radFunction 实现 deg2rad 函数
type Deg2radFunction struct{}

func NewDeg2radFunction() data.FuncStmt { return &Deg2radFunction{} }

func (f *Deg2radFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(fv * math.Pi / 180), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(float64(ai.Value) * math.Pi / 180), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *Deg2radFunction) GetName() string { return "deg2rad" }
var deg2radFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *Deg2radFunction) GetParams() []data.GetValue {
	return deg2radFunctionGetParams
}
var deg2radFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *Deg2radFunction) GetVariables() []data.Variable {
	return deg2radFunctionGetVariables
}

// Rad2degFunction 实现 rad2deg 函数
type Rad2degFunction struct{}

func NewRad2degFunction() data.FuncStmt { return &Rad2degFunction{} }

func (f *Rad2degFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(fv * 180 / math.Pi), nil
	}
	if ai, ok := v.(*data.IntValue); ok {
		return data.NewFloatValue(float64(ai.Value) * 180 / math.Pi), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *Rad2degFunction) GetName() string { return "rad2deg" }
var rad2degFunctionGetParams = []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}

func (f *Rad2degFunction) GetParams() []data.GetValue {
	return rad2degFunctionGetParams
}
var rad2degFunctionGetVariables = []data.Variable{node.NewVariable(nil, "num", 0, nil)}

func (f *Rad2degFunction) GetVariables() []data.Variable {
	return rad2degFunctionGetVariables
}

// ---------- 进制转换函数 ----------

// BaseConvertFunction 实现 base_convert 函数
type BaseConvertFunction struct{}

func NewBaseConvertFunction() data.FuncStmt { return &BaseConvertFunction{} }

func (f *BaseConvertFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	numV, _ := ctx.GetIndexValue(0)
	fromV, _ := ctx.GetIndexValue(1)
	toV, _ := ctx.GetIndexValue(2)

	if numV == nil || fromV == nil || toV == nil {
		return data.NewStringValue(""), nil
	}

	numStr := numV.AsString()
	var fromBase, toBase int
	if iv, ok := fromV.(*data.IntValue); ok {
		fromBase = iv.Value
	} else if ai, ok := fromV.(data.AsInt); ok {
		iv, _ := ai.AsInt()
		fromBase = iv
	}
	if iv, ok := toV.(*data.IntValue); ok {
		toBase = iv.Value
	} else if ai, ok := toV.(data.AsInt); ok {
		iv, _ := ai.AsInt()
		toBase = iv
	}

	if fromBase < 2 || fromBase > 36 || toBase < 2 || toBase > 36 {
		return data.NewStringValue(""), nil
	}

	// 先转换到十进制
	decimal, err := strconv.ParseInt(numStr, fromBase, 64)
	if err != nil {
		// 尝试无符号
		udecimal, uerr := strconv.ParseUint(numStr, fromBase, 64)
		if uerr != nil {
			return data.NewStringValue(""), nil
		}
		return data.NewStringValue(strconv.FormatUint(udecimal, toBase)), nil
	}
	return data.NewStringValue(strconv.FormatInt(decimal, toBase)), nil
}

func (f *BaseConvertFunction) GetName() string { return "base_convert" }
var baseConvertFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "num", 0, nil, nil),
	node.NewParameter(nil, "from_base", 1, nil, data.NewBaseType("int")),
	node.NewParameter(nil, "to_base", 2, nil, data.NewBaseType("int")),
}

func (f *BaseConvertFunction) GetParams() []data.GetValue {
	return baseConvertFunctionGetParams
}
var baseConvertFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "num", 0, nil),
	node.NewVariable(nil, "from_base", 1, data.NewBaseType("int")),
	node.NewVariable(nil, "to_base", 2, data.NewBaseType("int")),
}

func (f *BaseConvertFunction) GetVariables() []data.Variable {
	return baseConvertFunctionGetVariables
}

// BindecFunction 实现 bindec 函数
type BindecFunction struct{}

func NewBindecFunction() data.FuncStmt { return &BindecFunction{} }

func (f *BindecFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}
	str := v.AsString()
	val, err := strconv.ParseInt(str, 2, 64)
	if err != nil {
		udec, uerr := strconv.ParseUint(str, 2, 64)
		if uerr != nil {
			return data.NewIntValue(0), nil
		}
		return data.NewFloatValue(float64(udec)), nil
	}
	return data.NewIntValue(int(val)), nil
}

func (f *BindecFunction) GetName() string { return "bindec" }
var bindecFunctionGetParams = []data.GetValue{node.NewParameter(nil, "binary_string", 0, nil, nil)}

func (f *BindecFunction) GetParams() []data.GetValue {
	return bindecFunctionGetParams
}
var bindecFunctionGetVariables = []data.Variable{node.NewVariable(nil, "binary_string", 0, nil)}

func (f *BindecFunction) GetVariables() []data.Variable {
	return bindecFunctionGetVariables
}
