package sender

import (
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/prometheus/prompb"
)

func TestN9ESender_Convert(t *testing.T) {
	s := N9ESender{}
	wr := generateOneRequest(2, 4, 5)
	for _, ts := range wr.GetTimeseries() {
		points := s.Convert(ts)
		t.Log("prometheus timeseries:", ts)
		t.Log("n9e points:")
		for _, i := range points {
			t.Log(i.(*N9EMetricValue))
		}
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

		labels = append(labels, &prompb.Label{
			Name:  "__name__",
			Value: "metric-" + strconv.Itoa(se),
		})

		labels = append(labels, &prompb.Label{
			Name:  "instance",
			Value: "endpoint-" + strconv.Itoa(se),
		})

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
