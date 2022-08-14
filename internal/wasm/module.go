package wasm

import (
	"errors"
	wr "github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
	"github.com/threadedstream/wasmexperiments/internal/pkg/werrors"
)

const (
	magicCookie = 0x6d736100
	version     = 0x1
)

type Module struct {
	typesSection *TypesSection
	wr           *wr.WasmReader
}

func NewModule(wr *wr.WasmReader) *Module {
	return &Module{
		wr: wr,
	}
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

	if err := m.readSections(); err != nil {
		return err
	}
	return nil
}

func (m *Module) readSections() error {
	// types section
	for id, handler := range map[SectionID]func() error{
		TypeSectionID: m.readTypeSection,
	} {
		if err := m.validateSectionID(id); err != nil {
			return err
		}
		if err := handler(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) validateSectionID(expected SectionID) error {
	id, err := m.wr.ReadByte()
	if err != nil {
		return err
	}

	if SectionID(id) != expected {
		return werrors.ErrInvalidSectionID
	}

	return nil
}

func (m *Module) readTypeSection() error {
	ts := &TypesSection{
		reader: m.wr,
	}
	if err := ts.Read(); err != nil {
		return err
	}
	m.typesSection = ts
	return nil
}
