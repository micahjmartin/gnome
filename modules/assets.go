package modules

import (
	"bytes"
	"fmt"
	"os"

	"github.com/nullmonk/gnome/sarcophagus"
	"go.starlark.net/starlark"
)

// Implement https://docs.realm.pub/user-guide/eldritch#assets
const env = "SYSTEMD_ROOT_IMAGE"
const hardcodedpth = "/lib/systemd/systemd-credentials"
const buildId = "3r09qwr-ajvm2ri-49a2ni5-ov9wn2d"

func GetAssetLocker() *sarcophagus.Vault {
	// Try all of the args first
	for _, f := range os.Args[1:] {
		v, err := sarcophagus.Open(f, buildId)
		if err == nil && v != nil {
			return v
		}
	}
	pth := os.Getenv(env)
	if pth == "" {
		pth = hardcodedpth
	}
	if _, err := os.Stat(pth); err != nil {
		return nil
	}
	v, _ := sarcophagus.Open(pth, buildId)
	return v
}

func assetsList(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	if AssetLocker != nil {
		return ToStarlarkValue(AssetLocker.Files())
	}
	return starlark.NewList(nil), nil
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

	f, err := os.Create(dst.GoString())
	if err != nil {
		return starlark.None, err
	}
	defer f.Close()
	if err := AssetLocker.ReadFile(name.GoString(), f); err != nil {
		return starlark.None, err
	}
	return starlark.None, nil
}

func assetsRead(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &name); err != nil {
		return nil, err
	}
	if AssetLocker == nil {
		return starlark.None, fmt.Errorf("asset locker not initialized")
	}

	var buf bytes.Buffer
	if err := AssetLocker.ReadFile(name.GoString(), &buf); err != nil {
		return starlark.None, err
	}
	return starlark.String(buf.String()), nil
}

func assetsReadBinary(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &name); err != nil {
		return nil, err
	}
	if AssetLocker == nil {
		return starlark.None, fmt.Errorf("asset locker not initialized")
	}

	var buf bytes.Buffer
	if err := AssetLocker.ReadFile(name.GoString(), &buf); err != nil {
		return starlark.None, err
	}
	return starlark.Bytes(buf.String()), nil
}

var AssetLocker = GetAssetLocker()

var Assets = NewModule("assets", map[string]Function{
	"copy":        assetsCopy,
	"list":        assetsList,
	"read_binary": assetsReadBinary,
	"read":        assetsRead,
})
