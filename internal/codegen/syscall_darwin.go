package codegen

import "tinygo.org/x/go-llvm"

type DarwinSyscall struct {
	context *llvm.Context
	builder *llvm.Builder
	module  *llvm.Module
}

func (s *DarwinSyscall) Name() string {
	return "darwin"
}

func (s *DarwinSyscall) Write(fd int, bufPtr llvm.Value, length llvm.Value) (int, error) {
	i64 := s.context.Int64Type()

	// Function type: i64(i64, i64, i64, i64)
	fnType := llvm.FunctionType(i64, []llvm.Type{i64, i64, i64, i64}, false)

	// Ensure syscall function exists
	syscallFn := s.module.NamedFunction("syscall")
	if syscallFn.IsNil() {
		syscallFn = llvm.AddFunction(*s.module, "syscall", fnType)
	}

	// Darwin syscall numbers
	const SYS_WRITE = 0x2000004

	// Cast buffer pointer and length to i64
	bufAsI64 := s.builder.CreatePtrToInt(bufPtr, i64, "")
	lenAsI64 := length
	if length.Type() != i64 {
		lenAsI64 = s.builder.CreateIntCast(length, i64, "")
	}

	// Perform the syscall: write(fd, buf, len)
	s.builder.CreateCall(fnType, syscallFn, []llvm.Value{
		llvm.ConstInt(i64, SYS_WRITE, false),
		llvm.ConstInt(i64, uint64(fd), false),
		bufAsI64,
		lenAsI64,
	}, "")

	return 0, nil
}

func (s *DarwinSyscall) Exit(status int) {
	i64 := s.context.Int64Type()

	fnType := llvm.FunctionType(i64, []llvm.Type{i64, i64, i64, i64}, false)

	syscallFn := s.module.NamedFunction("syscall")
	if syscallFn.IsNil() {
		syscallFn = llvm.AddFunction(*s.module, "syscall", fnType)
	}

	const SYS_EXIT = 0x2000001

	s.builder.CreateCall(fnType, syscallFn, []llvm.Value{
		llvm.ConstInt(i64, SYS_EXIT, false),
		llvm.ConstInt(i64, uint64(status), false),
		llvm.ConstInt(i64, 0, false),
		llvm.ConstInt(i64, 0, false),
	}, "")
}
