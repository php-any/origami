package annotation

import "github.com/php-any/origami/data"

// CompiledRouteValue 构建预编译的 @Route 注解实例（供 compile 子命令使用）
func CompiledRouteValue(prefix string) *data.ClassValue {
	r := &Route{prefix: prefix}
	rc := &RouteClass{source: r, construct: &RouteConstructMethod{route: r}}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: rc}
}

// CompiledControllerValue 构建预编译的 @Controller 注解实例
func CompiledControllerValue(name string) *data.ClassValue {
	c := &Controller{name: name}
	cc := &ControllerClass{
		process:   &ControllerProcessMethod{controller: c},
		register:  &ControllerRegisterMethod{controller: c},
		construct: &ControllerConstructMethod{controller: c},
	}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: cc}
}

// CompiledGetMappingValue 构建预编译的 @GetMapping 注解实例
func CompiledGetMappingValue(path string) *data.ClassValue {
	m := &GetMapping{path: path}
	gc := &GetMappingClass{source: m, construct: &GetMappingConstructMethod{mapping: m}}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: gc}
}

// CompiledPostMappingValue 构建预编译的 @PostMapping 注解实例
func CompiledPostMappingValue(path string) *data.ClassValue {
	m := &PostMapping{path: path}
	pc := &PostMappingClass{source: m, construct: &PostMappingConstructMethod{mapping: m}}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: pc}
}

// CompiledPutMappingValue 构建预编译的 @PutMapping 注解实例
func CompiledPutMappingValue(path string) *data.ClassValue {
	m := &PutMapping{path: path}
	pc := &PutMappingClass{source: m, construct: &PutMappingConstructMethod{mapping: m}}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: pc}
}

// CompiledDeleteMappingValue 构建预编译的 @DeleteMapping 注解实例
func CompiledDeleteMappingValue(path string) *data.ClassValue {
	m := &DeleteMapping{path: path}
	dc := &DeleteMappingClass{source: m, construct: &DeleteMappingConstructMethod{mapping: m}}
	return &data.ClassValue{ObjectValue: data.NewObjectValue(), Class: dc}
}
