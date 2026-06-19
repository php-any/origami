package compile

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// 新增 AST 节点的扩展方式见 doc.go。

// Emit 将 AST 节点反射转译为 Go 源码字面量。
func (g *Generator) Emit(v data.GetValue) error {
	if v == nil {
		g.printf("nil")
		return nil
	}

	typ := reflect.TypeOf(v)
	if h, ok := specialHandlers[typ]; ok {
		return h(g, v)
	}

	if emit, ok := dataValueEmitters[typ]; ok {
		return emit(g, v)
	}

	if err := g.emitStructLiteral(v); err == nil {
		return nil
	}

	return newEmitError(g.file, v, "no special handler and cannot emit via constructor or struct literal")
}

// genGetValue 是 Emit 的别名，供遗留辅助函数调用。
func (g *Generator) genGetValue(v data.GetValue) error {
	return g.Emit(v)
}

func (g *Generator) emitStructLiteral(v data.GetValue) error {
	rv := reflect.ValueOf(v)
	typ := rv.Type()
	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer type")
	}
	elem := typ.Elem()
	if elem.PkgPath() == "" {
		return fmt.Errorf("no package for %s", elem.Name())
	}
	pkgName := elem.PkgPath()[strings.LastIndex(elem.PkgPath(), "/")+1:]
	if pkgName == "data" || strings.HasSuffix(elem.PkgPath(), "/data") {
		pkgName = "data"
	} else if strings.HasSuffix(elem.PkgPath(), "/node") {
		pkgName = "node"
	}

	needsNode := false
	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)
		if f.Anonymous && f.Name == "Node" && f.Tag.Get("pp") == "-" {
			needsNode = true
		}
		if !f.IsExported() && !(f.Anonymous && f.Name == "Node") {
			return fmt.Errorf("unexported field %s", f.Name)
		}
	}

	g.printf("&%s.%s{\n", pkgName, elem.Name())
	g.indent++
	if needsNode {
		g.printf("Node: node.NewNode(from),\n")
	}

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.Anonymous && field.Name == "Node" {
			continue
		}
		if field.Tag.Get("pp") == "-" {
			continue
		}
		if !field.IsExported() {
			return fmt.Errorf("unexported field %s", field.Name)
		}
		fv := rv.Elem().Field(i)
		g.printf("%s: ", field.Name)
		if err := g.emitReflectValue(fv); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}")
	return nil
}

func (g *Generator) emitReflectValue(rv reflect.Value) error {
	if !rv.IsValid() {
		g.printf("nil")
		return nil
	}

	// 接口 / 指针 / 值 统一处理
	if rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			g.printf("nil")
			return nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			g.printf("nil")
			return nil
		}
		if rv.Type().Implements(reflect.TypeOf((*data.GetValue)(nil)).Elem()) {
			return g.Emit(rv.Interface().(data.GetValue))
		}
	}

	// data.Variable 接口
	if rv.Type().Implements(reflect.TypeOf((*data.Variable)(nil)).Elem()) {
		return g.emitVariable(rv.Interface().(data.Variable))
	}

	// data.Types 接口
	if rv.Type().Implements(reflect.TypeOf((*data.Types)(nil)).Elem()) {
		g.genTypes(rv.Interface().(data.Types))
		return nil
	}

	switch rv.Kind() {
	case reflect.String:
		g.printf("%q", rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		g.printf("%d", rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		g.printf("%d", rv.Uint())
	case reflect.Bool:
		g.printf("%v", rv.Bool())
	case reflect.Float32, reflect.Float64:
		g.printf("%g", rv.Float())
	case reflect.Slice:
		return g.emitSlice(rv)
	case reflect.Map:
		return g.emitMap(rv)
	case reflect.Struct:
		return g.emitStructValue(rv)
	case reflect.Ptr:
		if rv.Type().Implements(reflect.TypeOf((*data.GetValue)(nil)).Elem()) {
			return g.Emit(rv.Interface().(data.GetValue))
		}
		return g.emitStructLiteral(rv.Interface().(data.GetValue))
	default:
		return fmt.Errorf("unsupported reflect kind %s for %s", rv.Kind(), rv.Type())
	}
	return nil
}

func (g *Generator) emitSlice(rv reflect.Value) error {
	if rv.Type().Elem().Kind() == reflect.Uint8 {
		g.printf("%q", string(rv.Bytes()))
		return nil
	}

	elemType := rv.Type().Elem()
	pkg := ""
	typeName := elemType.Name()
	if elemType.PkgPath() != "" {
		if strings.HasSuffix(elemType.PkgPath(), "/node") {
			pkg = "node."
		} else if strings.HasSuffix(elemType.PkgPath(), "/data") {
			pkg = "data."
		}
	}

	if elemType.Kind() == reflect.Interface && elemType.NumMethod() > 0 {
		// []data.GetValue
		if elemType.Implements(reflect.TypeOf((*data.GetValue)(nil)).Elem()) {
			g.printf("[]data.GetValue{\n")
			g.indent++
			for i := 0; i < rv.Len(); i++ {
				if err := g.Emit(rv.Index(i).Interface().(data.GetValue)); err != nil {
					return err
				}
				g.printf(",\n")
			}
			g.indent--
			g.printf("}")
			return nil
		}
		if elemType.Implements(reflect.TypeOf((*data.Variable)(nil)).Elem()) {
			g.printf("[]data.Variable{\n")
			g.indent++
			for i := 0; i < rv.Len(); i++ {
				if err := g.emitVariable(rv.Index(i).Interface().(data.Variable)); err != nil {
					return err
				}
				g.printf(",\n")
			}
			g.indent--
			g.printf("}")
			return nil
		}
		if elemType.Implements(reflect.TypeOf((*data.Method)(nil)).Elem()) {
			g.printf("[]data.Method{\n")
			g.indent++
			for i := 0; i < rv.Len(); i++ {
				if err := g.emitMethod(rv.Index(i).Interface().(data.Method)); err != nil {
					return err
				}
				g.printf(",\n")
			}
			g.indent--
			g.printf("}")
			return nil
		}
	}

	g.printf("[]%s%s{\n", pkg, typeName)
	g.indent++
	for i := 0; i < rv.Len(); i++ {
		if err := g.emitReflectValue(rv.Index(i)); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}")
	return nil
}

func (g *Generator) emitMap(rv reflect.Value) error {
	keyType := rv.Type().Key()
	valType := rv.Type().Elem()
	if keyType.Kind() != reflect.String {
		return fmt.Errorf("unsupported map key type %s", keyType)
	}

	valPkg := ""
	if valType.PkgPath() != "" {
		if strings.HasSuffix(valType.PkgPath(), "/data") {
			valPkg = "data."
		} else if strings.HasSuffix(valType.PkgPath(), "/node") {
			valPkg = "node."
		}
	}

	g.printf("map[string]%s%s{\n", valPkg, valType.Name())
	g.indent++
	for _, key := range rv.MapKeys() {
		g.printf("%q: ", key.String())
		if err := g.emitReflectValue(rv.MapIndex(key)); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}")
	return nil
}

func (g *Generator) emitStructValue(rv reflect.Value) error {
	typ := rv.Type()
	pkgName := "node"
	if strings.HasSuffix(typ.PkgPath(), "/data") {
		pkgName = "data"
	} else if strings.HasSuffix(typ.PkgPath(), "/node") {
		pkgName = "node"
	}

	needsNode := false
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Anonymous && f.Name == "Node" {
			needsNode = true
		}
	}

	g.printf("%s.%s{\n", pkgName, typ.Name())
	g.indent++
	if needsNode {
		g.printf("Node: node.NewNode(from),\n")
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous && field.Name == "Node" {
			continue
		}
		if !field.IsExported() {
			return fmt.Errorf("unexported struct field %s", field.Name)
		}
		g.printf("%s: ", field.Name)
		if err := g.emitReflectValue(rv.Field(i)); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}")
	return nil
}

func (g *Generator) emitVariable(v data.Variable) error {
	if v == nil {
		g.printf("nil")
		return nil
	}
	if gv, ok := v.(data.GetValue); ok {
		return g.Emit(gv)
	}
	g.printf("data.NewVariable(%q, %d, nil)", v.GetName(), v.GetIndex())
	return nil
}

func (g *Generator) emitMethod(method data.Method) error {
	if method == nil {
		g.printf("nil")
		return nil
	}
	if im, ok := method.(*node.InterfaceMethod); ok {
		return g.emitInterfaceMethod(im)
	}
	if am, ok := method.(*node.AbstractMethod); ok {
		g.printf("node.NewAbstractMethod(")
		if err := g.emitClassMethod(am.ClassMethod); err != nil {
			return err
		}
		g.printf(")")
		return nil
	}
	if cm, ok := method.(*node.ClassMethod); ok {
		return g.emitClassMethod(cm)
	}
	return newEmitError(g.file, nil, fmt.Sprintf("unsupported method type %T", method))
}
