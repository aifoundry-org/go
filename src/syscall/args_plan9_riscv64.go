// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall

// The riscv64 Plan 9 C compiler passes int and long (32-bit) arguments
// in 4-byte stack slots and aligns 8-byte arguments to 8, so two int
// arguments in a row share one 8-byte slot of Syscall's arguments.

//go:nosplit
func dupArgs(oldfd, newfd int) (a1, a2 uintptr) {
	return uintptr(uint32(oldfd)) | uintptr(newfd)<<32, 0
}

//go:nosplit
func createArgs(mode int, perm uint32) (a2, a3 uintptr) {
	return uintptr(uint32(mode)) | uintptr(perm)<<32, 0
}

//go:nosplit
func mountArgs(fd, afd int) uintptr {
	return uintptr(uint32(fd)) | uintptr(afd)<<32
}
