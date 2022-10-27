package types

import "github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"

type ValueType uint8

const (
	ValueTypeI32  ValueType = 0x7f
	ValueTypeI64            = 0x7e
	ValueTypeF32            = 0x7d
	ValueTypeF64            = 0x7c
	ValueTypeFunc           = 0x60
)

var (
	vtmap = map[ValueType]string{
		ValueTypeI32:  "i32",
		ValueTypeI64:  "i64",
		ValueTypeF32:  "f32",
		ValueTypeF64:  "f64",
		ValueTypeFunc: "func",
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
