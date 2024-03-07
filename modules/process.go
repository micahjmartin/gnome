package modules

import (
	"os"
	"syscall"

	"go.starlark.net/starlark"
)

func processInfo(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	// TODO: Fill out the missing fields
	cwd, _ := os.Getwd()
	d := map[string]interface{}{
		"pid":                  os.Getpid(),
		"name":                 os.Args[0],
		"cmd":                  os.Args,
		"exe":                  nil,
		"environ":              os.Environ(),
		"cwd":                  cwd,
		"root":                 nil,
		"memory_usage":         nil,
		"virtual_memory_usage": nil,
		"status":               nil,
		"start_time":           nil,
		"ppid":                 os.Getppid(),
		"run_time":             nil,
		"uid":                  os.Getuid(),
		"euid":                 os.Geteuid(),
		"gid":                  os.Getgid(),
		"egid":                 os.Getegid(),
		"sid":                  nil,
	}
	return ToStarlarkValue(d)
}

func processKill(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var pid starlark.Int
	signal := starlark.MakeInt(int(syscall.SIGKILL))
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &pid, &signal); err != nil {
		return nil, err
	}
	pidActual, _ := pid.Int64()
	sigActual, _ := pid.Int64()
	return nil, syscall.Kill(int(pidActual), syscall.Signal(sigActual))
}

var Process = NewModule("process", map[string]Function{
	"info":    processInfo,
	"kill":    processKill,
	"list":    nil,
	"name":    nil,
	"netstat": nil,
})
