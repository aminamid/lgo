    package mlog
     
    import (
    	"fmt"
    	"io"
    	"log"
    	"os"
    	"runtime"
    	"strconv"
    	"strings"
    	"time"
    )
     
    func LoggingSetting(logFile string) {
    	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    	if err != nil {
    		log.Fatalf("file=logFile err=%s", err.Error())
    	}
    	multiLogFile := io.MultiWriter(os.Stdout, logfile)
    	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    	log.SetOutput(multiLogFile)
    }
     
    func Mlog(logFile string) io.Writer {
    	logfile, err := os.OpenFile(fmt.Sprintf(logFile, time.Now().Format("20060102_150405")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    	if err != nil {
    		log.Fatalf("file=logFile err=%s", err.Error())
    	}
    	return io.MultiWriter(os.Stdout, logfile)
    }
     
    func Prefix() string {
    	pc, file, line, ok := runtime.Caller(1)
    	if !ok {
    		return "FileTo runtime.Caller"
    	}
    	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
    	funcname := runtime.FuncForPC(pc).Name()
    	fn := funcname[strings.LastIndex(funcname, ".")+1:]
    	return fmt.Sprintf("%s %s", filename, fn)
    }
