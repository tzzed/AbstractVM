package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BYTE_OBJ    = " INT8"
	SHORT_OBJ   = "INT16"
	FLOAT_OBJ   = "FLOAT32"
	DOUBLE_OBJ  = "DOUBLE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int32
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type Short struct {
	Value int16
}

func (s *Short) Inspect() string {
	return fmt.Sprintf("%d", s.Value)
}

func (s *Short) Type() ObjectType {
	return SHORT_OBJ
}

type Byte struct {
	Value int8
}

func (b *Byte) Inspect() string {
	return fmt.Sprintf("%d", b.Value)
}

func (b *Byte) Type() ObjectType {
	return BYTE_OBJ
}

type Float struct {
	Value float32
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
}

type Double struct {
	Value float64
}

func (d *Double) Inspect() string {
	return fmt.Sprintf("%f", d.Value)
}

func (d *Double) Type() ObjectType {
	return DOUBLE_OBJ
}
