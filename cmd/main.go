package main

import (
	"github.com/yaza-putu/online-test-bookandlink/internal/config"
	"github.com/yaza-putu/online-test-bookandlink/internal/core"
	"runtime"
)

func main() {
	// set max cpu
	runtime.GOMAXPROCS(config.App().MaxCpu)

	// load env
	core.Env()

	// init database
	core.Database()

	// redis
	//core.Redis()
	// start server
	core.HttpServe()
}
