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

	time.Sleep(3 * time.Second)
	println("finish")
}

func loop(ctx context.Context) {
	conn, err := net.Dial("tcp", "munepow.com:443")
	if err != nil {
		fmt.Printf("connection failed: %#v\n", err)
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
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("test ng: %#v\n", err)
			}
			fmt.Printf("test ok: %07d %s\n", count, string(buf[:n]))
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
func Retry(ctx context.Context, n uint, interval time.Duration, fn func() error) error {
	var err error
	for n > 0 {
		n -= 1
		if err = fn(); err == nil {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
	return err
}
