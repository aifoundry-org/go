// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plan9obj

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

type fileTest struct {
	file     string
	hdr      FileHeader
	sections []*SectionHeader
}

var fileTests = []fileTest{
	{
		"testdata/386-plan9-exec",
		FileHeader{Magic386, 0x324, 0x14, 4, 0x1000, 32},
		[]*SectionHeader{
			{"text", 0x4c5f, 0x20},
			{"data", 0x94c, 0x4c7f},
			{"syms", 0x2c2b, 0x55cb},
			{"spsz", 0x0, 0x81f6},
			{"pcsz", 0xf7a, 0x81f6},
		},
	},
	{
		"testdata/amd64-plan9-exec",
		FileHeader{MagicAMD64, 0x618, 0x13, 8, 0x200000, 40},
		[]*SectionHeader{
			{"text", 0x4213, 0x28},
			{"data", 0xa80, 0x423b},
			{"syms", 0x2c8c, 0x4cbb},
			{"spsz", 0x0, 0x7947},
			{"pcsz", 0xca0, 0x7947},
		},
	},
}

func TestOpen(t *testing.T) {
	for i := range fileTests {
		tt := &fileTests[i]

		f, err := Open(tt.file)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(f.FileHeader, tt.hdr) {
			t.Errorf("open %s:\n\thave %#v\n\twant %#v\n", tt.file, f.FileHeader, tt.hdr)
			continue
		}

		for i, sh := range f.Sections {
			if i >= len(tt.sections) {
				break
			}
			have := &sh.SectionHeader
			want := tt.sections[i]
			if !reflect.DeepEqual(have, want) {
				t.Errorf("open %s, section %d:\n\thave %#v\n\twant %#v\n", tt.file, i, have, want)
			}
		}
		tn := len(tt.sections)
		fn := len(f.Sections)
		if tn != fn {
			t.Errorf("open %s: len(Sections) = %d, want %d", tt.file, fn, tn)
		}
	}
}

func TestOpenFailure(t *testing.T) {
	filename := "file.go"    // not a Plan 9 a.out file
	_, err := Open(filename) // don't crash
	if err == nil {
		t.Errorf("open %s: succeeded unexpectedly", filename)
	}
}

func TestRISCV64(t *testing.T) {
	const loadAddress = 0x1000000000
	const entry = loadAddress + 40
	text := []byte{0x13, 0, 0, 0, 0x67, 0x80, 0, 0}
	data := []byte{1, 2, 3, 4}
	var symbols bytes.Buffer
	binary.Write(&symbols, binary.BigEndian, uint64(entry))
	symbols.WriteByte('T' | 0x80)
	symbols.WriteString("main\x00")

	var image bytes.Buffer
	binary.Write(&image, binary.BigEndian, prog{
		Magic: MagicRISCV64,
		Text:  uint32(len(text)),
		Data:  uint32(len(data)),
		Bss:   16,
		Syms:  uint32(symbols.Len()),
		Entry: uint32(entry & 0xffffffff),
	})
	binary.Write(&image, binary.BigEndian, uint64(entry))
	image.Write(text)
	image.Write(data)
	image.Write(symbols.Bytes())

	f, err := NewFile(bytes.NewReader(image.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	wantHeader := FileHeader{MagicRISCV64, 16, entry, 8, loadAddress, 40}
	if f.FileHeader != wantHeader {
		t.Fatalf("header = %#v, want %#v", f.FileHeader, wantHeader)
	}
	wantSections := []SectionHeader{
		{"text", uint32(len(text)), 40},
		{"data", uint32(len(data)), 40 + uint32(len(text))},
		{"syms", uint32(symbols.Len()), 40 + uint32(len(text)+len(data))},
		{"spsz", 0, uint32(image.Len())},
		{"pcsz", 0, uint32(image.Len())},
	}
	for i, want := range wantSections {
		if got := f.Sections[i].SectionHeader; got != want {
			t.Errorf("section %d = %#v, want %#v", i, got, want)
		}
	}
	for name, want := range map[string][]byte{"text": text, "data": data} {
		got, err := f.Section(name).Data()
		if err != nil || !bytes.Equal(got, want) {
			t.Errorf("%s data = %x, %v; want %x, nil", name, got, err, want)
		}
	}
	syms, err := f.Symbols()
	wantSymbols := []Sym{{entry, 'T', "main"}}
	if err != nil || !reflect.DeepEqual(syms, wantSymbols) {
		t.Errorf("symbols = %#v, %v; want %#v, nil", syms, err, wantSymbols)
	}

	t.Run("truncated entry", func(t *testing.T) {
		if _, err := NewFile(bytes.NewReader(image.Bytes()[:39])); err == nil {
			t.Fatal("accepted truncated 64-bit entry")
		}
	})
	t.Run("missing expanded header flag", func(t *testing.T) {
		b := bytes.Clone(image.Bytes())
		binary.BigEndian.PutUint32(b, MagicRISCV64&^Magic64)
		if _, err := NewFile(bytes.NewReader(b)); err == nil {
			t.Fatal("accepted riscv64 magic without expanded header flag")
		}
	})
}
