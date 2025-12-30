package data

// Control 表示基础控制流接口
type Control interface {
	Value
}

type AddStack interface {
	AddStackWithInfo(f From, className, methodName string)
}

// BreakControl 表示中断语句控制流
type BreakControl interface {
	Control
	// IsBreak 是否为中断语句
	IsBreak() bool
	// GetLabel 获取中断标签
	GetLabel() string
}

// ContinueControl 表示继续语句控制流
type ContinueControl interface {
	Control
	// IsContinue 是否为继续语句
	IsContinue() bool
	// GetLabel 获取继续标签
	GetLabel() string
}

// ExitControl 表示程序退出控制流
type ExitControl interface {
	Control
	// IsExit 是否为程序退出
	IsExit() bool
	// GetCode 获取退出码
	GetCode() int
}

// GotoControl 表示跳转语句控制流
type GotoControl interface {
	Control
	// IsGoto 是否为跳转语句
	IsGoto() bool
	// GetLabel 获取跳转标签
	GetLabel() string
}

// YieldValueControl 表示生成器yield控制流
type YieldValueControl interface {
	Control
	GetYieldKey() Value
	GetYieldValue() Value
}

// YieldControl 函数Yield中断
type YieldControl interface {
	Control
	// GetBodyStackState 模拟返回堆栈状态的值
	GetBodyStackState(ctx Context) Generator
	// SetBodyIndex 添加函数 Body 索引状态
	SetBodyIndex(index int, body []GetValue)
	// GetBodyIndex 获取函数 Body 索引
	GetBodyIndex() int
}
