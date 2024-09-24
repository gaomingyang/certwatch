package scheduler

import (
	"fmt"
	"time"
)

// Task 定义任务函数类型
type Task func()

// Scheduler 任务调度器结构体
type Scheduler struct {
	ticker      *time.Ticker
	interval    time.Duration
	tasks       map[string]Task
	stopChannel chan bool
}

// NewScheduler 创建并返回一个新的 Scheduler 实例
func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		ticker:      time.NewTicker(interval),
		interval:    interval,
		tasks:       make(map[string]Task),
		stopChannel: make(chan bool),
	}
}

// AddTask 添加一个任务到调度器
func (s *Scheduler) AddTask(domain string, task Task) {
	s.tasks[domain] = task
	fmt.Printf("Task added for domain: %s\n", domain)
}

// Start 开始运行调度器，定期执行任务
func (s *Scheduler) Start() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.runTasks() // 定期执行任务
			case <-s.stopChannel:
				fmt.Println("Scheduler stopped.")
				return
			}
		}
	}()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.stopChannel <- true
}

// runTasks 执行所有已添加的任务
func (s *Scheduler) runTasks() {
	for domain, task := range s.tasks {
		fmt.Printf("Running task for domain: %s\n", domain)
		task()
	}
}
