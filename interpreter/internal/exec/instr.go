package exec

import (
	"fmt"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/types"
	"strings"
)

type Instr interface {
	Op() Op
	String() string
}

type commonI struct {
	op Op
}

func (ci commonI) Op() Op {
	return ci.op
}

func (ci commonI) String() string {
	return ""
}

type doubleArgI struct {
	commonI
	arg0, arg1 any
}

func (di doubleArgI) String() string {
	s := strings.Builder{}
	s.WriteString(di.Op().Name + " ")
	if di.arg0 == nil && di.arg1 == nil {
		return s.String()
	}
	s.WriteString(fmt.Sprint(di.arg0) + "," + fmt.Sprint(di.arg1))
	return s.String()
}

func newDoubleArgI(op Op, arg0, arg1 any) Instr {
	inner := doubleArgI{
		commonI: commonI{op: op},
		arg0:    arg0,
		arg1:    arg1,
	}
	switch op.Code {
	default:
		reporter.ReportError("unexpected double arg instruction with opcode %v", op.Code)
	case i32AddOp:
		return &I32AddI{inner}
	case i32SubOp:
		return &I32SubI{inner}
	case i32MulOp:
		return &I32MulI{inner}
	case i32DivUOp:
		return &I32DivUI{inner}
	case i32DivSOp:
		return &I32DivSI{inner}
	}
	panic("unreachable")
}

type singleArgI struct {
	commonI
	arg0 any
}

func (si singleArgI) String() string {
	s := strings.Builder{}
	s.WriteString(si.Op().Name + " ")
	if si.arg0 == nil {
		return s.String()
	}
	s.WriteString(fmt.Sprint(si.arg0))
	return s.String()
}

func newSingleArgI(op Op, arg0 any) Instr {
	inner := singleArgI{
		commonI: commonI{op: op},
		arg0:    arg0,
	}
	switch op.Code {
	default:
		reporter.ReportError("unexpected single arg instruction with opcode %v", op.Code)
	case i32ConstOp:
		return &I32ConstI{inner}
	case callOp:
		return &CallI{inner}
	case localGetOp:
		return &LocalGetI{inner}
	case globalGetOp:
		return &GlobalGetI{inner}
	case localSetOp:
		return &LocalSetI{inner}
	case brOp:
		return &BrI{inner, ""}
	case brIfOp:
		return &BrIfI{inner, ""}
	}
	panic("unreachable")
}

type noArgI struct {
	commonI
}

func (na noArgI) String() string {
	return na.Op().Name + " "
}

func newNoArgI(op Op) Instr {
	inner := noArgI{
		commonI: commonI{op},
	}
	switch op.Code {
	default:
		reporter.ReportError("unexpected no arg instruction with opcode %v", op.Code)
	case i32EqOp:
		return &I32EqI{inner}
	case endOp:
		return &EndI{inner}
	case returnOp:
		return &RetI{inner}
	case i32LtSOp:
		return &I32LtSI{inner}
	}
	panic("unreachable")
}

type blockTypedI struct {
	commonI
	blockType types.BlockType
}

func (i *blockTypedI) resolveBlockType(reader *wasm_reader.WasmReader) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	switch valty := types.ValueType(b); valty {
	default:
		i.blockType = types.OtherBlockType{X: int64(valty)}
	case types.ValueTypeEmpty:
		i.blockType = types.EmptyBlockType{}
	case types.ValueTypeI32, types.ValueTypeF32, types.ValueTypeI64, types.ValueTypeF64, types.ValueTypeVector, types.ValueTypeFuncRef, types.ValueTypeExternRef:
		i.blockType = types.ResultBlockType{Ty: valty}
	}
	return nil
}

type (
	I32AddI struct {
		doubleArgI
	}

	I32SubI struct {
		doubleArgI
	}

	I32MulI struct {
		doubleArgI
	}

	I32DivUI struct {
		doubleArgI
	}

	I32DivSI struct {
		doubleArgI
	}

	I32LoadI struct {
		doubleArgI
	}

	I32StoreI struct {
		doubleArgI
	}

	I32ConstI struct {
		singleArgI
	}

	LocalGetI struct {
		singleArgI
	}

	GlobalGetI struct {
		singleArgI
	}

	LocalSetI struct {
		singleArgI
	}

	GlobalSetI struct {
		singleArgI
	}

	CallI struct {
		singleArgI
	}

	I32EqI struct {
		noArgI
	}

	EndI struct {
		noArgI
	}

	RetI struct {
		noArgI
	}

	I32LtSI struct {
		noArgI
	}

	IfI struct {
		blockTypedI
		body     []Instr
		elseBody []Instr
	}

	BlockI struct {
		blockTypedI
		body []Instr
	}

	LoopI struct {
		blockTypedI
		body []Instr
	}

	BrI struct {
		singleArgI
		context string // useful when executing "br" inside "block" or "loop" constructs
	}

	BrIfI struct {
		singleArgI
		context string // useful when executing "br_if" inside "block" or "loop" constructs
	}
)

func Dump(is []Instr) {
	s := strings.Builder{}
	for _, i := range is {
		s.WriteString(i.String())
		s.WriteRune('\n')
	}
	println(s.String())
}
