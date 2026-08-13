package annotation

import "github.com/php-any/origami/data"

const (
	lifetimeTransient = iota
	lifetimeSingleton
	lifetimeScoped
)

// RegisterClassLifetime 由 std/container 在 Load 时注入，用于注册类级生命周期注解。
var RegisterClassLifetime func(ctx data.Context, lifetime int) data.Control

// BindClassAnnotation 由 std/container 在 Load 时注入，用于 #[Bind] 注解。
var BindClassAnnotation func(ctx data.Context) data.Control

// InjectParameterAnnotation 由 std/container 在 Load 时注入，用于 #[Inject] 注解。
var InjectParameterAnnotation func(ctx data.Context) data.Control

// NamedParameterAnnotation 由 std/container 在 Load 时注入，用于 #[Named] 注解。
var NamedParameterAnnotation func(ctx data.Context) data.Control
