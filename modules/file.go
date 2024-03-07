package modules

// Implement https://docs.realm.pub/user-guide/eldritch#file

var File = NewModule("file", map[string]Function{
	"append":      nil,
	"compress":    nil,
	"decompress":  nil,
	"copy":        nil,
	"exists":      nil,
	"follow":      nil,
	"is_dir":      nil,
	"is_file":     nil,
	"list":        nil,
	"mkdir":       nil,
	"moveto":      nil,
	"parent_dir":  nil,
	"read":        nil,
	"remove":      nil,
	"replace":     nil,
	"replace_all": nil,
	"template":    nil,
	"timestomp":   nil,
	"write":       nil,
	"find":        nil,
})
