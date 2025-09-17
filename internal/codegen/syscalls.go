package codegen

import (
	"strings"

	"tinygo.org/x/go-llvm"
)

type TargetSystem int

const (
	SystemInvalid TargetSystem = iota
	SystemDarwin
	SystemIOS
	SystemLinux
	SystemWindows
	SystemUnknown
)

// Syscall represents a minimal syscall across platforms
type Syscall struct {
	Name     string // Logical name
	LinuxX86 int    // x86_64 Linux syscall number
	LinuxA64 int    // ARM64 Linux syscall number
	Darwin   int    // macOS syscall number (x86_64 & ARM64)
	Windows  string // Windows API function (kernel32/ntdll)
	WASI     string // WASI import
}

type SyscallWriter interface {
	// Write issues a raw write syscall with buffer pointer and length values already in IR
	Write(fd int, bufPtr llvm.Value, length llvm.Value) (n int, err error)
	Exit(status int)
}

var MinimalSyscalls = []Syscall{
	// Program lifecycle
	{"exit", 60, 93, 0x2000001, "ExitProcess", "proc_exit"},
	{"yield", 35, 101, 0, "", "sched_yield"},

	// I/O
	{"write", 1, 64, 0x2000004, "WriteFile", "fd_write"},
	{"read", 0, 63, 0x2000003, "ReadFile", "fd_read"},

	// Memory
	{"mmap", 9, 222, 0x20001, "VirtualAlloc", "memory_grow"},
	{"munmap", 11, 215, 0x20002, "VirtualFree", ""},

	// Process info / optional
	{"getpid", 39, 172, 0x200020, "", ""},

	// Networking (sockets)
	{"socket", 41, 198, 0x200061, "socket", "sock_socket"},
	{"connect", 42, 203, 0x20006c, "connect", "sock_connect"},
	{"accept", 43, 202, 0x200068, "accept", "sock_accept"},
	{"bind", 49, 200, 0x200067, "bind", "sock_bind"},
	{"listen", 50, 201, 0x200066, "listen", "sock_listen"},
	{"sendto", 44, 208, 0x20006a, "send", "sock_sendto"},
	{"recvfrom", 45, 205, 0x200069, "recv", "sock_recvfrom"},
}

func syscall(target string, context llvm.Context, builder llvm.Builder, module llvm.Module) (SyscallWriter, error) {

	switch parseTargetTriple(target) {
	case SystemDarwin:
		return &DarwinSyscall{
			&context,
			&builder,
			&module,
		}, nil
	default:
		panic("unhandled default case")
	}
	return nil, nil
}

func parseTargetTriple(target string) TargetSystem {
	parts := strings.Split(target, "-")
	if len(parts) < 3 {
		return SystemInvalid
	}
	switch strings.ToLower(parts[2]) {
	case "darwin", "ios", "macos", "darwin25.0.0":
		return SystemDarwin
	case "linux":
		return SystemLinux
	case "windows", "win32":
		return SystemWindows
	default:
		return SystemUnknown
	}
}
