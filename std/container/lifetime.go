package container

// Lifetime 服务生命周期。
type Lifetime int

const (
	LifetimeTransient Lifetime = iota
	LifetimeSingleton
	LifetimeScoped
)

func (l Lifetime) shared() bool {
	return l == LifetimeSingleton || l == LifetimeScoped
}
