package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task 接口：任何"有名字、能执行、可能失败"的东西都是任务
type Task interface {
	Name() string
	Run() error
}

// Job：Task 的一种实现（模拟耗时工作，id==failAt 时失败）
type Job struct {
	id     int
	failAt int
}

func (j Job) Name() string { return fmt.Sprintf("任务%d", j.id) }
func (j Job) Run() error {
	time.Sleep(100 * time.Millisecond)
	if j.id == j.failAt {
		return fmt.Errorf("执行失败")
	}
	return nil
}

// Processor：N 个 worker 并发处理，安全地收集失败清单
type Processor struct {
	workers int
	mu      sync.Mutex // 保护下面的 errors map
	errors  map[string]error
}

func NewProcessor(workers int) *Processor {
	return &Processor{workers: workers, errors: make(map[string]error)}
}

func (p *Processor) Process(ctx context.Context, tasks []Task) { // 空1
	jobs := make(chan Task)
	var wg sync.WaitGroup

	for w := 0; w < p.workers; w++ {
		wg.Add(1) // 空2：登记一个 worker（注意位置！）
		go func() {
			defer wg.Done() // 空3：worker 结束时注销
			for t := range jobs { // 空4：取任务，队列关闭自动下班
				select {
				case <-ctx.Done(): // 空5：超时/取消 → 立刻下班
					return
				default:
				}
				if err := t.Run(); err != nil {
					p.mu.Lock()// 空6：要改共享 map 了，先做什么？
					p.errors[t.Name()] = err
					p.mu.Unlock()// 空7：改完，别忘了
				}
			}
		}()
	}

	for _, t := range tasks {
		jobs <- t
	}
	close(jobs)// 空8：任务派完，宣布"没有新数据了"
	wg.Wait()// 空9：等所有 worker 收工
}

// Failed：安全读出失败清单（注意：读也要锁）
func (p *Processor) Failed() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	names := make([]string, 0, len(p.errors))
	for n := range p.errors {
		names = append(names, n)
	}
	return names
}

func main() {
	tasks := make([]Task, 10) // 空10：能装 10 个 Task 的切片（昨天错过一次的知识点）
	for i := 0; i < 10; i++ {
		tasks[i] = Job{id: i, failAt: 3} // 3 号任务会失败
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 空11：2秒超时的 ctx
	defer cancel()

	p := NewProcessor(3)
	p.Process(ctx, tasks)
	fmt.Println("失败的任务:", p.Failed())
}