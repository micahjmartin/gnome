package gnome

import (
	"fmt"
	"os"

	"github.com/nullmonk/gnome/modules"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func Run(script string) {
	thread := &starlark.Thread{Name: "my thread"}
	opts := &syntax.FileOptions{
		Set:             true,
		While:           true,
		TopLevelControl: true,
		GlobalReassign:  true,
		Recursion:       false,
	}
	libs := starlark.StringDict{
		"crypto": &modules.Crypto,
		"assets": &modules.Assets,
	}
	_, err := starlark.ExecFileOptions(opts, thread, script, nil, libs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Failed to run script: %v", err)
	}
}
