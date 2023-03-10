package main

import (
	_ "embed"
	"flag"
	"os"
	"fmt"
	"path/filepath"

	"github.com/aminamid/hello/mcfg"
)

var (
	//go:embed version.txt
	version string

	initCfg bool
	showVersion bool
	cfgFile string
	optCfgPort string
	
)

func main() {

	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.BoolVar(&initCfg, "init", false, "create default cfgfile if it does not exist")
	flag.StringVar(&optCfgPort, "p", "10080", "remote address to connect")
	flag.StringVar(&cfgFile, "f", "./cfg.yaml", "config file")
	flag.Parse()

	if showVersion {
		f1, _ := filepath.Abs(".")
		f2, _ := filepath.Abs("logs")
		fmt.Printf("%s\n%s\n%s\n", version, f1, f2)
		os.Exit(0)
	}

	go mcfg.ConfigStart(os.Stderr, optCfgPort, cfgFile, initCfg)
}
