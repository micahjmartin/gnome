package modules

import (
	"fmt"

	"go.starlark.net/starlark"
)

func assetsCopy(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	var dst starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &src, &dst); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("assets.copy not impemented")
}

func assetsList(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("assets.list not impemented")
}

func assetsReadbinary(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("assets.read_binary not impemented")
}

func assetsRead(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &src); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("assets.read not impemented")
}

var Assets = Module{
	"copy":        starlark.NewBuiltin("", assetsCopy),
	"list":        starlark.NewBuiltin("", assetsList),
	"read_binary": starlark.NewBuiltin("", assetsReadbinary),
	"read":        starlark.NewBuiltin("", assetsRead),
}
