// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// The riscv64 Plan 9 kernel enters a new image with R8 pointing to the
// Tos structure, argc at 0(SP) and the argv array at 8(SP).
TEXT _rt0_riscv64_plan9(SB),NOSPLIT|NOFRAME,$0
	MOV	X8, _tos(SB)
	MOV	0(X2), A0	// argc
	ADD	$8, X2, A1	// argv
	MOV	$runtime·rt0_go(SB), T0
	JALR	ZERO, T0

GLOBL _tos(SB), NOPTR, $8
