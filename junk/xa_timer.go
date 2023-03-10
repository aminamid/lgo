package main

import (
	"fmt"
	"time"
)

func main() {
	// チャンネルを作成
	que := make(chan int64)

	// タイムスタンプを生成するgoroutine
	go func(q chan int64, timeDelta int64, testCount int) {
		startTime := time.Now().UnixNano()
		ts := startTime
		for i := 0; i < 50; i++ {
			targetTime := ts + timeDelta
			ts = targetTime
			q <- targetTime
		}
		close(q)
	}(que, 100_000_000, 30)

	// タイムスタンプを処理するgoroutine
	for i := 0; i < 3; i++ {
		go func(q chan int64, id int) {
			for ts := range q {
				time.Sleep(time.Until(time.Unix(0, ts)))
				fmt.Printf("%s: %d\n", time.Unix(0, ts).UTC().Format("2006-01-02 15:04:05.999999999"), id)
			}
		}(que, i)
	}

	// 終了を待つ
	time.Sleep(time.Second * 6)
}
