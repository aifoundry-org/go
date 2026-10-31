// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riscv64

import (
	"cmd/internal/objabi"
	"cmd/internal/sys"
	"cmd/link/internal/ld"
)

func Init() (*sys.Arch, ld.Arch) {
	arch := sys.ArchRISCV64

	theArch := ld.Arch{
		Funcalign:  funcAlign,
		Maxalign:   maxAlign,
		Minalign:   minAlign,
		Dwarfregsp: dwarfRegSP,
		Dwarfreglr: dwarfRegLR,

		// B_MAGIC, _MAGIC(HDR_MAGIC, 30) in 9front's a.out.h.
		Plan9Magic:  uint32(4*30*30 + 7),
		Plan9_64Bit: true,

		Adddynrel:        adddynrel,
		Archinit:         archinit,
		Archreloc:        archreloc,
		Archrelocvariant: archrelocvariant,
		Extreloc:         extreloc,

		// TrampLimit is set such that we always run the trampoline
		// generation code. This is necessary since calls to external
		// symbols require the use of trampolines, regardless of the
		// text size.
		TrampLimit: 1,
		Trampoline: trampoline,

		Gentext:     gentext,
		GenSymsLate: genSymsLate,
		Machoreloc1: machoreloc1,

		ELF: ld.ELFArch{
			Linuxdynld:     "/lib/ld-linux-riscv64-lp64d.so.1",
			LinuxdynldMusl: "/lib/ld-musl-riscv64.so.1",

			Freebsddynld:   "/usr/libexec/ld-elf.so.1",
			Netbsddynld:    "XXX",
			Openbsddynld:   "/usr/libexec/ld.so",
			Dragonflydynld: "XXX",
			Solarisdynld:   "XXX",

			Reloc1:    elfreloc1,
			RelocSize: 24,
			SetupPLT:  elfsetupplt,
		},
	}

	return arch, theArch
}

func archinit(ctxt *ld.Link) {
	switch ctxt.HeadType {
	case objabi.Hplan9: /* plan 9 */
		// a.out with the 64-bit header: text, header included,
		// at the kernel's UTZERO; data on the next page.
		ld.HEADR = 32 + 8
		if *ld.FlagRound == -1 {
			*ld.FlagRound = 0x1000
		}
		if *ld.FlagTextAddr == -1 {
			*ld.FlagTextAddr = ld.Rnd(0x1000000000, *ld.FlagRound) + int64(ld.HEADR)
		}

	case objabi.Hlinux, objabi.Hfreebsd, objabi.Hopenbsd:
		ld.Elfinit(ctxt)
		ld.HEADR = ld.ELFRESERVE
		if *ld.FlagRound == -1 {
			*ld.FlagRound = 0x10000
		}
		if *ld.FlagTextAddr == -1 {
			*ld.FlagTextAddr = ld.Rnd(0x10000, *ld.FlagRound) + int64(ld.HEADR)
		}
	default:
		ld.Exitf("unknown -H option: %v", ctxt.HeadType)
	}
}
