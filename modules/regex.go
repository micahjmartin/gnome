package modules

// Implement https://docs.realm.pub/user-guide/eldritch#regex

var Regex = NewModule("regex", map[string]Function{
	"match":       nil,
	"match_all":   nil,
	"replace":     nil,
	"replace_all": nil,
})
