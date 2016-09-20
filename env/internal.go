package env

// version represents the running applications version.
var version string = "devel"

var buildTime string = ""

// osBackend implements the env interfaces using the real local os state to
// fulfill the interface.
type osBackend struct{}
