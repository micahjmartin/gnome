package modules

// Implement https://docs.realm.pub/user-guide/eldritch#http

var Pivot = NewModule("pivot", map[string]Function{
	"arp_scan":           nil,
	"bind_proxy":         nil,
	"ncat":               nil,
	"port_forward":       nil,
	"port_scan":          nil,
	"smb_exec":           nil,
	"ssh_copy":           nil,
	"ssh_exec":           nil,
	"ssh_password_spray": nil,
})
