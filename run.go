package gnome

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/nullmonk/gnome/modules"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// Kill the entire starlark thread
var exitCode int64

func exit(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var code starlark.Int
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &code); err != nil {
		return nil, err
	}
	exitCode, _ = code.Int64()
	thread.Cancel("user exit")
	return starlark.None, nil
}

func Run() {
	thread := &starlark.Thread{Name: os.Args[0]}
	opts := &syntax.FileOptions{
		Set:             true,
		While:           true,
		TopLevelControl: true,
		GlobalReassign:  true,
		Recursion:       false,
	}

	libs := starlark.StringDict{
		"assets":   &modules.Assets,
		"crypto":   &modules.Crypto,
		"file":     &modules.File,
		"http":     &modules.Http,
		"pivot":    &modules.Pivot,
		"process":  &modules.Process,
		"regex":    &modules.Regex,
		"report":   &modules.Report,
		"sys":      &modules.Sys,
		"time":     &modules.Time,
		"exit":     starlark.NewBuiltin("exit", exit),
		"fallback": starlark.NewBuiltin("fallback", fallback),
	}

	if modules.AssetLocker != nil {
		for _, f := range modules.AssetLocker.Files() {
			if strings.HasSuffix(f, ".eldr") || strings.HasSuffix(f, ".eldritch") {
				var buf bytes.Buffer
				modules.AssetLocker.ReadFile(f, &buf)
				_, err := starlark.ExecFileOptions(opts, thread, f, buf.Bytes(), libs)
				if err != nil {
					if _, ok := err.(*starlark.EvalError); ok {
						os.Exit(int(exitCode))
					}
					fmt.Fprintf(os.Stderr, "[!] Failed to run script: %v", err)
				}
			}
		}
	}
	if len(os.Args) > 1 {
		_, err := starlark.ExecFileOptions(opts, thread, os.Args[1], nil, libs)
		if err != nil {
			if _, ok := err.(*starlark.EvalError); ok {
				os.Exit(int(exitCode))
			}
			fmt.Fprintf(os.Stderr, "[!] Failed to run script: %v", err)
		}
	}

}
