package receiver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"prometheus-transporter/components"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"

	"github.com/prometheus/prometheus/prompb"
)

func Benchmark_Receiver(b *testing.B) {
	// Please note: MUST start receiver first
	err := components.ParseConfig("../dev.conf.toml")
	if err != nil {
		b.Error(err)
	}
	conf := components.GetConfig()
	body := generateOneRequest(100000, 4, 1)
	err = remoteWrite(conf.HTTP, body)
	if err != nil {
		b.Error(err)
	}
}

func Test_Receiver(t *testing.T) {
	startTestReceiver()
	body := generateOneRequest(100, 4, 1)
	conf := components.GetConfig()
	err := remoteWrite(conf.HTTP, body)
	if err != nil {
		t.Error(err)
	}
}

func generateOneRequest(seriesNum, labelNum, sampleNum int) *prompb.WriteRequest {
	samples := make([]prompb.Sample, 0)
	// the step of test samples is 10s
	for sa := 0; sa < sampleNum; sa++ {
		samples = append(samples, prompb.Sample{
			Timestamp: time.Now().Unix() + int64(sa*10),
			Value:     100,
		})
	}

	series := make([]*prompb.TimeSeries, 0)
	for se := 0; se < seriesNum; se++ {
		labels := make([]*prompb.Label, 0)
		for la := 0; la < labelNum; la++ {
			labels = append(labels, &prompb.Label{
				Name:  "key" + strconv.Itoa(se) + "-" + strconv.Itoa(la),
				Value: "value" + strconv.Itoa(se) + "-" + strconv.Itoa(la),
			})
		}

		series = append(series, &prompb.TimeSeries{
			Labels:  labels,
			Samples: samples,
		})
	}

	body := &prompb.WriteRequest{
		Timeseries: series,
	}

	return body
}

func startTestReceiver() error {
	// init
	err := components.ParseConfig("../dev.conf.toml")
	if err != nil {
		return err
	}
	conf := components.GetConfig()
	components.InitLogger(conf.Logger)
	logger := components.GetLogger()

	// start receiver
	go Start(logger, conf)

	// wait the receiver start
	time.Sleep(time.Second)
	return nil
}

func remoteWrite(address string, body *prompb.WriteRequest) error {
	target := components.MakeURL("http://"+address, "/api/v1/receive", map[string]string{})
	paramStr, _ := proto.Marshal(body)
	compressed := snappy.Encode(nil, paramStr)

	bd, errs := http.Post(target, "application/x-www-form-urlencoded", strings.NewReader(string(compressed)))
	if errs != nil {
		return fmt.Errorf("%+v", errs)
	}

	if bd.StatusCode != 200 {
		body, err := ioutil.ReadAll(bd.Body)
		if err != nil {
			body = []byte(fmt.Sprintf("[read body failed:%s]", err.Error()))
		}
		return fmt.Errorf("code is not 200. [code:%d][body:%s]", bd.StatusCode, string(body))
	}
	return nil
}
