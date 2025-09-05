module github.com/php-any/origami/tools/lsp

go 1.24.4

toolchain go1.24.6

require (
	github.com/php-any/origami v0.0.0-20250807021948-61648a4c5484
	github.com/sirupsen/logrus v1.9.3
	github.com/sourcegraph/jsonrpc2 v0.2.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/sys v0.34.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/php-any/origami => ../../
