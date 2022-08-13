package wasm

type ValueType uint8

const (
	ValueTypeI32  ValueType = 0x7f
	ValueTypeFunc           = 0x60
)

var (
	vtmap = map[ValueType]string{
		ValueTypeI32:  "i32",
		ValueTypeFunc: "func",
	}
)

func (v ValueType) String() string {
	return vtmap[v]
}
