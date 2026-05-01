package data

// TraitAlias 表示 trait use 语句中的方法别名
// 语法: use Trait { method as alias; }
type TraitAlias struct {
	Method string // 原方法名
	Alias  string // 别名
}
