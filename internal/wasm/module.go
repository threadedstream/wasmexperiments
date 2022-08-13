package wasm

import (
	"bytes"
	"errors"
	wr "github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
	"io"
)

const (
	magicCookie = 0x6d736100
	version     = 0x1
)

type Module struct {
	typesSection *TypesSection
	wr           *wr.WasmReader
}

func (m *Module) Read() error {
	// validate magic_cookie
	cookie, err := wbinary.ReadU32(m.wr)
	if err != nil {
		return err
	}
	if cookie != magicCookie {
		return errors.New("cookies do not match")
	}
	ver, err := wbinary.ReadU32(m.wr)
	if err != nil {
		return err
	}
	if ver != version {
		return errors.New("versions do not match")
	}

	if err := m.readSection(); err != nil {
		return err
	}
	return nil
}

func (m *Module) readSection() error {
	id, err := m.wr.ReadByte()
	if err != nil {
		return err
	}
	switch SectionID(id) {
	case TypeSection:
		if err := m.readTypeSection(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) readTypeSection() error {
	ts := &TypesSection{}
	dataLen, err := wbinary.ReadVarUint32(m.wr)
	if err != nil {
		return err
	}

	sectionData := new(bytes.Buffer)
	sectionData.Grow(int(dataLen))
	sectionReader := io.LimitReader(io.TeeReader(m.wr.Peek().(io.Reader), sectionData), int64(dataLen))
	m.wr.Push(sectionReader)
	_, err = m.wr.ReadBytes(int(dataLen))
	if err != nil {
		return err
	}

	// read arr length
	arrLen, err := wbinary.ReadVarUint32(m.wr)
	if err != nil {
		return err
	}
	ts.sigs = make([]*FunctionSig, arrLen, arrLen)

	for i, sig := range ts.sigs {
		sig = new(FunctionSig)
		// force value type to be a func
		valType, err := wbinary.ReadU8(m.wr)
		if err != nil {
			return err
		}
		if valType != ValueTypeFunc {
			return errors.New("readTypeSection: value type must be a function")
		}

		// start filling out function signatures slice
		paramsLen, err := wbinary.ReadVarUint32(m.wr)
		if err != nil {
			return err
		}

		sig.params = make([]ValueType, paramsLen, paramsLen)
		for i := 0; i < int(paramsLen); i++ {
			valTyp, err := wbinary.ReadVarUint32(m.wr)
			if err != nil {
				return err
			}
			sig.params[i] = ValueType(valTyp)
		}
		resultsLen, err := wbinary.ReadVarUint32(m.wr)
		if err != nil {
			return err
		}
		if resultsLen > 1 {
			return errors.New("readTypeSection: length of results array can't exceed the value of 1 (yet)")
		}
		for i := 0; i < int(resultsLen); i++ {
			valTyp, err := wbinary.ReadVarUint32(m.wr)
			if err != nil {
				return err
			}
			sig.results[i] = ValueType(valTyp)
		}
		ts.sigs[i] = sig
	}

	return nil
}
