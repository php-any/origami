package data

// Control 表示基础控制流接口
type Control interface {
	Value
}

type AddStack interface {
	AddStack(f From)
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

// YieldControl 表示生成器yield控制流
type YieldControl interface {
	Control
	// IsYield 是否为生成器yield
	IsYield() bool
	// GetKey 获取生成器键
	GetKey() Value
}
