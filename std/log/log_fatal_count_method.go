package log

import (
	"github.com/php-any/origami/data"
)

// LogFatalCountMethod 实现 Log::fatalCount(): int
//
// 供测试套件统计「本文件内是否触发过 Log::fatal」。fatal 现在是可被
// catch 捕获的异常，用例自身的 catch (\Throwable) 可能把它吞掉，
// 用计数兜底才能保证失败一定被上报。
type LogFatalCountMethod struct {
	source *Log
}

func (h *LogFatalCountMethod) Call(_ data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(int(FatalCount())), nil
}

func (h *LogFatalCountMethod) GetName() string {
	return "fatalCount"
}

func (h *LogFatalCountMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *LogFatalCountMethod) GetIsStatic() bool {
	return true
}

func (h *LogFatalCountMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *LogFatalCountMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *LogFatalCountMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}
