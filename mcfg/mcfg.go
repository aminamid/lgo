package mcfg

import (
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mattn/go-jsonpointer"
	"gopkg.in/yaml.v3"
)

var mu sync.RWMutex
var sfcmCfg interface{}
var cfgPath string

//go:embed cfg.yml
var cfgDefault string

func wrap(errp *error, format string, args ...any) {
    if *errp != nil {
        *errp = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *errp)
    }
}

func loadJson(inputPath string, jsonObj *interface{}) error {
	byteArray, _ := ioutil.ReadFile(inputPath)
	return yaml.Unmarshal(byteArray, &jsonObj)
}

func saveJson(jsonObj interface{}, outputPath string) error {
	file, _ := os.Create(outputPath)
	defer file.Close()
	enc := yaml.NewEncoder(file)
	return enc.Encode(jsonObj)
}

func getCfg(c echo.Context) error {
	mu.RLock()
	defer mu.RUnlock()
	return c.JSON(http.StatusOK, sfcmCfg)
}

func ReadCfg(pathJson string, i interface{}) interface{} {
	mu.RLock()
	defer mu.RUnlock()
	it, err := jsonpointer.Get(sfcmCfg, pathJson)
	if err != nil {
		//fmt.Printf("%s\n",err)
		return i
	}
	return it
}

func CfgAsYaml(pathJson string) []byte {
	i := ReadCfg(pathJson, sfcmCfg)
	d, err := yaml.Marshal(i)
	if err != nil {
		return []byte(fmt.Sprintf("{ error: %s}", err.Error()))
	}
	return d
}


func putCfg(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Printf("ERROR parse failed :%v\n", err)
		return c.JSON(http.StatusInternalServerError, sfcmCfg)
	}
	fmt.Printf("%v\n----%v\n", c.Request().URL.Path, body)
	var val interface{}
	_ = yaml.Unmarshal(body, &val)
	mu.Lock()
	defer mu.Unlock()
	err = jsonpointer.Set(sfcmCfg, c.Request().URL.Path, val)
	if err != nil {
		fmt.Printf("ERROR Set failed :%v\n", err)
		return c.JSON(http.StatusInternalServerError, sfcmCfg)
	}
	saveJson(sfcmCfg, cfgPath)
	return c.JSON(http.StatusOK, sfcmCfg)
}

func ConfigStart(cherr chan error, logf io.Writer, port string, cfgpath string, initCfg bool) {
	cfgPath = cfgpath
	if _, err := os.Stat(cfgPath); err != nil {
		if !initCfg {
			wrap(&err, "config file does not exist %s", cfgPath)
			cherr <- err
			return
		} else {
			log.Printf("Creating config file %s", cfgPath)
			if err := yaml.Unmarshal([]byte(cfgDefault), &sfcmCfg) ; err != nil {
				wrap(&err, "failed to load default config file")
				cherr <- err
				return
			}
			if err := saveJson(sfcmCfg, cfgPath); err != nil {
				wrap(&err, "failed to save default config file")
				cherr <- err
				return
			}
		}

	}
	if err := loadJson(cfgPath, &sfcmCfg); err != nil {
		wrap(&err, "failed to load config file")
		cherr <- err
		return
	}
	//mu = sync.RWMutex{}
	e := echo.New()
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
		l.SetOutput(logf)
	}
	//e.Logger.Printf("sfcmCfg=%#v\n", sfcmCfg)
	e.HideBanner = true
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logf,
		Format: `${time_rfc3339} remote_ip=${remote_ip}:` +
			`host=${host}:method=${method}:uri=${uri}:user_agent=${user_agent}:` +
			`status=${status}:error=${error}:latency=${latency}:latency_human=${latency_human}:` +
			`bytes_in=${bytes_in}:bytes_out=${bytes_out}` + "\n",
	}))
	e.GET("/", getCfg)
	e.PUT("/*", putCfg)
	//e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	//	fmt.Fprintf(os.Stderr, "Request: %v\n", string(reqBody))
	//}))
	cherr <- nil

	if err := e.Start(":" + port); err != http.ErrServerClosed {
		wrap(&err, "server error")
		cherr <- err
		return
	}
	return

}
