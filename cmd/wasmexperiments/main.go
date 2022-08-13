package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	magicCookie = 0x6d736100
	version     = 0x1
)

const (
	pageSize = 64 // 64k page size
)

func assert(cond bool, msg string) {
	if !cond {
		reportError(msg)
	}
}

type SectionID uint

const (
	CustomSection SectionID = iota
	TypeSection
	ImportSection
	FunctionSection
	TableSection
	LinearMemorySection
	GlobalSection
	ExportSection
	StartSection
	ElementSection
	CodeSection
	DataSection
)

type ValueType uint8

const (
	ValueTypeI32  ValueType = 0x7f
	ValueTypeFunc           = 0x60
)

func readWasm(path string) (bs []byte, err error) {
	bs, err = os.ReadFile(path)
	return
}

func byteArrEqual(p, q []byte) bool {
	if len(p) != len(q) {
		return false
	}
	n := len(p)
	for i := 0; i < n; i++ {
		if p[i] != q[i] {
			return false
		}
	}
	return true
}

func (m *Module) Read(reader io.Reader) error {
	// validate magic_cookie
	cookie, err := readU32(reader)
	if err != nil {
		return err
	}
	if cookie != magicCookie {
		return errors.New("cookies do not match")
	}
	ver, err := readU32(reader)
	if err != nil {
		return err
	}
	if ver != version {
		return errors.New("versions do not match")
	}

	if err := m.readSection(reader); err != nil {
		return err
	}
	return nil
}

func reportError(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func readU32(reader io.Reader) (uint32, error) {
	val, err := ReadBytes(reader, 4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(val), nil
}

func readU8(reader io.Reader) (uint8, error) {
	val, err := ReadByte(reader)
	if err != nil {
		return 0, err
	}
	return uint8(val), nil
}

func readVarUint(reader io.Reader, n int) (uint64, error) {
	if n > 64 {
		reportError("readVarUint: n can't be larger than 64")
	}
	var p byte
	var shift uint
	var res uint64
	for {
		p, _ = ReadByte(reader)
		b := uint64(p)
		switch {
		default:
			return 0, errors.New("readVarUint: invalid uint")
		case b < 1<<7 && b <= 1<<n-1:
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		}
	}
}

func readVarUint32(reader io.Reader) (uint32, error) {
	val, err := readVarUint(reader, 32)
	return uint32(val), err
}

func readVarUint64(reader io.Reader) (uint64, error) {
	return readVarUint(reader, 64)
}

type Module struct {
	typesSection *TypesSection
}

type Section interface {
	IsSection() bool
}

type FunctionSig struct {
	params  []ValueType
	results [1]ValueType // has at most one element (before rolling support for multiple return types)
}

type TypesSection struct {
	sigs   []*FunctionSig
	reader io.Reader
}

func (ts TypesSection) IsSection() bool { return true }

func (m *Module) readSection(reader io.Reader) error {
	id, err := ReadByte(reader)
	if err != nil {
		return err
	}
	switch SectionID(id) {
	case TypeSection:
		if err := m.readTypeSection(reader); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) readTypeSection(reader io.Reader) error {
	ts := &TypesSection{}
	dataLen, err := readVarUint32(reader)
	if err != nil {
		return err
	}

	sectionData := new(bytes.Buffer)
	sectionData.Grow(int(dataLen))
	sectionReader := io.LimitReader(io.TeeReader(reader, sectionData), int64(dataLen))
	_, err = ReadBytes(sectionReader, int(dataLen))
	if err != nil {
		return err
	}

	// read arr length
	arrLen, err := readVarUint32(sectionData)
	if err != nil {
		return err
	}
	ts.sigs = make([]*FunctionSig, arrLen, arrLen)

	for i, sig := range ts.sigs {
		sig = new(FunctionSig)
		// force value type to be a func
		valType, err := readU8(sectionData)
		if err != nil {
			return err
		}
		if valType != ValueTypeFunc {
			return errors.New("readTypeSection: value type must be a function")
		}

		// start filling out function signatures slice
		paramsLen, err := readVarUint32(sectionData)
		if err != nil {
			return err
		}

		sig.params = make([]ValueType, paramsLen, paramsLen)
		for i := 0; i < int(paramsLen); i++ {
			valTyp, err := readVarUint32(sectionData)
			if err != nil {
				return err
			}
			sig.params[i] = ValueType(valTyp)
		}
		resultsLen, err := readVarUint32(sectionData)
		if err != nil {
			return err
		}
		if resultsLen > 1 {
			return errors.New("readTypeSection: length of results array can't exceed the value of 1 (yet)")
		}
		for i := 0; i < int(resultsLen); i++ {
			valTyp, err := readVarUint32(sectionData)
			if err != nil {
				return err
			}
			sig.results[i] = ValueType(valTyp)
		}
		ts.sigs[i] = sig
	}

	return nil
}

func ReadByte(reader io.Reader) (byte, error) {
	bs, err := ReadBytes(reader, 1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func ReadBytes(reader io.Reader, n int) ([]byte, error) {
	bs := make([]byte, n, n)
	if _, err := reader.Read(bs); err != nil {
		return nil, err
	}
	return bs, nil
}

func main() {
	bs, err := readWasm("/Users/gildarov/toys/wasmexperiments/cmd/wasmexperiments/checkers.wasm")
	if err != nil {
		panic(err)
	}
	m := &Module{}
	reader := bytes.NewReader(bs)
	if err := m.Read(reader); err != nil {
		reportError(err.Error())
	}
}
