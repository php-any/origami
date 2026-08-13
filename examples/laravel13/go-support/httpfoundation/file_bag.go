package httpfoundation

import (
	"fmt"
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var fileKeys = []string{"error", "full_path", "name", "size", "tmp_name", "type"}

// FileBagClass 实现 Symfony\Component\HttpFoundation\FileBag。
type FileBagClass struct {
	node.Node
	source     *ParamBagData
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

func NewFileBagClass() data.ClassStmt {
	return NewFileBagClassFrom(nil)
}

func NewFileBagClassFrom(source *ParamBagData) data.ClassStmt {
	c := &FileBagClass{
		source:     source,
		properties: []data.Property{protectedArrayProp("parameters")},
	}
	c.methods, c.methodList = fileBagMethods()
	return c
}

func (c *FileBagClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := c.source
	if src == nil {
		src = newParamBagData()
	} else {
		src = src.clone()
	}
	return data.NewProxyValue(NewFileBagClassFrom(src), ctx.CreateBaseContext()), nil
}

func (c *FileBagClass) GetName() string {
	return fqnFileBag
}
func (c *FileBagClass) GetExtend() *string {
	parent := fqnParameterBag
	return &parent
}
func (c *FileBagClass) GetImplements() []string          { return nil }
func (c *FileBagClass) GetSource() any                   { return c.source }
func (c *FileBagClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *FileBagClass) GetPropertyList() []data.Property { return c.properties }
func (c *FileBagClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *FileBagClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *FileBagClass) GetMethods() []data.Method { return c.methodList }

func fileBagMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{param("parameters", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("parameters", 0, nil)},
			nil, fileBagConstruct),
		pubMethod("replace",
			[]data.GetValue{param("files", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("files", 0, nil)},
			nil, fileBagReplaceMethod),
		pubMethod("set",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("value", 1, nil, nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("value", 1, nil),
			},
			nil, fileBagSet),
		pubMethod("add",
			[]data.GetValue{param("files", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("files", 0, nil)},
			nil, fileBagAdd),
	}
	m := make(map[string]data.Method, len(list))
	for _, method := range list {
		m[method.GetName()] = method
	}
	return m, list
}

func fileBagConstruct(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	m, err := valueToAssocMap(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, fileBagReplace(bagClassValue(ctx), m)
}

func fileBagReplaceMethod(ctx data.Context) (data.GetValue, data.Control) {
	raw, _ := ctx.GetIndexValue(0)
	m, err := valueToAssocMap(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, fileBagReplace(bagClassValue(ctx), m)
}

func fileBagReplace(cv *data.ClassValue, files map[string]data.Value) data.Control {
	store := ParamBagFrom(cv)
	if store == nil {
		return nil
	}
	store.replace(map[string]data.Value{})
	for k, v := range files {
		converted, ctl := convertFileInformation(v)
		if ctl != nil {
			return ctl
		}
		if converted != nil {
			store.set(k, converted)
		}
	}
	if cv != nil {
		cv.ObjectValue.SetProperty("parameters", store.toArrayValue())
	}
	return nil
}

func fileBagSet(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	val := defaultValueParam(ctx, 1, data.NewNullValue())
	if !isArrayValue(val) && !isUploadedFile(val) {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("An uploaded file must be an array or an instance of UploadedFile."))
	}
	converted, ctl := convertFileInformation(val)
	if ctl != nil {
		return nil, ctl
	}
	if converted == nil {
		store.remove(key)
	} else {
		store.set(key, converted)
	}
	syncParametersProperty(ctx, store)
	return nil, nil
}

func fileBagAdd(ctx data.Context) (data.GetValue, data.Control) {
	store := paramData(ctx)
	if store == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	m, err := valueToAssocMap(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	for k, v := range m {
		if !isArrayValue(v) && !isUploadedFile(v) {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("An uploaded file must be an array or an instance of UploadedFile."))
		}
		converted, ctl := convertFileInformation(v)
		if ctl != nil {
			return nil, ctl
		}
		if converted == nil {
			store.remove(k)
		} else {
			store.set(k, converted)
		}
	}
	syncParametersProperty(ctx, store)
	return nil, nil
}

func isUploadedFile(v data.Value) bool {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return false
	}
	name := cv.GetName()
	return name == "Symfony\\Component\\HttpFoundation\\File\\UploadedFile" ||
		stringsHasSuffix(name, "\\UploadedFile")
}

func stringsHasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func convertFileInformation(file data.Value) (data.Value, data.Control) {
	if file == nil || isNull(file) {
		return nil, nil
	}
	if isUploadedFile(file) {
		return file, nil
	}
	m, err := valueToAssocMap(file)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	m = fixPhpFilesArray(m)

	keys := make([]string, 0, len(m)+1)
	for k := range m {
		keys = append(keys, k)
	}
	if _, ok := m["full_path"]; !ok {
		keys = append(keys, "full_path")
	}
	sort.Strings(keys)

	if equalStringSlices(keys, fileKeys) {
		errVal := 0
		if e, ok := m["error"]; ok {
			if iv, ok := e.(data.AsInt); ok {
				errVal, _ = iv.AsInt()
			}
		}
		if errVal == 4 { // UPLOAD_ERR_NO_FILE
			return nil, nil
		}
		// 未实现 UploadedFile Go 类时保留规范化数组；有类时尝试构造
		return assocMapToArrayValue(m), nil
	}

	outKeys := make([]string, 0, len(m))
	outVals := make(map[string]data.Value, len(m))
	isList := true
	i := 0
	for k, v := range m {
		outKeys = append(outKeys, k)
		if k != fmt.Sprintf("%d", i) {
			isList = false
		}
		i++
		if isUploadedFile(v) || isArrayValue(v) {
			converted, ctl := convertFileInformation(v)
			if ctl != nil {
				return nil, ctl
			}
			if converted != nil {
				outVals[k] = converted
			}
		} else {
			outVals[k] = v
		}
	}
	if isList {
		filtered := make([]*data.ZVal, 0, len(outKeys))
		for _, k := range outKeys {
			if v, ok := outVals[k]; ok && v != nil && !isNull(v) {
				filtered = append(filtered, data.NewZVal(v))
			}
		}
		return &data.ArrayValue{List: filtered}, nil
	}
	return orderedAssocToArrayValue(outKeys, outVals), nil
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func fixPhpFilesArray(dataMap map[string]data.Value) map[string]data.Value {
	keys := make([]string, 0, len(dataMap)+1)
	for k := range dataMap {
		keys = append(keys, k)
	}
	if _, ok := dataMap["full_path"]; !ok {
		keys = append(keys, "full_path")
	}
	sort.Strings(keys)
	if !equalStringSlices(keys, fileKeys) {
		return dataMap
	}
	nameVal, ok := dataMap["name"]
	if !ok || !isArrayValue(nameVal) {
		return dataMap
	}

	nameMap, _ := valueToAssocMap(nameVal)
	files := map[string]data.Value{}
	for k := range dataMap {
		isFileKey := false
		for _, fk := range fileKeys {
			if k == fk {
				isFileKey = true
				break
			}
		}
		if !isFileKey {
			files[k] = dataMap[k]
		}
	}

	errorMap, _ := valueToAssocMap(dataMap["error"])
	typeMap, _ := valueToAssocMap(dataMap["type"])
	tmpMap, _ := valueToAssocMap(dataMap["tmp_name"])
	sizeMap, _ := valueToAssocMap(dataMap["size"])
	var fullPathMap map[string]data.Value
	if fp, ok := dataMap["full_path"]; ok {
		fullPathMap, _ = valueToAssocMap(fp)
	}

	for k, name := range nameMap {
		entry := map[string]data.Value{
			"error":    errorMap[k],
			"name":     name,
			"type":     typeMap[k],
			"tmp_name": tmpMap[k],
			"size":     sizeMap[k],
		}
		if fullPathMap != nil {
			if v, ok := fullPathMap[k]; ok {
				entry["full_path"] = v
			}
		}
		files[k] = assocMapToArrayValue(fixPhpFilesArray(entry))
	}
	return files
}
