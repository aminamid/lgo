package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go loop(ctx)

	time.Sleep(2500 * time.Millisecond)

	cancel()

	cancel()

	time.Sleep(3 * time.Second)
	println("finish")
}

func loop(ctx context.Context) {
	conn, err := net.Dial("tcp", "munepow.com:443")
	if err != nil {
		println("dial failed")
		return
	}
	defer conn.Close()
	println("connected")

	go func() {
		buf := make([]byte, 1024)
		count := 0
		for {
			count += 1
			time.Sleep(100 * time.Millisecond)
			conn.Read(buf)
			fmt.Printf("test %07d %s\n", count, buf)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			println("cancel")
			if err := conn.Close(); err != nil {
				fmt.Printf("failed to close connection: %+v", err)
			}
			return
		}
	}
}
