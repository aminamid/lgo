package main

import (
	"fmt"
	"time"
)

func main() {
	currentTime := time.Now()                                      // Get the current time
	microSeconds := currentTime.UnixNano() / 1000                  // Convert to microseconds by dividing nanoseconds by 1000
	fmt.Printf("Current time in microseconds: %d\n", microSeconds) // Print the current time in microseconds
	for i := 0; i < 10; i++ {
		fmt.Printf("%10d\n", currentTime.UnixNano()%1000)
		time.Sleep(1)
	}
}
