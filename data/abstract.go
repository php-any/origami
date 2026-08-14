package data

// AbstractMethodMarker 标记抽象方法。
// 由 node 包的 AbstractMethod 类型实现，用于在 data 包中
// 通过类型断言判断一个方法是否为抽象方法（避免 data -> node 的循环依赖）。
type AbstractMethodMarker interface {
	IsAbstractMethod() bool
}
