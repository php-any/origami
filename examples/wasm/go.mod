module wasm-example

go 1.24

toolchain go1.24.6

// 始终依赖远端主干；构建脚本会执行 `go get github.com/php-any/origami@main` 进行刷新
require github.com/php-any/origami v0.0.0-00010101000000-000000000000


