package main

import (
	"os"
	"os/signal"
	"prometheus-transporter/components"
	"prometheus-transporter/queue"
	"prometheus-transporter/receiver"
	"syscall"
	"time"
)

func main() {
	// init components
	components.InitComponents()
	logger := components.GetLogger()
	conf := components.GetConfig()

	// queue
	queue := queue.InitQueue()

	// indexer

	// sender

	// receiver
	recv := receiver.Start(logger, conf)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	logger.Info(
		"module", "main",
		"msg", "get system signal...",
		"signal", <-ch)

	// gracefully shutdown
	recv.Shutdown(nil)
	timeout := 60 * time.Second
	for {
		select {
		case <-time.After(timeout):
			logger.Error("module", "main", "msg", "quit timeout", "queue_lost_ts_num", len(queue))
			break
		default:
			if len(queue) == 0 {
				logger.Info("module", "main", "msg", "queue is empty, quit successfully")
				break
			}
			time.Sleep(time.Second)
		}
	}
}
