package data

type Serializer interface {
	MarshalInt(*IntValue) ([]byte, error)
	UnmarshalInt([]byte, *IntValue) error
	MarshalString(*StringValue) ([]byte, error)
	UnmarshalString([]byte, *StringValue) error
	MarshalNull(*NullValue) ([]byte, error)
	UnmarshalNull([]byte, *NullValue) error
	MarshalArray(*ArrayValue) ([]byte, error)
	UnmarshalArray([]byte, *ArrayValue) error
	MarshalObject(*ObjectValue) ([]byte, error)
	UnmarshalObject([]byte, *ObjectValue) error
	MarshalBool(*BoolValue) ([]byte, error)
	UnmarshalBool([]byte, *BoolValue) error
	MarshalFloat(*FloatValue) ([]byte, error)
	UnmarshalFloat([]byte, *FloatValue) error
	MarshalAny(*AnyValue) ([]byte, error)
	UnmarshalAny([]byte, *AnyValue) error
	MarshalClass(*ClassValue) ([]byte, error)
	UnmarshalClass([]byte, *ClassValue) error
}

// 每个Value类型都需要实现这个接口
type ValueSerializer interface {
	Marshal(serializer Serializer) ([]byte, error)
	Unmarshal(data []byte, serializer Serializer) error
}
