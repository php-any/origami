package stream

import "github.com/php-any/origami/data"

type StreamGetWrappersFunction struct{}

func (*StreamGetWrappersFunction) GetName() string               { return "stream_get_wrappers" }
func (*StreamGetWrappersFunction) GetParams() []data.GetValue    { return nil }
func (*StreamGetWrappersFunction) GetVariables() []data.Variable { return nil }
func (*StreamGetWrappersFunction) Call(data.Context) (data.GetValue, data.Control) {
	return data.NewArrayValue([]data.Value{data.NewStringValue("file"), data.NewStringValue("php"), data.NewStringValue("phar"), data.NewStringValue("http"), data.NewStringValue("https")}), nil
}
