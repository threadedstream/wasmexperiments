package types

import "github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"

type ValueType uint8

const (
	ValueTypeI32       ValueType = 0x7f
	ValueTypeI64                 = 0x7e
	ValueTypeF32                 = 0x7d
	ValueTypeF64                 = 0x7c
	ValueTypeFunc                = 0x60
	ValueTypeVector              = 0x7b
	ValueTypeFuncRef             = 0x70
	ValueTypeExternRef           = 0x6f
	ValueTypeEmpty               = 0x40
)

// some shortcuts
var (
	ValueTypeVoid      = []ValueType{0x0}
	ValueTypeSingleI32 = []ValueType{ValueTypeI32}
	ValueTypeDoubleI32 = []ValueType{ValueTypeI32, ValueTypeI32}
	ValueTypeSingleF32 = []ValueType{ValueTypeF32}
	ValueTypeDoubleF32 = []ValueType{ValueTypeF32, ValueTypeF32}
)

var (
	vtmap = map[ValueType]string{
		ValueTypeI32:       "i32",
		ValueTypeI64:       "i64",
		ValueTypeF32:       "f32",
		ValueTypeF64:       "f64",
		ValueTypeFunc:      "func",
		ValueTypeVector:    "vector",
		ValueTypeFuncRef:   "func_ref",
		ValueTypeExternRef: "extern_ref",
	}
)

func (v ValueType) String() string { return vtmap[v] }

func (_ ValueType) Serialize() error { return nil }

func (v *ValueType) Deserialize(reader *wasm_reader.WasmReader) error {
	val, err := reader.ReadByte()
	if err != nil {
		return err
	}
	*v = ValueType(val)
	return nil
}

type BlockType interface {
	Empty() ValueType
	ValType() ValueType
	Other() int64 // 33-bit signed integer
}

type EmptyBlockType struct{}

func (e EmptyBlockType) Empty() ValueType {
	return ValueTypeEmpty
}

func (e EmptyBlockType) ValType() ValueType {
	//TODO implement me
	panic("implement me")
}

func (e EmptyBlockType) Other() int64 {
	//TODO implement me
	panic("implement me")
}

type ResultBlockType struct {
	Ty ValueType
}

func (r ResultBlockType) Empty() ValueType {
	//TODO implement me
	panic("implement me")
}

func (r ResultBlockType) ValType() ValueType {
	return r.Ty
}

func (r ResultBlockType) Other() int64 {
	//TODO implement me
	panic("implement me")
}

type OtherBlockType struct {
	X int64
}

func (o OtherBlockType) Empty() ValueType {
	//TODO implement me
	panic("implement me")
}

func (o OtherBlockType) ValType() ValueType {
	//TODO implement me
	panic("implement me")
}

func (o OtherBlockType) Other() int64 {
	return o.X
}
