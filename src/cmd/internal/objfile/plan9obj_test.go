// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objfile

import (
	"bytes"
	"debug/plan9obj"
	"encoding/binary"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPlan9RISCV64(t *testing.T) {
	const loadAddress = 0x1000000000
	const textStart = loadAddress + 40
	instructions := []byte{0x13, 0, 0, 0, 0x67, 0x80, 0, 0}
	pclntab := []byte{1, 2, 3, 4}
	text := append(bytes.Clone(instructions), pclntab...)
	wantSymbols := []Sym{
		{Name: "runtime.text", Addr: textStart, Size: 8, Code: 'T'},
		{Name: "runtime.pclntab", Addr: textStart + 8, Size: 4, Code: 'T'},
		{Name: "runtime.epclntab", Addr: textStart + 12, Code: 'T'},
	}
	var symbols bytes.Buffer
	for _, s := range wantSymbols {
		binary.Write(&symbols, binary.BigEndian, s.Addr)
		symbols.WriteByte(byte(s.Code) | 0x80)
		symbols.WriteString(s.Name)
		symbols.WriteByte(0)
	}

	var image bytes.Buffer
	binary.Write(&image, binary.BigEndian, [8]uint32{
		plan9obj.MagicRISCV64, uint32(len(text)), 0, 0,
		uint32(symbols.Len()), uint32(textStart & 0xffffffff), 0, 0,
	})
	binary.Write(&image, binary.BigEndian, uint64(textStart))
	image.Write(text)
	image.Write(symbols.Bytes())
	name := filepath.Join(t.TempDir(), "riscv64-plan9-exec")
	if err := os.WriteFile(name, image.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	f, err := Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if got := f.GOARCH(); got != "riscv64" {
		t.Errorf("GOARCH = %q, want riscv64", got)
	}
	if got, err := f.LoadAddress(); err != nil || got != loadAddress {
		t.Errorf("LoadAddress = %#x, %v; want %#x, nil", got, err, loadAddress)
	}
	if start, got, err := f.Text(); err != nil || start != textStart || !bytes.Equal(got, text) {
		t.Errorf("Text = %#x, %x, %v; want %#x, %x, nil", start, got, err, textStart, text)
	}
	if got, err := f.Symbols(); err != nil || !reflect.DeepEqual(got, wantSymbols) {
		t.Errorf("Symbols = %#v, %v; want %#v, nil", got, err, wantSymbols)
	}
	if start, got, err := f.entries[0].raw.pcln(); err != nil || start != textStart || !bytes.Equal(got, pclntab) {
		t.Errorf("pcln = %#x, %x, %v; want %#x, %x, nil", start, got, err, textStart, pclntab)
	}
}
