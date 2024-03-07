package gnome

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/nullmonk/gnome/modules"
	"go.starlark.net/starlark"
	"golang.org/x/sys/unix"
)

// Execveat executes at program at a path using a file descriptor.
// The go runtime process image is replaced by the executable described
// by the directory file descriptor and pathname.
func Execveat(fd uintptr, pathname string, argv []string, envv []string, flags int) error {
	pathnamep, err := syscall.BytePtrFromString(pathname)
	if err != nil {
		return err
	}

	argvp, err := syscall.SlicePtrFromStrings(argv)
	if err != nil {
		return err
	}

	envvp, err := syscall.SlicePtrFromStrings(envv)
	if err != nil {
		return err
	}

	_, _, errno := syscall.Syscall6(
		unix.SYS_EXECVEAT,
		fd,
		uintptr(unsafe.Pointer(pathnamep)),
		uintptr(unsafe.Pointer(&argvp[0])),
		uintptr(unsafe.Pointer(&envvp[0])),
		uintptr(flags),
		0,
	)
	return errno
}

func Execute(buf []byte, argv []string, envv []string) error {
	flag := unix.MFD_CLOEXEC
	// if close-on-exec flag has been set when fd points to a script,
	// then fexecve() fails with the error ENOENT. Peek this to undo if its a script
	if buf[0] == '#' && buf[1] == '!' {
		flag &= ^unix.MFD_CLOEXEC
	}

	fd, err := unix.MemfdCreate("", flag)
	if err != nil {
		return err
	}

	f := os.NewFile(uintptr(fd), fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), fd))
	if _, err := f.Write(buf); err != nil {
		return err
	}
	return Execveat(f.Fd(), "", argv, os.Environ(), unix.AT_EMPTY_PATH)
}

// Stop execution of all the thread, then fexec a binary from either the system or the asset locker
func fallback(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var code starlark.String
	var cmdArgs *starlark.List
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &code, &cmdArgs); err != nil {
		return nil, err
	}

	thread.Cancel("fallback")
	// Convert args to []string
	argsActual := make([]string, 0, cmdArgs.Len())
	iter := cmdArgs.Iterate()
	defer iter.Done()
	var x starlark.Value
	for iter.Next(&x) {
		i := x.(starlark.String).GoString()
		argsActual = append(argsActual, i)
	}

	// First, if this matches an asset, call the asset
	if modules.AssetLocker != nil && modules.AssetLocker.Exists(code.GoString()) {
		// Winner winner
		flag := unix.MFD_CLOEXEC
		// if close-on-exec flag has been set when fd points to a script,
		// then fexecve() fails with the error ENOENT. Peek this to undo if its a script
		// TODO undo this
		fd, err := unix.MemfdCreate("", flag)
		if err != nil {
			return nil, err
		}

		memfd := os.NewFile(uintptr(fd), code.GoString())
		if err := modules.AssetLocker.ReadFile(code.GoString(), memfd); err != nil {
			return nil, err
		}
		fmt.Println("calling", argsActual)
		return starlark.None, Execveat(memfd.Fd(), "", argsActual, os.Environ(), unix.AT_EMPTY_PATH)
	}

	// See if its a binary on disk
	f, err := os.Open(code.GoString())
	if err == nil {
		return starlark.None, Execveat(f.Fd(), "", argsActual, os.Environ(), unix.AT_EMPTY_PATH)
	}
	return starlark.None, nil
}
