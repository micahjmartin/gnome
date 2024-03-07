package modules

import (
	"time"

	strftime "github.com/itchyny/timefmt-go"
	"go.starlark.net/starlark"
)

// Implement https://docs.realm.pub/user-guide/eldritch#report

func timeFormatToEpoch(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var input starlark.String
	var format starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &input, &format); err != nil {
		return nil, err
	}
	t, err := strftime.Parse(input.GoString(), format.GoString())
	if err != nil {
		return nil, err
	}
	return starlark.MakeInt64(t.Unix()), nil
}

func timeFormatToReadable(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var input starlark.Int
	var format starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &input, &format); err != nil {
		return nil, err
	}
	sec, _ := input.Int64()
	s := strftime.Format(time.Unix(sec, 0).UTC(), format.GoString())
	return starlark.String(s), nil
}

func timeNow(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 0); err != nil {
		return nil, err
	}

	return starlark.MakeInt64(time.Now().Unix()), nil
}

func timeSleep(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var input starlark.Int
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &input); err != nil {
		return nil, err
	}
	i, _ := input.Int64()
	time.Sleep(time.Second * time.Duration(i))
	return starlark.None, nil
}

// Intentionally not implemented. These functions dont, error, they just return nil
var Time = NewModule("time", map[string]Function{
	"format_to_readable": timeFormatToReadable,
	"format_to_epoch":    timeFormatToEpoch,
	"now":                timeNow,
	"sleep":              timeSleep,
})
