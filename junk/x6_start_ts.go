package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now().UnixNano()
	currentTime := startTime
	for i := 0; i < 10; i++ {
		targetTime := currentTime + 1_000_000_000 // 1 seconds later
		currentTime = targetTime
		time.Sleep(time.Until(time.Unix(0, targetTime)))
		fmt.Println("Finished after", time.Now(), time.Now().UnixNano()-targetTime, "nanoseconds")
	}
}
