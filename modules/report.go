package modules

import (
	"go.starlark.net/starlark"
)

// Implement https://docs.realm.pub/user-guide/eldritch#report

func report(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, nil
}

// Intentionally not implemented. These functions dont, error, they just return nil
var Report = NewModule("report", map[string]Function{
	"match":       report,
	"match_all":   report,
	"replace":     report,
	"replace_all": report,
})
