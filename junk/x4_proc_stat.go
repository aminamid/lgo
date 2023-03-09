package main

import (
	"bufio"
	//"time"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	for i := 0; i < 10; i++ {
		file, err := os.Open("/proc/stat")
		if err != nil {
			fmt.Println("Error opening /proc/stat:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if fields[0] == "btime" { // boot time
				btime := fields[1]
				fmt.Println("Boot time:", btime)
			} else if fields[0] == "cpu" { // system wide info
				if len(fields) < 8 {
					continue
				}
				user, _ := strconv.ParseUint(fields[1], 10, 64)
				nice, _ := strconv.ParseUint(fields[2], 10, 64)
				system, _ := strconv.ParseUint(fields[3], 10, 64)
				idle, _ := strconv.ParseUint(fields[4], 10, 64)
				iowait, _ := strconv.ParseUint(fields[5], 10, 64)
				irq, _ := strconv.ParseUint(fields[6], 10, 64)
				softirq, _ := strconv.ParseUint(fields[7], 10, 64)
				fmt.Printf("%4d %4d %4d %4d %4d %4d %4d\n", user%1000, nice%1000, system%1000, idle%1000, iowait%1000, irq%1000, softirq%1000)
				break
			}
		}
		//time.Sleep(1 * time.Second)

	}
}
