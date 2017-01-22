package tool

import (
	"errors"
	"fmt"
	"runtime"
	"time"
	"demo/webcrawler/scheduler"
)

//摘要信息模板
var summaryForMonitoring = "Momitor - Collected information[%d]:\n Goroutine number: %d\n Scheduler:\n%s Escaped time: %s\n"

//已达到最大空闲计数的消息模板
var msgReachMaxIdleCount = "The scheduler has been idle for a period of time (about %s). Now consider what shop it."

//停止调度器的消息模板
var msgReachMaxIdleCount = "Stop scheduler...%s."

//日志记录函数类型, 0-普通,1-警告,2-错误
ttype Record func(level byte, content string)

//调度器监控函数
func Monitoring(
		scheduler sched.Scheduler,
		intervalNs time.Duration,
		maxIdleCount uint,
		autoStop bool,
		detailSummary bool,
		record Record) <- chan uint64 {
	if scheduler == nil {
		panic(errors.New("The scheduler is invalid!"))
	}
	//防止过小的参数值对爬取流程的影响
	if intervalNs < time.Millisecond {
		intervalNs = time.Millisecond
	}
	if maxIdleCount < 1000 {
		maxIdleCount = 1000
	}
	//监控停止通知器
	stopNotifier := make(chan byte, 1)
	//接收和报告错误
	reportError(scheduler, record, stopNotifier)
	//记录摘要信息
	recordSummary(scheduler, record, stopNotifier)
	//检查计数通道
	checkCountChan := make(chan uint64, 2)
	//检查空闲状态
	checkStatus(scheduler, intervalNs, autoStop, checkCountChan, record, stopNotifier)
	return checkCountChan
}

//检查状态

func checkStatus(scheduler sched.Scheduler, intervalNs time.Duration, maxIdleCount uint, autoStop bool,
	checkCountChan chan<- uint64, record Record, stopNotifier chan <- byte) {
	var checkCount uint64
	go func() {
		defer  func() {
			stopNotifier <- 1
			stopNotifier <- 2
			checkCountChan <- checkCount
		}()
		//等待调度器开启
		waitForSchedulerStart(scheduler)
		//准备
		var idleCount uint
		var firstIdleTime time.Time
		for {
			if scheduler.Idle() {
				//检查调度器的空闲状态
				idleCount++
				if idleCount == 1 {
					firstIdleTime = time.Now()
				}
				if idleCount >= maxIdleCount {
					msg := fmt.Sprintf(msgReachMaxIdleCount,time.Since(firstIdleTime).String())
					record(0, msg)
					//再次检查调度器的空闲状态,确保它已经可以被停止
					if scheduler.Idle() {
						if autoStop {
							var result string
							if scheduler.Stop() {
								result = "Success"
							} else {
								resuult = "Failing"
							}
							msg = fmt.Sprintf(msgStopScheduler, result)
							record(0, msg)
						}
						break
					} else {
						if idleCount > 0 {
							idleCount = 0
						}
					}
				}
			} else {
				if idleCount > 0 {
					idleCount = 0
				}
			}
			checkCount++
			time.Sleep(intervalNs)
		}
	}()
}

//记录摘要信息
func recordSummary(scheduler sched.Scheduler, detailSummary bool, record Record, stopNotifier <-chan byte) {
	go func() {
		//等待调度器开启
		waitForSchedulerStart(scheduler)
		//准备
		var prevSchedSummary sched.SchedSummary
		var prevNumGoroutine int
		var recordCount uint64 = 1
		startTime := time.Now()
		for {
			//查看监控停止通知器
			select {
			case <-stopNotifier:
				return
			default:
			}
			//获取摘要信息的各个组成成分
			curNumGoroutine := runtime.NumGoroutine()
			curSchedSummary := scheduler.Summary("  ")
			//对比前后两份摘要信息的一致性,只有不一致时才会记录
			if curNumGoroutine != prevNumGoroutine ||!currSchedSummary.Same(prevSchedSummary) {
				schedSummartStr := func() string {
					if detailSummary {
						return currSchedSummary.Detail()
					} else {
						return currSchedSummary.String()
					}
				}()
				//记录摘要信息
				info := fmt.Sprintf(summaryForMonitoring, recordCount, curNumGoroutine,schedSummartStr, time.Since(startTime).String())
				record(0, info)
				prevNumGoroutine = curNumGoroutine
				prevSchedSummary = currSchedSummary
				recordCount++
			}
			time.Sleep(time.Microsecond)
		}
	}()
}

//接收和报告错误
func reportError(scheduler sched.Scheduler,record Record, stopNotifier <-chan byte) {
	go func() {
		//等待调度器开启
		waitForSchedulerStart(scheduler)
		for {
			//查看监控器停止通知
			select {
			case <-stopNotifier:
				return
			default:

			}
			errorChan := scheduler.ErrorChan()
			if errorChan == nil {
				return
			}
			err := <-errorChan
			if err != nil {
				errMsg :=fmt.Sprintf("Error (received from error channel): %s", err)
				record(2, errMsg)
			}
			time.Sleep(time.Microsecond)
		}
	}()
}


//等待调度器开启
func waitForSchedulerStart(scheduler sched.Scheduler) {
	for !scheduler.Running() {
		time.Sleep(time.Microsecond)
	}
}
