package queue

import (
	"fmt"
	"github.com/prometheus/prometheus/prompb"
	"time"
)

func InitQueue(size int64) chan *prompb.TimeSeries {
	q := make(chan *prompb.TimeSeries, size)
	return q
}

func metricQueueLength(q chan *prompb.TimeSeries) {
	go func() {
		for {
			fmt.Println("the length of the queue is :", len(q))
		}
		time.Sleep(10)
	}()
}
