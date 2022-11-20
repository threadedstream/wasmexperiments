package exec

import (
	"fmt"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
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

type doubleArgI[T any] struct {
	commonI
	arg0 T
	arg1 T
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

func newDoubleArgI[T any](op Op, arg0, arg1 T) Instr {
	switch op.Code {
	default:
		reporter.ReportError("unexpected double arg instruction with opcode %v", op.Code)
	case i32AddOp:
		return &I32AddI{doubleArgI[*uint32]{
			commonI: commonI{op: op},
			arg0:    arg0.(*uint32),
			arg1:    arg1.(*uint32),
		}}
	case i32SubOp:
		return &I32SubI{doubleArgI[*uint32]{
			commonI: commonI{op: op},
			arg0:    arg0.(*uint32),
			arg1:    arg1.(*uint32),
		}}
	case i32MulOp:
		return &I32MulI{doubleArgI[*uint32]{
			commonI: commonI{op: op},
			arg0:    arg0.(*uint32),
			arg1:    arg1.(*uint32),
		}}

	case i32DivUOp:
		return &I32DivUI{doubleArgI[*uint32]{
			commonI: commonI{op: op},
			arg0:    arg0.(*uint32),
			arg1:    arg1.(*uint32),
		}}

	case i32DivSOp:
		return &I32DivSI{doubleArgI[*uint32]{
			commonI: commonI{op: op},
			arg0:    arg0.(*uint32),
			arg1:    arg1.(*uint32),
		}}
	}
	panic("unreachable")
}

type singleArgI[T any] struct {
	commonI
	arg0 T
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

func newSingleArgI[T any](op Op, arg0 T) Instr {
	switch op.Code {

	}
	return nil
}

type (
	I32AddI struct {
		doubleArgI[*uint32]
	}

	I32SubI struct {
		doubleArgI[*uint32]
	}

	I32MulI struct {
		doubleArgI[*uint32]
	}

	I32DivUI struct {
		doubleArgI[*uint32]
	}

	I32DivSI struct {
		doubleArgI[*uint32]
	}

	I32EqI struct {
		singleArgI[*uint32]
	}

	I32LoadI struct {
		doubleArgI[uint32] // arg0 represents alignment, whereas arg1 represents offset
	}

	I32StoreI struct {
		doubleArgI[uint32]
	}

	I32ConstI struct {
		singleArgI[uint32]
	}

	I32LocalGetI struct {
		singleArgI[uint32]
	}

	I32GlobalGetI struct {
		singleArgI[uint32]
	}

	I32LocalSetI struct {
		doubleArgI[uint32]
	}

	I32GlobalSetI struct {
		doubleArgI[uint32]
	}

	IfI struct {
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
