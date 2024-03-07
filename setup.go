package gnome

import (
	"fmt"
	"os"

	"github.com/nullmonk/gnome/modules"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func Run(script string) {
	thread := &starlark.Thread{Name: os.Args[0]}
	opts := &syntax.FileOptions{
		Set:             true,
		While:           true,
		TopLevelControl: true,
		GlobalReassign:  true,
		Recursion:       false,
	}
	libs := starlark.StringDict{
		"assets":  &modules.Assets,
		"crypto":  &modules.Crypto,
		"file":    &modules.File,
		"http":    &modules.Http,
		"pivot":   &modules.Pivot,
		"process": &modules.Process,
		"regex":   &modules.Regex,
		"report":  &modules.Report,
		"sys":     &modules.Sys,
		"time":    &modules.Time,
	}

	_, err := starlark.ExecFileOptions(opts, thread, script, nil, libs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Failed to run script: %v", err)
	}
}
