package modules

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"syscall"

	"go.starlark.net/starlark"
)

// Implement https://docs.realm.pub/user-guide/eldritch#file
func fileAppend(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path starlark.String
	var content starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &path, &content); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path.GoString(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err = f.WriteString(content.GoString()); err != nil {
		return nil, err
	}
	return starlark.None, nil
}

func fileMoveTo(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	var dst starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &src, &dst); err != nil {
		return nil, err
	}
	return starlark.None, os.Rename(src.GoString(), dst.GoString())
}

func fileIsFile(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src); err != nil {
		return nil, err
	}

	if st, err := os.Stat(src.GoString()); err != nil || st.IsDir() {
		return starlark.False, nil
	}
	return starlark.True, nil
}

func fileIsDir(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src); err != nil {
		return nil, err
	}

	if st, err := os.Stat(src.GoString()); err != nil || !st.IsDir() {
		return starlark.False, nil
	}
	return starlark.True, nil
}

func fileMkDir(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	var parent starlark.Bool
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src, &parent); err != nil {
		return nil, err
	}
	if parent {
		return starlark.None, os.MkdirAll(src.GoString(), 0755)
	}
	return starlark.None, os.Mkdir(src.GoString(), 0755)
}

func fileExists(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src); err != nil {
		return nil, err
	}
	if _, err := os.Stat(src.GoString()); err != nil {
		return starlark.False, nil
	}
	return starlark.True, nil
}

func fileRemove(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src); err != nil {
		return nil, err
	}
	return starlark.None, os.RemoveAll(src.GoString())
}

func fileList(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src); err != nil {
		return nil, err
	}

	if st, err := os.Stat(src.GoString()); err == nil && st.IsDir() {
		src = src + starlark.String("/*")
	}
	files, err := filepath.Glob(src.GoString())
	if err != nil {
		return nil, err
	}
	res := make([]interface{}, 0, len(files))
	for _, f := range files {
		abs, _ := filepath.Abs(f)
		st, err := os.Stat(f)
		if err != nil {
			return nil, err
		}
		typ := "File"
		if st.IsDir() {
			typ = "Directory"
		}
		usern := ""
		group := ""
		group_name := ""
		if runtime.GOOS != "windows" {
			adv, ok := st.Sys().(*syscall.Stat_t)
			if !ok {
				break
			}
			u, err := user.LookupId(fmt.Sprint(adv.Uid))
			if err != nil {
				return nil, err
			}
			usern = u.Username
			g, err := user.LookupGroupId(fmt.Sprint(adv.Uid))
			if err != nil {
				return nil, err
			}
			group = g.Gid
			group_name = g.Name

		}
		fil := map[string]interface{}{
			"size":          st.Size(),
			"file_name":     st.Name(),
			"absolute_path": abs,
			"permissions":   fmt.Sprintf("0%o", st.Mode().Perm()),
			"type":          typ,
			"modified":      st.ModTime().UTC().Format("2006-01-02 15:04:05 MST"),
			"owner":         usern,
			"group":         group,
			"group_name":    group_name,
		}
		res = append(res, fil)
	}
	return ToStarlarkValue(res)
}

func fileRead(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src); err != nil {
		return nil, err
	}

	b, err := os.ReadFile(src.GoString())
	return starlark.String(string(b)), err
}

func fileWrite(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	var contents starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &src, &contents); err != nil {
		return nil, err
	}

	return starlark.None, os.WriteFile(src.GoString(), []byte(contents.GoString()), 0644)
}

// Chmod a file *nix only
func fileChmod(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var file starlark.String
	var permissions starlark.Int
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &file, &permissions); err != nil {
		return nil, err
	}

	if runtime.GOOS != "linux" && runtime.GOOS != "bsd" {
		return starlark.None, fmt.Errorf("file.chmod not implemented on PLATFORM_WINDOWS")
	}
	perm, ok := permissions.Int64()
	if !ok {
		return starlark.None, fmt.Errorf("invalid int: %v", permissions.String())
	}
	return starlark.None, os.Chmod(file.GoString(), os.FileMode(perm))
}

var File = NewModule("file", map[string]Function{
	"append":      fileAppend,
	"compress":    nil,
	"decompress":  nil,
	"copy":        nil,
	"exists":      fileExists,
	"follow":      nil,
	"is_dir":      fileIsDir,
	"is_file":     fileIsFile,
	"list":        fileList,
	"mkdir":       fileMkDir,
	"moveto":      fileMoveTo,
	"parent_dir":  nil,
	"read":        fileRead,
	"remove":      fileRemove,
	"replace":     nil,
	"replace_all": nil,
	"template":    nil,
	"timestomp":   nil,
	"write":       fileWrite,
	"find":        nil,
	"chmod":       fileChmod,
})
