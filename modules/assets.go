package modules

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"go.starlark.net/starlark"
)

// Implement https://docs.realm.pub/user-guide/eldritch#assets

func assetsList(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	if AssetLocker == nil {
		return starlark.None, fmt.Errorf("asset locker not initialized")
	}

	assets := make([]string, 0, 64)
	fs.WalkDir(AssetLocker, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		assets = append(assets, path)
		return nil
	})
	return ToStarlarkValue(assets)
}

func assetsCopy(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String
	var dst starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &name, &dst); err != nil {
		return nil, err
	}
	if AssetLocker == nil {
		return starlark.None, fmt.Errorf("asset locker not initialized")
	}

	f, err := AssetLocker.Open(name.GoString())
	if err != nil {
		return starlark.None, err
	}

	d, err := os.Create(dst.GoString())
	if err != nil {
		return starlark.None, err
	}
	defer f.Close()
	_, err = io.Copy(d, f)
	return starlark.None, err
}

func assetsRead(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &name); err != nil {
		return nil, err
	}
	if AssetLocker == nil {
		return starlark.None, fmt.Errorf("asset locker not initialized")
	}

	f, err := AssetLocker.Open(name.GoString())
	if err != nil {
		return starlark.None, err
	}
	buf, err := io.ReadAll(f)
	if err != nil {
		return starlark.None, err
	}
	return starlark.String(string(buf)), nil
}

func assetsReadBinary(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &name); err != nil {
		return nil, err
	}
	if AssetLocker == nil {
		return starlark.None, fmt.Errorf("asset locker not initialized")
	}

	f, err := AssetLocker.Open(name.GoString())
	if err != nil {
		return starlark.None, err
	}
	buf, err := io.ReadAll(f)
	if err != nil {
		return starlark.None, err
	}
	return starlark.Bytes(buf), nil
}

var AssetLocker fs.FS

func SetAssetLocker(f fs.FS) {
	AssetLocker = f
}

var Assets = NewModule("assets", map[string]Function{
	"copy":        assetsCopy,
	"list":        assetsList,
	"read_binary": assetsReadBinary,
	"read":        assetsRead,
})
