// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

// System calls on riscv64 Plan 9: the call number goes in R8 (X8), the
// arguments are on the stack from 8(SP) as the Plan 9 C compiler (jc)
// lays them out (4-byte types in 4-byte slots, 8-byte types 8-aligned,
// which is also Go's ABI0 layout for the same types), ECALL, and the
// result comes back in R8.  The kernel preserves all other registers.
// In a NOFRAME function the arguments are already at 8(X2) = 0(FP).

// from ../syscall/zsysnum_plan9.go

#define SYS_SYSR1       0
#define SYS_BIND        2
#define SYS_CHDIR       3
#define SYS_CLOSE       4
#define SYS_DUP         5
#define SYS_ALARM       6
#define SYS_EXEC        7
#define SYS_EXITS       8
#define SYS_FAUTH       10
#define SYS_SEGBRK      12
#define SYS_OPEN        14
#define SYS_OSEEK       16
#define SYS_SLEEP       17
#define SYS_RFORK       19
#define SYS_PIPE        21
#define SYS_CREATE      22
#define SYS_FD2PATH     23
#define SYS_BRK_        24
#define SYS_REMOVE      25
#define SYS_NOTIFY      28
#define SYS_NOTED       29
#define SYS_SEGATTACH   30
#define SYS_SEGDETACH   31
#define SYS_SEGFREE     32
#define SYS_SEGFLUSH    33
#define SYS_RENDEZVOUS  34
#define SYS_UNMOUNT     35
#define SYS_SEMACQUIRE  37
#define SYS_SEMRELEASE  38
#define SYS_SEEK        39
#define SYS_FVERSION    40
#define SYS_ERRSTR      41
#define SYS_STAT        42
#define SYS_FSTAT       43
#define SYS_WSTAT       44
#define SYS_FWSTAT      45
#define SYS_MOUNT       46
#define SYS_AWAIT       47
#define SYS_PREAD       50
#define SYS_PWRITE      51
#define SYS_TSEMACQUIRE 52
#define SYS_NSEC        53

//func open(name *byte, mode, perm int32) int32
TEXT runtime·open(SB),NOSPLIT|NOFRAME,$0-20
	MOV	$SYS_OPEN, X8
	ECALL
	MOVW	X8, ret+16(FP)
	RET

//func pread(fd int32, buf unsafe.Pointer, nbytes int32, offset int64) int32
TEXT runtime·pread(SB),NOSPLIT|NOFRAME,$0-36
	MOV	$SYS_PREAD, X8
	ECALL
	MOVW	X8, ret+32(FP)
	RET

//func pwrite(fd int32, buf unsafe.Pointer, nbytes int32, offset int64) int32
TEXT runtime·pwrite(SB),NOSPLIT|NOFRAME,$0-36
	MOV	$SYS_PWRITE, X8
	ECALL
	MOVW	X8, ret+32(FP)
	RET

// The kernel's seek is _seek(vlong *ret, int fd, vlong off, int type).
//func seek(fd int32, offset int64, whence int32) int64
TEXT runtime·seek(SB),NOSPLIT,$32-32
	MOV	$ret+24(FP), T0
	MOV	T0, 8(X2)
	MOVW	fd+0(FP), T0
	MOV	T0, 16(X2)
	MOV	offset+8(FP), T0
	MOV	T0, 24(X2)
	MOVW	whence+16(FP), T0
	MOV	T0, 32(X2)
	MOV	$SYS_SEEK, X8
	ECALL
	MOVW	X8, T0
	MOV	$-1, T1
	BNE	T0, T1, seekok
	MOV	T1, ret+24(FP)
seekok:
	RET

//func closefd(fd int32) int32
TEXT runtime·closefd(SB),NOSPLIT|NOFRAME,$0-12
	MOV	$SYS_CLOSE, X8
	ECALL
	MOVW	X8, ret+8(FP)
	RET

//func dupfd(old, new int32) int32
TEXT runtime·dupfd(SB),NOSPLIT|NOFRAME,$0-12
	MOV	$SYS_DUP, X8
	ECALL
	MOVW	X8, ret+8(FP)
	RET

//func exits(msg *byte)
TEXT runtime·exits(SB),NOSPLIT|NOFRAME,$0-8
	MOV	$SYS_EXITS, X8
	ECALL
	RET

//func brk_(addr unsafe.Pointer) int32
TEXT runtime·brk_(SB),NOSPLIT|NOFRAME,$0-12
	MOV	$SYS_BRK_, X8
	ECALL
	MOVW	X8, ret+8(FP)
	RET

//func sleep(ms int32) int32
TEXT runtime·sleep(SB),NOSPLIT|NOFRAME,$0-12
	MOV	$SYS_SLEEP, X8
	ECALL
	MOVW	X8, ret+8(FP)
	RET

//func plan9_semacquire(addr *uint32, block int32) int32
TEXT runtime·plan9_semacquire(SB),NOSPLIT|NOFRAME,$0-20
	MOV	$SYS_SEMACQUIRE, X8
	ECALL
	MOVW	X8, ret+16(FP)
	RET

//func plan9_tsemacquire(addr *uint32, ms int32) int32
TEXT runtime·plan9_tsemacquire(SB),NOSPLIT|NOFRAME,$0-20
	MOV	$SYS_TSEMACQUIRE, X8
	ECALL
	MOVW	X8, ret+16(FP)
	RET

//func plan9_semrelease(addr *uint32, count int32) int32
TEXT runtime·plan9_semrelease(SB),NOSPLIT|NOFRAME,$0-20
	MOV	$SYS_SEMRELEASE, X8
	ECALL
	MOVW	X8, ret+16(FP)
	RET

//func notify(fn unsafe.Pointer) int32
TEXT runtime·notify(SB),NOSPLIT|NOFRAME,$0-12
	MOV	$SYS_NOTIFY, X8
	ECALL
	MOVW	X8, ret+8(FP)
	RET

//func noted(mode int32) int32
TEXT runtime·noted(SB),NOSPLIT|NOFRAME,$0-12
	MOV	$SYS_NOTED, X8
	ECALL
	MOVW	X8, ret+8(FP)
	RET

//func rfork(flags int32) int32
TEXT runtime·rfork(SB),NOSPLIT|NOFRAME,$0-12
	MOV	$SYS_RFORK, X8
	ECALL
	MOVW	X8, ret+8(FP)
	RET

// func timesplit(u uint64) (sec int64, nsec int32)
TEXT runtime·timesplit(SB),NOSPLIT|NOFRAME,$0-20
	MOV	u+0(FP), T0
	MOV	$1000000000, T1
	DIVU	T1, T0, T2
	REMU	T1, T0, T3
	MOV	T2, sec+8(FP)
	MOVW	T3, nsec+16(FP)
	RET

//func tstart_plan9(newm *m)
TEXT runtime·tstart_plan9(SB),NOSPLIT,$8-8
	MOV	newm+0(FP), A0
	MOV	m_g0(A0), g

	// Layout new m scheduler stack on os stack.
	MOV	X2, T0
	MOV	T0, (g_stack+stack_hi)(g)
	MOV	$(64*1024), T1
	SUB	T1, T0
	MOV	T0, (g_stack+stack_lo)(g)
	MOV	T0, g_stackguard0(g)
	MOV	T0, g_stackguard1(g)

	// Initialize procid from TOS struct.
	MOV	_tos(SB), T0
	MOVWU	64(T0), T0
	MOV	T0, m_procid(A0)	// save pid as m->procid

	CALL	runtime·mstart(SB)

	// Exit the thread.
	MOV	ZERO, 8(X2)
	CALL	runtime·exits(SB)
	JMP	0(PC)

// The kernel calls the note handler on the current stack with the
// saved Ureg's address at 8(SP) and the note's at 16(SP), g (X27) as
// it was; the handler continues the process with noted(2).
//func sigtramp(ureg, note unsafe.Pointer)
TEXT runtime·sigtramp(SB),NOSPLIT|NOFRAME,$0-16
	// check that g and m exist
	BEQ	g, ZERO, nom
	MOV	g_m(g), T0
	BNE	T0, ZERO, havem
nom:
	CALL	runtime·badsignal2(SB)	// will exit

havem:
	// save args
	MOV	ureg+0(FP), A1
	MOV	note+8(FP), A2

	// change stack
	MOV	m_gsignal(T0), T1
	MOV	(g_stack+stack_hi)(T1), X2

	// make room for args, retval and g
	// 8: ureg, 16: note, 24: gp, 32: retval, 40: saved g
	ADD	$-48, X2

	MOV	A1, 8(X2)
	MOV	A2, 16(X2)
	MOV	g, 24(X2)
	MOV	g, 40(X2)

	// g = m->gsignal
	MOV	T1, g

	CALL	runtime·sighandler(SB)
	MOV	32(X2), T0			// retval

	// restore g
	MOV	40(X2), g

	// call noted(retval)
	MOV	T0, 8(X2)
	CALL	runtime·noted(SB)
	RET

// sighandler set ra to the faulting pc and stashed the old ra in T0:
// push it as if the faulting function had saved it, then enter
// sigpanic as if called from the faulting pc.
//func sigpanictramp()
TEXT runtime·sigpanictramp(SB),NOSPLIT|NOFRAME,$0-0
	ADD	$-8, X2
	MOV	T0, 0(X2)
	JMP	runtime·sigpanic<ABIInternal>(SB)

//func setfpmasks()
// Only used by the 64-bit runtime.  RISC-V has no FP exception traps.
TEXT runtime·setfpmasks(SB),NOSPLIT|NOFRAME,$0
	RET

#define ERRMAX 128	/* from os_plan9.h */

// func errstr() string
// Only used by package syscall.
// Grab error string due to a syscall made
// in entersyscall mode, without going
// through the allocator (issue 4994).
// See ../syscall/asm_plan9_riscv64.s:/·Syscall/
// The results are where errstr(2) wants its arguments: buf at 0(FP)
// and nbuf (the low half of the length word) at 8(FP).
TEXT runtime·errstr(SB),NOSPLIT|NOFRAME,$0-16
	MOV	g_m(g), T0
	MOV	(m_mOS+mOS_errstr)(T0), T1
	MOV	T1, ret_base+0(FP)
	MOV	$ERRMAX, T2
	MOV	T2, ret_len+8(FP)
	MOV	$SYS_ERRSTR, X8
	ECALL
	MOV	T1, T2
loop:
	MOVBU	0(T2), T3
	BEQ	T3, ZERO, done
	ADD	$1, T2
	JMP	loop
done:
	SUB	T1, T2
	MOV	T2, ret_len+8(FP)
	RET
