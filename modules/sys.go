package modules

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"go.starlark.net/starlark"
)

func getos() (os string, arch string) {
	if strings.HasSuffix(runtime.GOOS, "bsd") {
		os = "PLATFORM_BSD"
	} else if runtime.GOOS == "darwin" {
		os = "PLATFORM_MACOS"
	} else if runtime.GOOS == "linux" {
		os = "PLATFORM_LINUX"
	} else if runtime.GOOS == "windows" {
		os = "PLATFORM_WINDOWS"
	} else {
		os = "PLATFORM_UNSPECIFIED"
	}

	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "i386"
	default:
		arch = runtime.GOOS
	}
	return
}

func run(cmd string, args []string, disown bool) (starlark.Value, error) {
	var stdout, stderr bytes.Buffer
	c := exec.Command(cmd, args...)
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	c.Dir = cwd
	c.Stdout = &stdout
	c.Stderr = &stderr
	if disown {
		if err = c.Start(); err != nil {
			return starlark.None, err
		}
		return starlark.None, c.Process.Release()
	}

	statusCode := starlark.MakeInt(0)
	if err := c.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			statusCode = starlark.MakeInt(exiterr.ExitCode())
		} else {
			return starlark.None, err
		}
	}

	d := starlark.NewDict(3)
	d.SetKey(starlark.String("stdout"), starlark.String(stdout.String()))
	d.SetKey(starlark.String("stderr"), starlark.String(stderr.String()))
	d.SetKey(starlark.String("status"), statusCode)
	return d, nil
}

func SysExec(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var cmd starlark.String
	var cmdArgs *starlark.List
	var disown starlark.Bool
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &cmd, &cmdArgs, &disown); err != nil {
		return nil, err
	}

	argsActual := make([]string, 0, cmdArgs.Len())
	iter := cmdArgs.Iterate()
	defer iter.Done()
	var x starlark.Value
	for iter.Next(&x) {
		i := x.(starlark.String).GoString()
		argsActual = append(argsActual, i)
	}

	return run(cmd.GoString(), argsActual, bool(disown))
}

func SysShell(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var cmd starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &cmd); err != nil {
		return nil, err
	}
	if runtime.GOOS == "windows" {
		return run("cmd.exe", []string{"/c", cmd.GoString()}, false)
	}
	return run("/bin/bash", []string{"-c", cmd.GoString()}, false)
}

// Implement https://docs.realm.pub/user-guide/eldritch#sys
func SysGetEnv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	env := os.Environ()
	d := starlark.NewDict(len(env))
	for _, k := range env {
		vals := strings.SplitN(k, "=", 2)
		d.SetKey(starlark.String(vals[0]), starlark.String(vals[1]))
	}
	return d, nil
}

func SysSetEnv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key starlark.String
	var value starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &key, &value); err != nil {
		return nil, err
	}
	return starlark.None, os.Setenv(key.GoString(), value.GoString())
}

func SysGetOs(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	o, arch := getos()
	r := map[string]interface{}{
		"arch":        arch,
		"desktop_env": "Unknown: Unknown",
		"distro":      "",
		"platform":    o,
	}
	return ToStarlarkValue(r)
}

func SysIsLinux(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	o, _ := getos()
	return starlark.Bool(o == "PLATFORM_LINUX"), nil
}

func SysIsWindows(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	o, _ := getos()
	return starlark.Bool(o == "PLATFORM_WINDOWS"), nil
}

func SysIsBSD(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	o, _ := getos()
	return starlark.Bool(o == "PLATFORM_BSD"), nil
}

func SysIsMacos(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	o, _ := getos()
	return starlark.Bool(o == "PLATFORM_MACOS"), nil
}

func SysGetPid(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	return starlark.MakeInt(os.Getpid()), nil
}

func SysHostname(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	h, err := os.Hostname()
	return starlark.String(h), err
}

func userFromUid(uid int) (map[string]interface{}, error) {
	uname, err := user.LookupId(fmt.Sprint(uid))
	if err != nil {
		return nil, err
	}
	gids, err := uname.GroupIds()
	if err != nil {
		return nil, err
	}
	group_ids := make([]string, 0, len(gids))
	groups := make([]string, 0, len(gids))

	for _, gid := range gids {
		group_ids = append(group_ids, gid)
		g, err := user.LookupGroupId(gid)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g.Name)
	}
	res := map[string]interface{}{
		"gid":       uname.Gid,
		"uid":       uname.Uid,
		"name":      uname.Username,
		"group_ids": group_ids,
		"groups":    groups,
	}
	return res, nil
}

func SysGetUser(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}

	user, err := userFromUid(os.Getuid())
	if err != nil {
		return nil, err
	}
	euser, err := userFromUid(os.Geteuid())
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{
		"uid":  user,
		"euid": euser,
		"gid":  os.Getgid(),
		"egid": os.Getegid(),
	}
	return ToStarlarkValue(m)
}

func SysGetIp(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}

	ifs, err := net.Interfaces()
	if err != nil {
		return starlark.None, err
	}

	res := make([]map[string]interface{}, 0, len(ifs))
	for _, i := range ifs {
		ips, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		ipsStr := make([]string, 0, len(ips))
		for _, ip := range ips {
			ipsStr = append(ipsStr, ip.String())
		}
		res = append(res, map[string]interface{}{
			"mac":  i.HardwareAddr.String(),
			"name": i.Name,
			"ips":  ipsStr,
		})
	}
	return ToStarlarkValue(res)
}

// Intentionally not implemented. These functions dont, error, they just return nil
var Sys = NewModule("sys", map[string]Function{
	"dll_inject":    nil,
	"dll_reflect":   nil,
	"exec":          SysExec,
	"get_env":       SysGetEnv,
	"set_env":       SysSetEnv,
	"get_ip":        SysGetIp,
	"get_os":        SysGetOs,
	"get_pid":       SysGetPid,
	"get_reg":       nil,
	"get_user":      SysGetUser,
	"hostname":      SysHostname,
	"is_windows":    SysIsWindows,
	"is_linux":      SysIsLinux,
	"is_bsd":        SysIsBSD,
	"is_macos":      SysIsMacos,
	"shell":         SysShell,
	"write_reg_hex": nil,
	"write_reg_int": nil,
	"write_reg_str": nil,
})
