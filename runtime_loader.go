package main

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/laravel/framework"
	netannotation "github.com/php-any/origami/std/net/annotation"
	"github.com/php-any/origami/std/net/grpc"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/tcp"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func loadDefaultRuntime(vm data.VM) {
	std.Load(vm)
	framework.Load(vm)
	php.Load(vm)
	http.Load(vm)
	websocket.Load(vm)
	tcp.Load(vm)
	grpc.Load(vm)
	netannotation.Load(vm)
	system.Load(vm)
}
