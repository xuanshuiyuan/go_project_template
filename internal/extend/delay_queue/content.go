// @Author xuanshuiyuan 2025/7/23 11:21:00
package delay_queue

import (
	"context"
	"fmt"
	"time"
)

// 传播机制示例
func demonstratePropagation() {
	// 创建根 Context
	root := context.Background()

	// 创建可取消的 Context
	ctx1, cancel1 := context.WithCancel(root)

	// 创建带超时的子 Context
	ctx2, cancel2 := context.WithTimeout(ctx1, 5*time.Second)
	defer cancel2()

	// 创建带值的子 Context
	ctx3 := context.WithValue(ctx2, "userID", "123")

	// 启动工作 Goroutine
	go worker(ctx3, "worker1")
	go worker(ctx3, "worker2")

	// 2 秒后取消根 Context
	time.Sleep(2 * time.Second)
	cancel1() // 这会导致 ctx1, ctx2, ctx3 都被取消

	time.Sleep(1 * time.Second)
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: 收到取消信号: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("%s: 正在工作...\n", name)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
