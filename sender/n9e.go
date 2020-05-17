package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/prometheus/prometheus/prompb"
	"github.com/toolkits/pkg/logger"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type N9ESender struct {
	BaseSender
}

type N9EMetricValue struct {
	Metric       string      `json:"metric"`
	Endpoint     string      `json:"endpoint"`
	Timestamp    int64       `json:"timestamp"`
	Step         int64       `json:"step"`
	ValueUntyped interface{} `json:"value"`
	Tags         string      `json:"tags"`
}

func NewN9ESender() *N9ESender {
	return &N9ESender{}
}

func (n *N9ESender) Convert(ts *prompb.TimeSeries) interface{} {
	if len(ts.Labels) == 0 || len(ts.Samples) == 0 {
		return nil
	}

	tagsMap := map[string]string{}
	metric := ""
	endpoint := ""
	// TODO: get step or push to cache
	step := int64(15)

	// to map "__name__" to "n9e metric"
	// to map "instance" to "n9e endpoint"
	for _, l := range ts.Labels {
		if l.Name == "__name__" {
			metric = l.Value
			continue
		}

		if l.Name == "instance" {
			endpoint = l.Value
			continue
		}

		if l.Name == "" || l.Value == "" {
			continue
		}
		tagsMap[l.Name] = l.Value
	}

	tagString := sortedTags(tagsMap)

	n9ePoints := make([]*N9EMetricValue, 0)
	for _, sample := range ts.Samples {
		tmp := &N9EMetricValue{
			Metric:       metric,
			Endpoint:     endpoint,
			Timestamp:    sample.Timestamp,
			Step:         step,
			ValueUntyped: sample.Value,
			Tags:         tagString,
		}
		n9ePoints = append(n9ePoints, tmp)
	}

	return n9ePoints
}

func (n *N9ESender) Start() {
	// start send worker
	for i := 0; i < n.concurrency; i++ {
		go func(index int) {
			// TODO log: Xth worker started
			for {
				func() {
					defer func() {
						if err := recover(); err != nil {
							// TODO panic reason log
							fmt.Println(err)
						}
					}()

					select {
					case data := <-n.q:
						n.push(data.([]*N9EMetricValue))
					}
				}()
			}
		}(i)
	}
}

func (n N9ESender) push(items []*N9EMetricValue) {
	bs, err := json.Marshal(items)
	if err != nil {
		logger.Warning(err)
		return
	}

	bf := bytes.NewBuffer(bs)

	resp, err := http.Post(n.addr, "application/json", bf)
	if err != nil {
		logger.Warning(err)
		return
	}
	// TODO retry

	defer resp.Body.Close()
}

func sortedTags(tags map[string]string) string {
	if tags == nil {
		return ""
	}

	size := len(tags)
	if size == 0 {
		return ""
	}

	ret := bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if size == 1 {
		for k, v := range tags {
			ret.WriteString(k)
			ret.WriteString("=")
			ret.WriteString(v)
		}
		return ret.String()
	}

	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for j, key := range keys {
		ret.WriteString(key)
		ret.WriteString("=")
		ret.WriteString(tags[key])
		if j != size-1 {
			ret.WriteString(",")
		}
	}

	return ret.String()
}
