package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aminamid/lgo/mcfg"
	"golang.org/x/exp/slog"
)

var (
	//go:embed version.txt
	version string

	initCfg     bool
	showVersion bool
	cfgFile     string
	optCfgPort  string
)

func main() {

	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.BoolVar(&initCfg, "init", false, "create default cfgfile if it does not exist")
	flag.StringVar(&optCfgPort, "p", "10080", "remote address to connect")
	flag.StringVar(&cfgFile, "f", "./cfg.yml", "config file")
	flag.Parse()

	if showVersion {
		f1, _ := filepath.Abs(".")
		f2, _ := filepath.Abs("logs")
		fmt.Printf("%s\n%s\n%s\n", version, f1, f2)
		os.Exit(0)
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr))

	errChan := make(chan error, 1)
	go mcfg.ConfigStart(errChan, os.Stderr, optCfgPort, cfgFile, initCfg)
	select {
	case err := <-errChan:
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	//from now  mcfg.ReadCfg(pathJson string, i interface{}) interface{} is available
	//fmt.Printf("----\n")
	//fmt.Printf(string(mcfg.CfgAsYaml("/params")))
	//fmt.Printf("----\n")
	//fmt.Printf("%#v", mcfg.ReadCfg("/params", interface{}(nil)))
	//fmt.Printf("----\n")

	params := make(map[string]chan interface{})
	for k, v := range mcfg.ReadCfg("/params", interface{}(make(map[string]interface{}))).(map[string]interface{}) {
		switch {
		case strings.HasPrefix(k, "set_"):
			vl := v.([]interface{})
			//fmt.Printf("%#v\n", v)
			params[k] = ChanParamSetGen(vl[0].(map[string]interface{}), vl[1].([]interface{}))
		case strings.HasPrefix(k, "weight_"):
			vl := v.([]interface{})
			//fmt.Printf("%#v\n", v)
			params[k] = ChanParamWeightGen(vl)
		default:
			fmt.Printf("Unknown param generator %s", k)
			os.Exit(1)
		}
	}

}

func ChanParamWeightGen(arg1 []interface{}) chan interface{} {
	var total int32
	n := len(arg1)
	sliceWeight := make([]int32, n,n)
	sliceKey := make([]string,n,n)
	total = 0
	for i,rawl := range arg1 {
		l := rawl.([]interface{})
		//fmt.Printf("k=%#v, weight=%#v\n", l[0].(string), int32(l[1].(int)))
		total += int32(l[1].(int))
		sliceWeight[i]=total
		sliceKey[i]=l[0].(string)
	}
	chanId := range_random(1, total)
	chanOut := make(chan interface{}, 10)
	go func() {
		for {
		id := <- chanId
		for i := range sliceWeight {
			if id > sliceWeight[i] {
				continue
			}
			chanOut <- interface{}(sliceKey[i])
			break
		}
		}
	}()
	return chanOut
}


func ChanParamSetGen(arg1 map[string]interface{}, arg2 []interface{}) chan interface{} {

	var chanId chan int32
	switch {
	case arg2[0].(string) == "range_random":
		chanId = range_random(arg2[1].(int32), arg2[2].(int32))
	}

	templates := make(map[string]interface{})
	for k, v := range arg1 {
		//fmt.Printf("### k=%#v:v=%#v\n", k,v)
		switch {
		case strings.HasPrefix(k, "listfile_"):
			//fmt.Printf("######  listfile_\n")
			l, err := readFileToSlice(k)
			if err != nil {
				fmt.Printf("faile to readFileToSlice:  %#v", err)
				os.Exit(1)
			}
			templates[k] = interface{}(l)
		default:
			//fmt.Printf("###### default\n")
			templates[k] = interface{}(v)
		}
	}
	chanOut := make(chan interface{}, 10)
	go func() {
		for {
		id := <-chanId
			s := make(map[string]interface{})
			for k, _ := range arg1 {
				switch {
				case strings.HasPrefix(k, "tmpl_"):
					s[k] = fmt.Sprintf(templates[k].(string), id)
				case strings.HasPrefix(k, "listfile_"):
					s[k] = templates[k].([]string)[id]
				case strings.HasPrefix(k, "list_"):
					s[k] = templates[k].([]string)[id]
				default:
					fmt.Println("Unknown prefix %s", k)
					os.Exit(1)
				}
			}
			chanOut <- interface{}(s)
		}
	}()
	return chanOut
}
func range_random(base int32, num int32) chan int32 {
	ch := make(chan int32, 10)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	go func() {
		for {
			ch <- base + r.Int31n(num)
		}
	}()
	return ch
}
func readFileToSlice(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
