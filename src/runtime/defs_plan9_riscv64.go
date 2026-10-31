// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const _PAGESIZE = 0x1000

// ureg is the riscv64 Plan 9 kernel's Ureg (/riscv64/include/ureg.h):
// the 31 integer registers with the pc in place of the zero register,
// then supervisor CSRs.
type ureg struct {
	pc      uint64
	ra      uint64 // r1, link
	sp      uint64 // r2
	gp      uint64 // r3
	tp      uint64 // r4
	t0      uint64 // r5
	t1      uint64 // r6
	t2      uint64 // r7
	s0      uint64 // r8: system call number and return value
	s1      uint64 // r9
	a0      uint64 // r10
	a1      uint64
	a2      uint64
	a3      uint64
	a4      uint64
	a5      uint64
	a6      uint64
	a7      uint64 // r17
	s2      uint64 // r18
	s3      uint64
	s4      uint64
	s5      uint64
	s6      uint64
	s7      uint64
	s8      uint64
	s9      uint64
	s10     uint64
	s11     uint64 // r27: g
	t3      uint64 // r28
	t4      uint64
	t5      uint64
	t6      uint64 // r31
	status  uint64
	ie      uint64
	cause   uint64
	tval    uint64 // faulting address
	curmode uint64
	scratch uint64
}

type sigctxt struct {
	u *ureg
}

//go:nosplit
//go:nowritebarrierrec
func (c *sigctxt) pc() uintptr { return uintptr(c.u.pc) }

func (c *sigctxt) sp() uintptr { return uintptr(c.u.sp) }
func (c *sigctxt) lr() uintptr { return uintptr(c.u.ra) }

func (c *sigctxt) setpc(x uintptr) { c.u.pc = uint64(x) }
func (c *sigctxt) setsp(x uintptr) { c.u.sp = uint64(x) }
func (c *sigctxt) setlr(x uintptr) { c.u.ra = uint64(x) }

// savelr stashes the link register in T0, from where sigpanictramp
// pushes it onto the stack before entering sigpanic.
func (c *sigctxt) savelr(x uintptr) { c.u.t0 = uint64(x) }

func dumpregs(u *ureg) {
	print("pc      ", hex(u.pc), "\n")
	print("ra      ", hex(u.ra), "\n")
	print("sp      ", hex(u.sp), "\n")
	print("gp      ", hex(u.gp), "\n")
	print("tp      ", hex(u.tp), "\n")
	print("t0      ", hex(u.t0), "\n")
	print("t1      ", hex(u.t1), "\n")
	print("t2      ", hex(u.t2), "\n")
	print("s0      ", hex(u.s0), "\n")
	print("s1      ", hex(u.s1), "\n")
	print("a0      ", hex(u.a0), "\n")
	print("a1      ", hex(u.a1), "\n")
	print("a2      ", hex(u.a2), "\n")
	print("a3      ", hex(u.a3), "\n")
	print("a4      ", hex(u.a4), "\n")
	print("a5      ", hex(u.a5), "\n")
	print("a6      ", hex(u.a6), "\n")
	print("a7      ", hex(u.a7), "\n")
	print("s2      ", hex(u.s2), "\n")
	print("s3      ", hex(u.s3), "\n")
	print("s4      ", hex(u.s4), "\n")
	print("s5      ", hex(u.s5), "\n")
	print("s6      ", hex(u.s6), "\n")
	print("s7      ", hex(u.s7), "\n")
	print("s8      ", hex(u.s8), "\n")
	print("s9      ", hex(u.s9), "\n")
	print("s10     ", hex(u.s10), "\n")
	print("g       ", hex(u.s11), "\n")
	print("t3      ", hex(u.t3), "\n")
	print("t4      ", hex(u.t4), "\n")
	print("t5      ", hex(u.t5), "\n")
	print("t6      ", hex(u.t6), "\n")
	print("status  ", hex(u.status), "\n")
	print("cause   ", hex(u.cause), "\n")
	print("tval    ", hex(u.tval), "\n")
}

func sigpanictramp()
