package modules

// Implement https://docs.realm.pub/user-guide/eldritch#http

var Http = NewModule("http", map[string]Function{
	"download": nil,
	"get":      nil,
	"post":     nil,
})
