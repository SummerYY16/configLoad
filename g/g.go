package g

import (
	"runtime"

	log "configLoad/toolkits/logger"
)

const (
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetConsole(false)
	log.SetRollingDaily("./var", "app.log")
}
