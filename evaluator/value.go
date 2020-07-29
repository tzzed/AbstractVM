package evaluator

import "fmt"

type ValueType uint8

const (
	CharValue ValueType = 0x10

	ShortValue = 0x20

	// integer family: 0x10 to 0x1F
	IntegerValue = 0x30

	FloatValue = 0x40

	// double family: 0x20 to 0x2F
	DoubleValue = 0x50
)

type Value struct {
	V    interface{}
	Type ValueType
}

func (t ValueType) String() string {
	switch t {
	case CharValue:
		return "int8"
	case ShortValue:
		return "int16"
	case IntegerValue:
		return "int32"
	case FloatValue:
		return "float"
	case DoubleValue:
		return "double"
	}

	return ""
}

func NewInt8Value(x int8) Value {
	return Value{V: x,
		Type: CharValue}
}

func NewInt16Value(x int16) Value {
	return Value{V: x,
		Type: ShortValue}
}

func NewInt32Value(x int32) Value {
	return Value{V: int32(x), Type: IntegerValue}
}

func NewFloatValue(x float32) Value {
	return Value{V: x,
		Type: FloatValue}
}

func NewDoubleValue(x float64) Value {
	return Value{V: x,
		Type: DoubleValue}
}

func GetBiggerType(a, b Value) ValueType {
	if a.Type > b.Type {
		return a.Type
	}

	return b.Type
}

func (v Value) ConvertToInteger() (int32, error) {
	switch v.Type {
	case CharValue:
		return int32(v.V.(int8)), nil
	case ShortValue:
		return int32(v.V.(int16)), nil
	case IntegerValue:
		return int32(v.V.(int32)), nil
	}

	return 0, fmt.Errorf("cannot convert Type %d into int32", v.Type)
}

func (v Value) ConvertToFloat() (float32, error) {
	fmt.Println(v.Type)
	switch v.Type {
	case CharValue:
		return float32(v.V.(int8)), nil
	case ShortValue:
		return float32(v.V.(int16)), nil
	case IntegerValue:
		return float32(v.V.(int32)), nil
	case FloatValue:
		return v.V.(float32), nil
	}

	return 0, fmt.Errorf("cannot convert Type %d into float32", v.Type)
}

func (v Value) ConvertToShort() (int16, error) {
	switch v.Type {
	case CharValue:
		return int16(v.V.(int8)), nil
	case ShortValue:
		return int16(v.V.(int16)), nil
	}

	return 0, fmt.Errorf("cannot convert Type %d into in16", v.Type)
}

func (v Value) ConvertToDouble() (float64, error) {
	switch v.Type {
	case CharValue:
		return float64(v.V.(int8)), nil
	case ShortValue:
		return float64(v.V.(int16)), nil
	case IntegerValue:
		return float64(v.V.(int32)), nil
	case FloatValue:
		return float64(v.V.(float32)), nil
	case DoubleValue:
		return v.V.(float64), nil
	}

	return 0, fmt.Errorf("cannot convert Type %d into float64", v.Type)
}

func (v Value) ConvertToChar() (int8, error) {
	if v.Type == CharValue {
		return v.V.(int8), nil
	}

	return 0, fmt.Errorf("cannot convert Type %d into int8", v.Type)
}
