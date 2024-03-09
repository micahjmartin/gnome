package gnome

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/nullmonk/gnome/modules"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

//go:embed docs
var assets embed.FS

type script struct {
	name string
	src  interface{}
}

// Run a stark script, passing in the previous globals if specified
func Run() {
	// Set the asset locker to whatever we have specified
	modules.SetAssetLocker(assets)

	scripts_to_run := make([]script, 0, 1)
	if modules.AssetLocker != nil {
		err := fs.WalkDir(modules.AssetLocker, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !strings.HasSuffix(path, ".eldr") && !strings.HasSuffix(path, ".eldritch") {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			buf, err := fs.ReadFile(modules.AssetLocker, path)
			if err != nil {
				return nil
			}
			scripts_to_run = append(scripts_to_run, script{path, buf})
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] loading script from fs: %v\n", err)
		}
	}
	if len(os.Args) > 1 {
		scripts_to_run = append(scripts_to_run, script{os.Args[1], nil})
	}

	globals := starlark.StringDict{}
	var err error
	for _, s := range scripts_to_run {
		globals, err = run(s.name, s.src, globals)
		if err != nil {
			fmt.Printf("[!] '%s' failed: %v\n", s.name, err)
		}
	}
}

func run(name string, src interface{}, globals starlark.StringDict) (starlark.StringDict, error) {
	thread := &starlark.Thread{Name: name}
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
		"quit":     starlark.NewBuiltin("exit", quit),
		"fallback": starlark.NewBuiltin("fallback", fallback),
	}

	// Add the globals into the environment
	for k, v := range globals {
		libs[k] = v
	}
	res, err := starlark.ExecFileOptions(opts, thread, name, src, libs)
	if err != nil {
		if e, ok := err.(*starlark.EvalError); ok {
			// Check what the error message is, that is how we determined if we quit or exited
			lines := strings.SplitN(e.Msg, ": ", 2)
			if len(lines) < 2 {
				return nil, err
			}
			if lines[1] == "user exit" {
				// On exit calls, the interpreter also dies
				os.Exit(int(exitCode))
			} else if lines[1] == "user quit" {
				// on quit calls, only the script exits, not an error
				err = nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if globals == nil {
		globals = make(starlark.StringDict, len(res))
	}
	// Update globals with the results of this script
	for k, v := range res {
		globals[k] = v
	}
	return globals, nil
}
