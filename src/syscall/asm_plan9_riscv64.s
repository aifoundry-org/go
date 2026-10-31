// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"

#define SYS_ERRSTR 41	/* from zsysnum_plan9.go */
#define SYS_SEEK 39	/* from zsysnum_plan9.go */

// System call support for Plan 9 on riscv64: the call number in X8,
// the arguments in 8-byte slots from 8(SP), ECALL, the result in X8.
// See ../runtime/sys_plan9_riscv64.s.

// Frame of Syscall and Syscall6: 8(X2) up to 56(X2) system call
// arguments, 56(X2) up to 184(X2) error buffer.

//func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err ErrorString)
TEXT	·Syscall(SB),NOSPLIT,$184-64
	NO_LOCAL_POINTERS
	CALL	runtime·entersyscall<ABIInternal>(SB)
	MOV	a1+8(FP), T0
	MOV	T0, 8(X2)
	MOV	a2+16(FP), T0
	MOV	T0, 16(X2)
	MOV	a3+24(FP), T0
	MOV	T0, 24(X2)
	MOV	trap+0(FP), X8
	ECALL
	MOV	X8, r1+32(FP)
	MOV	ZERO, r2+40(FP)
	MOVW	X8, T0
	MOV	$-1, T1
	BNE	T0, T1, ok3

	ADD	$56, X2, T0
	MOV	T0, 8(X2)
	MOV	$128, T0
	MOV	T0, 16(X2)
	MOV	$SYS_ERRSTR, X8
	ECALL
	CALL	runtime·exitsyscall(SB)	// call via ABI wrapper
	ADD	$56, X2, T0
	MOV	T0, 8(X2)
	CALL	runtime·gostring(SB)
	MOV	16(X2), T0
	MOV	24(X2), T1
	JMP	copyresult3

ok3:
	CALL	runtime·exitsyscall(SB)	// call via ABI wrapper
	MOV	$·emptystring(SB), T2
	MOV	0(T2), T0
	MOV	8(T2), T1

copyresult3:
	MOV	T0, err_base+48(FP)
	MOV	T1, err_len+56(FP)
	RET

//func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err ErrorString)
TEXT	·Syscall6(SB),NOSPLIT,$184-88
	NO_LOCAL_POINTERS
	CALL	runtime·entersyscall<ABIInternal>(SB)
	MOV	a1+8(FP), T0
	MOV	T0, 8(X2)
	MOV	a2+16(FP), T0
	MOV	T0, 16(X2)
	MOV	a3+24(FP), T0
	MOV	T0, 24(X2)
	MOV	a4+32(FP), T0
	MOV	T0, 32(X2)
	MOV	a5+40(FP), T0
	MOV	T0, 40(X2)
	MOV	a6+48(FP), T0
	MOV	T0, 48(X2)
	MOV	trap+0(FP), X8
	ECALL
	MOV	X8, r1+56(FP)
	MOV	ZERO, r2+64(FP)
	MOVW	X8, T0
	MOV	$-1, T1
	BNE	T0, T1, ok4

	ADD	$56, X2, T0
	MOV	T0, 8(X2)
	MOV	$128, T0
	MOV	T0, 16(X2)
	MOV	$SYS_ERRSTR, X8
	ECALL
	CALL	runtime·exitsyscall(SB)	// call via ABI wrapper
	ADD	$56, X2, T0
	MOV	T0, 8(X2)
	CALL	runtime·gostring(SB)
	MOV	16(X2), T0
	MOV	24(X2), T1
	JMP	copyresult4

ok4:
	CALL	runtime·exitsyscall(SB)	// call via ABI wrapper
	MOV	$·emptystring(SB), T2
	MOV	0(T2), T0
	MOV	8(T2), T1

copyresult4:
	MOV	T0, err_base+72(FP)
	MOV	T1, err_len+80(FP)
	RET

// RawSyscall and RawSyscall6 are NOFRAME: slide the arguments down on
// top of the call number, so that they start at 0(FP) = 8(X2).

//func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr)
TEXT ·RawSyscall(SB),NOSPLIT|NOFRAME,$0-56
	MOV	trap+0(FP), X8
	MOV	a1+8(FP), T0
	MOV	T0, trap+0(FP)
	MOV	a2+16(FP), T0
	MOV	T0, a1+8(FP)
	MOV	a3+24(FP), T0
	MOV	T0, a2+16(FP)
	ECALL
	MOV	X8, r1+32(FP)
	MOV	X8, r2+40(FP)
	MOV	X8, err+48(FP)
	RET

//func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr)
TEXT	·RawSyscall6(SB),NOSPLIT|NOFRAME,$0-80
	MOV	trap+0(FP), X8
	MOV	a1+8(FP), T0
	MOV	T0, trap+0(FP)
	MOV	a2+16(FP), T0
	MOV	T0, a1+8(FP)
	MOV	a3+24(FP), T0
	MOV	T0, a2+16(FP)
	MOV	a4+32(FP), T0
	MOV	T0, a3+24(FP)
	MOV	a5+40(FP), T0
	MOV	T0, a4+32(FP)
	MOV	a6+48(FP), T0
	MOV	T0, a5+40(FP)
	ECALL
	MOV	X8, r1+56(FP)
	MOV	X8, r2+64(FP)
	MOV	X8, err+72(FP)
	RET

// The kernel's seek is _seek(vlong *ret, int fd, vlong off, int type).
//func seek(placeholder uintptr, fd int, offset int64, whence int) (newoffset int64, err string)
TEXT ·seek(SB),NOSPLIT,$40-56
	NO_LOCAL_POINTERS
	MOV	$newoffset+32(FP), T0
	MOV	T0, 8(X2)
	MOV	fd+8(FP), T0
	MOV	T0, 16(X2)
	MOV	offset+16(FP), T0
	MOV	T0, 24(X2)
	MOV	whence+24(FP), T0
	MOV	T0, 32(X2)
	MOV	$SYS_SEEK, X8
	ECALL
	MOVW	X8, T0
	MOV	$-1, T1
	BNE	T0, T1, ok6
	MOV	T1, newoffset+32(FP)

	CALL	syscall·errstr(SB)
	MOV	8(X2), T0
	MOV	16(X2), T1
	JMP	copyresult6

ok6:
	MOV	$·emptystring(SB), T2
	MOV	0(T2), T0
	MOV	8(T2), T1

copyresult6:
	MOV	T0, err_base+40(FP)
	MOV	T1, err_len+48(FP)
	RET
