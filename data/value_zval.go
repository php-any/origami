package data

import "fmt"

func NewZValValue(v *ZVal) *ZValValue {
	return &ZValValue{ZVal: v}
}

type ZValValue struct {
	ZVal *ZVal
}

func (Z *ZValValue) GetValue(_ Context) (GetValue, Control) {
	return Z, nil
}

func (Z *ZValValue) AsString() string {
	return fmt.Sprintf("%d", Z.ZVal)
}

func (Z *ZValValue) Marshal(serializer Serializer) ([]byte, error) {
	if v, ok := Z.ZVal.Value.(ValueSerializer); ok {
		return v.Marshal(serializer)
	}
	return nil, fmt.Errorf("cannot marshal ZValValue")
}
func (Z *ZValValue) Unmarshal(data []byte, serializer Serializer) error {
	if v, ok := Z.ZVal.Value.(ValueSerializer); ok {
		return v.Unmarshal(data, serializer)
	}
	return fmt.Errorf("cannot unmarshal ZValValue")
}

func (Z *ZValValue) ToGoValue(serializer Serializer) (any, error) {
	if v, ok := Z.ZVal.Value.(ValueSerializer); ok {
		return v.ToGoValue(serializer)
	}

	return nil, fmt.Errorf("cannot marshal ZValValue")
}
