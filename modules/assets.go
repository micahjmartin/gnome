package modules

// Implement https://docs.realm.pub/user-guide/eldritch#assets

var Assets = NewModule("assets", map[string]Function{
	"copy":        nil,
	"list":        nil,
	"read_binary": nil,
	"read":        nil,
})
