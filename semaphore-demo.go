package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"runtime"
	"time"
)

var (
	// cpu核心线程数字，数字设置为0，可以保证获取的就等于cpu的线程数。
	worker = runtime.GOMAXPROCS(0)
	// 信号量
	sema = semaphore.NewWeighted(int64(worker))
	// 任务量
	task = make([]int, 16*worker)
)

func Work() {
	ctx := context.Background()
	for i := range task {
		// 将信号量从池中取出来，每次取出来1个
		if err := sema.Acquire(ctx, 1); err != nil {
			break
		}
		go func(i int) {
			// 将一个信号量重新返还给池。
			defer sema.Release(1)
			time.Sleep(time.Millisecond *50)
			task[i] = i + 1
		}(i)
	}
	// 这里是个技巧，这里获取的是全部的信号量，如果能获取成功，就证明上面全部执行，并且把信号量放入池子中了。
	// 如果这里无法获取全部的信号量就会阻塞，所以这种方法可以保证可以让上面的所有的获取信号的goroutine执行完毕。
	if err := sema.Acquire(ctx, int64(worker)); err != nil {
		fmt.Println("获取所有的值失败", err)
	}
	fmt.Println(task)
}
func main() {
	Work()
}
