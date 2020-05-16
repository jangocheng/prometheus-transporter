package reader

import (
	"fmt"
	"prometheus-transporter/components"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/prometheus/prompb"
)

type ReadResponse struct {
	Status string           `json:"status"`
	Data   ReadResponseData `json:"data"`
}

type ReadResponseData struct {
	Result []ReadResponseTimeSeries `json:"result"`
}

type ReadResponseTimeSeries struct {
	Metric map[string]string `json:"metric"`
	Values [][]interface{}   `json:"values"`
}

func (series ReadResponseTimeSeries) TranstoStdTimeSeries() *prompb.TimeSeries {
	tmp := &prompb.TimeSeries{
		Labels: make([]*prompb.Label, 0), Samples: make([]prompb.Sample, 0),
	}

	for lk, lv := range series.Metric {
		tmp.Labels = append(tmp.Labels, &prompb.Label{
			Name:  lk,
			Value: lv,
		})
	}

	for i, _ := range series.Values {
		v := series.Values[i]
		if len(v) != 2 {
			continue
		}
		ts, ok1 := v[0].(float64)
		value, ok2 := v[1].(string)
		valueFloat, err := strconv.ParseFloat(value, 64)
		if !ok1 || !ok2 || err != nil {
			// TODO: log
			fmt.Printf("parse error:[%v-%v][%v-%v][%v-%v]\n", v[0], ok1, v[1], ok2, value, err)
		}
		tmpSample := prompb.Sample{Timestamp: int64(ts * 1000), Value: valueFloat}
		tmp.Samples = append(tmp.Samples, tmpSample)
	}

	return tmp
}

func Query(address, expression string) (*ReadResponse, error) {
	target := components.MakeURL(address, "/api/v1/query", map[string]string{})
	param := map[string]string{"query": expression}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	paramStr, _ := json.Marshal(param)

	data := &ReadResponse{}
	_, _, errs := gorequest.New().Post(target).Type("multipart").Send(string(paramStr)).EndStruct(data)
	if errs != nil {
		fmt.Println(errs)
		return nil, fmt.Errorf("%v", errs)
	}
	return data, nil
}

func QueryRange(address, expression string, start, end, step int64) (*ReadResponse, error) {
	target := components.MakeURL(address, "/api/v1/query_range", map[string]string{})
	param := map[string]string{
		"query": expression,
		"start": strconv.Itoa(int(start)),
		"end":   strconv.Itoa(int(end)),
		"step":  strconv.Itoa(int(step)),
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	paramStr, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	data := &ReadResponse{}
	_, _, errs := gorequest.New().Post(target).Type("multipart").Send(string(paramStr)).EndStruct(data)
	if errs != nil {
		return nil, fmt.Errorf("%v", errs)
	}
	return data, nil
}

// TimeRangeSplit return a splited timestamp slice between start and end
func timeRangeSplit(start, end, step int64) []int64 {
	if start < 0 || end < 0 || start >= end {
		return []int64{}
	}
	timeSlice := make([]int64, 0)
	ts := start - start%step
	for ts < end {
		timeSlice = append(timeSlice, ts)
		ts = ts + step
	}
	return timeSlice
}
