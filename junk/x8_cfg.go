package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/aminamid/hello/mlog"
	"github.com/aminamid/hello/mcfg"
)

var (
	version     = "v1.0.3"
	showVersion bool
	optCfgPort string
)

func main() {
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.StringVar(&optCfgPort, "cfgport", "10080", "remote address to connect")
	flag.Parse()
	if showVersion {
		f1, _ := filepath.Abs(".")
		f2, _ := filepath.Abs("logs")
		fmt.Printf("%s\n%s\n%s\n", version, f1, f2)
	}
	logf := mlog.Mlog("logs/svr.%s.log")

	var wg sync.WaitGroup
	wg.Add(1)
	go mcfg.ConfigStart(logf, optCfgPort, "./cfg.yaml")
	time.Sleep(time.Second)
	wg.Wait()
}

