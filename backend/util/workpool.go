package util

import (
	"context"
	"fmt"
	"sync"
)

// 定义接口 可传任意参数
type TaskFunc func(args ...interface{})

// 定义任务实体，里面有方法和参数
type Task struct {
	F    TaskFunc
	Args interface{}
}

// 定义线程池对象
type WorkPool struct {
	Pool       chan *Task      //定义任务池
	WorkCount  int             //工作线程数量,决定初始化几个goroutine
	StopCtx    context.Context //上下文
	StopCancel context.CancelFunc
	WG         sync.WaitGroup //阻塞计数器
}

// 任务执行
func (t *Task) Execute(args ...interface{}) {
	t.F(args...)
}

// 实例化一个新线程池
func NewWrokPool(workerCount int, len int) *WorkPool {
	return &WorkPool{
		WorkCount: workerCount,
		Pool:      make(chan *Task, len),
	}
}

// 任务入队
func (w *WorkPool) PushTask(task *Task) {
	w.Pool <- task
}

// 任务调度 go协程从channel里取任务执行Execute方法
func (w *WorkPool) Work(wid int) {
	for {
		select {
		case <-w.StopCtx.Done():
			w.WG.Done()
			fmt.Printf("线程%d 退出执行了 \n", wid)
			return
		case t := <-w.Pool:
			if t != nil {
				t.Execute()
				fmt.Printf("f被线程%d执行了，参数为%v \n", wid, t.Args)
			}

		}
	}

}

// 启动线程池，触发任务调度
func (w *WorkPool) Start() *WorkPool {
	//定义好worker数量
	w.WG.Add(w.WorkCount)
	w.StopCtx, w.StopCancel = context.WithCancel(context.Background())
	for i := 0; i < w.WorkCount; i++ {
		//定义多少个协程来工作
		go w.Work(i)
	}
	return w
}

// 停止执行任务，回收正在执行任务的协程 协程计数器减1 直到变成0退出，否则阻塞
func (w *WorkPool) Stop() {
	w.StopCancel()
	w.WG.Wait()
}
