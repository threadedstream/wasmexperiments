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

func newSingleArgI[T any](op Op, arg0 T) Instr {
	switch op.Code {

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

	I32EqI struct {
		singleArgI
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

	I32LocalGetI struct {
		singleArgI
	}

	I32GlobalGetI struct {
		singleArgI
	}

	I32LocalSetI struct {
		doubleArgI
	}

	I32GlobalSetI struct {
		doubleArgI
	}

	CallI struct {
		singleArgI
	}

	IfI struct {
		body     []Instr
		elseBody []Instr
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
