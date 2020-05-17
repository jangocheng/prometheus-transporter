package sender

import (
	"prometheus-transporter/model"

	"github.com/prometheus/prometheus/prompb"
)

type Sender interface {
	Convert(*prompb.TimeSeries) interface{}
}

type BaseSender struct {
	q                   chan interface{}
	retryQ              chan interface{}
	addr                string
	concurrency         int
	maxSamplePerRequest int
	timeout             int
	retryNum            int
}

func (b *BaseSender) Init(addr string, qLength, retryQLength, conn, timeout, maxSamplePerRequest, retryNum int) {
	b.addr = addr
	b.q = make(chan interface{}, qLength)
	b.retryQ = make(chan interface{}, retryQLength)
	b.concurrency = conn
	b.maxSamplePerRequest = maxSamplePerRequest
	b.timeout = timeout
	b.retryNum = retryNum
}

func StartConsumeAndSend(ch chan *prompb.TimeSeries, conf *model.Config) {
	// TODO: judge the sender type and generate one
	// TODO: Init all the sender (including queue, sender workers and retryWorkers)
	// TODO: Start to consume and send
}
