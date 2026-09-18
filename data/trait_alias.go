package data

// TraitAlias 表示 trait use 语句中的方法别名
// 语法: use Trait { method as alias; } 或 { Trait::method as alias; }
type TraitAlias struct {
	Trait  string // 可选：TraitName::method 中的 trait 名（已解析的 FQCN 或短名）
	Method string // 原方法名
	Alias  string // 别名
}
