// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

import (
	"cmd/internal/objabi"
	"cmd/internal/sys"
	"cmd/link/internal/loader"
	"cmd/link/internal/sym"
	"fmt"
	"internal/buildcfg"
	"testing"
)

func TestPlan9RISCV64TextAlignment(t *testing.T) {
	// textaddress uses linker-wide state, so this test cannot run in parallel.
	oldArch, oldText, oldAlign, oldThearch := buildcfg.GOARCH, Segtext, Funcalign, thearch
	oldTextAddr, oldRandLayout := *FlagTextAddr, *flagRandLayout
	defer func() {
		buildcfg.GOARCH, Segtext, Funcalign, thearch = oldArch, oldText, oldAlign, oldThearch
		*FlagTextAddr, *flagRandLayout = oldTextAddr, oldRandLayout
	}()
	buildcfg.GOARCH = "riscv64"
	Funcalign = 8
	thearch = Arch{}
	*FlagTextAddr = 0x1000000028
	*flagRandLayout = 0

	for _, align := range []int32{8, 16, 32} {
		t.Run(fmt.Sprintf("align%d", align), func(t *testing.T) {
			Segtext = sym.Segment{}
			ctxt := setUpContext(sys.ArchRISCV64, false, objabi.Hplan9, "exe", "internal")
			ctxt.Loaded = true
			ldr := ctxt.loader
			for _, name := range []string{"first", "aligned"} {
				s := ldr.CreateSymForUpdate(name, 0)
				s.SetType(sym.STEXT)
				s.SetSize(8)
				if name == "aligned" {
					s.SetAlign(align)
				}
				ctxt.Textp = append(ctxt.Textp, s.Sym())
			}
			ctxt.textaddress()
			sect := Segtext.Sections[0]
			text := ldr.SymValue(ldr.Lookup("runtime.text", 0))
			// address will round the section base to sect.Align. It must
			// still equal runtime.text, or R_METHODOFF and pclntab offsets
			// will be relative to different bases.
			want := Rnd(*FlagTextAddr, int64(sect.Align))
			if text != want || int64(sect.Vaddr) != want {
				t.Fatalf("runtime.text = %#x, section = %#x; want %#x", text, sect.Vaddr, want)
			}
			if sect.Align != align {
				t.Errorf("section alignment = %d, want %d", sect.Align, align)
			}
		})
	}
}

func setUpContext(arch *sys.Arch, iself bool, ht objabi.HeadType, bm, lm string) *Link {
	ctxt := linknew(arch)
	ctxt.HeadType = ht
	er := loader.ErrorReporter{}
	ctxt.loader = loader.NewLoader(0, &er)
	ctxt.BuildMode.Set(bm)
	ctxt.LinkMode.Set(lm)
	ctxt.IsELF = iself
	ctxt.mustSetHeadType()
	ctxt.setArchSyms()
	return ctxt
}

// Make sure the addgotsym properly increases the symbols.
func TestAddGotSym(t *testing.T) {
	tests := []struct {
		arch    *sys.Arch
		ht      objabi.HeadType
		bm, lm  string
		rel     string
		relsize int
		gotsize int
	}{
		{
			arch:    sys.Arch386,
			ht:      objabi.Hlinux,
			bm:      "pie",
			lm:      "internal",
			rel:     ".rel",
			relsize: 2 * sys.Arch386.PtrSize,
			gotsize: sys.Arch386.PtrSize,
		},
		{
			arch:    sys.ArchAMD64,
			ht:      objabi.Hlinux,
			bm:      "pie",
			lm:      "internal",
			rel:     ".rela",
			relsize: 3 * sys.ArchAMD64.PtrSize,
			gotsize: sys.ArchAMD64.PtrSize,
		},
		{
			arch:    sys.ArchAMD64,
			ht:      objabi.Hdarwin,
			bm:      "pie",
			lm:      "external",
			gotsize: sys.ArchAMD64.PtrSize,
		},
	}

	// Save the architecture as we're going to set it on each test run.
	origArch := buildcfg.GOARCH
	defer func() {
		buildcfg.GOARCH = origArch
	}()

	for i, test := range tests {
		iself := len(test.rel) != 0
		buildcfg.GOARCH = test.arch.Name
		ctxt := setUpContext(test.arch, iself, test.ht, test.bm, test.lm)
		foo := ctxt.loader.CreateSymForUpdate("foo", 0)
		ctxt.loader.CreateExtSym("bar", 0)
		AddGotSym(&ctxt.Target, ctxt.loader, &ctxt.ArchSyms, foo.Sym(), 0)

		if iself {
			rel := ctxt.loader.Lookup(test.rel, 0)
			if rel == 0 {
				t.Fatalf("[%d] could not find symbol: %q", i, test.rel)
			}
			if s := ctxt.loader.SymSize(rel); s != int64(test.relsize) {
				t.Fatalf("[%d] expected ldr.Size(%q) == %v, got %v", i, test.rel, test.relsize, s)
			}
		}
		if s := ctxt.loader.SymSize(ctxt.loader.Lookup(".got", 0)); s != int64(test.gotsize) {
			t.Fatalf(`[%d] expected ldr.Size(".got") == %v, got %v`, i, test.gotsize, s)
		}
	}
}

func TestWriteULebFixedLength(t *testing.T) {
	flavs := []objabi.RelocType{
		objabi.R_DWTXTADDR_U1,
		objabi.R_DWTXTADDR_U2,
		objabi.R_DWTXTADDR_U3,
		objabi.R_DWTXTADDR_U4,
	}
	var clear, scratch [7]byte
	tmp := scratch[:]
	for i := range 5 {
		for _, rt := range flavs {
			scratch = clear
			_, leb128len := rt.DwTxtAddrRelocParams()
			_, n := objabi.FuncCountToDwTxtAddrFlavor(i)
			if n > leb128len {
				continue
			}
			err := writeUleb128FixedLength(tmp, uint64(i), leb128len)
			if err != nil {
				t.Errorf("unexpected err %v on val %d flav %s leb128len %d",
					err, i, rt.String(), leb128len)
				continue
			}
		}
	}
}
