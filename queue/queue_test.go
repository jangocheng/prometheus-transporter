package queue

import (
	"github.com/prometheus/prometheus/prompb"
	"strconv"
	"testing"
	"time"
)

func Test_Queue(t *testing.T) {
	q := InitQueue(100)
	for i := 0; i < 100; i++ {
		q <- generate_one_timeseries(6, 6)
		t.Log("the length of queue is ", len(q))
	}

	for i := range q {
		t.Log("the length of queue is ", len(q), ", i: ", i)
		if len(q) == 0 {
			close(q)
		}
	}
}

func generate_one_timeseries(labelNum, sampleNum int) *prompb.TimeSeries {
	samples := make([]prompb.Sample, 0)
	// the step of test samples is 10s
	for sa := 0; sa < sampleNum; sa++ {
		samples = append(samples, prompb.Sample{
			Timestamp: time.Now().Unix() + int64(sa*10),
			Value:     100,
		})
	}

	labels := make([]*prompb.Label, 0)
	for la := 0; la < labelNum; la++ {
		labels = append(labels, &prompb.Label{
			Name:  "key" + "-" + strconv.Itoa(la),
			Value: "value" + "-" + strconv.Itoa(la),
		})
	}

	return &prompb.TimeSeries{
		Labels:  labels,
		Samples: samples,
	}
}
