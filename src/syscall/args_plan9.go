// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9 && !riscv64

package syscall

// dupArgs and createArgs return the system call arguments for dup(2)
// and the mode and permission arguments of create(2).  The Plan 9 C
// compilers for 386, amd64 and arm give every argument a word-sized
// stack slot, so they are passed as they are.

//go:nosplit
func dupArgs(oldfd, newfd int) (a1, a2 uintptr) {
	return uintptr(oldfd), uintptr(newfd)
}

//go:nosplit
func createArgs(mode int, perm uint32) (a2, a3 uintptr) {
	return uintptr(mode), uintptr(perm)
}
