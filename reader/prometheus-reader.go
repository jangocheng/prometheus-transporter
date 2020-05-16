package reader

import (
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/niean/gotools/concurrent/semaphore"
	"github.com/prometheus/prometheus/prompb"
)

/*
 * Query and Aggregate Data from multi-prometheus Backend
 * Streaming return the Prometheus data sorted by time restrictly
 * Can specify time range, query expression, step, data cache DIR, multi-prometheus BE
 * Can support raw data
 */

type PromReader struct {
	Address       []string
	StartTs       int64
	EndTs         int64
	DataStep      int64
	MigrationStep int64
	Expression    string
	outer         chan *PromReaderOutput
}
type PromReaderOutput struct {
	Start         int64
	End           int64
	MigrationStep int64
	DataStep      int64
	TimeSeries    *[]*prompb.TimeSeries
}

func NewPromReader(address []string, start, end int64, dStep, mStep int64, expression string, outer chan *PromReaderOutput) *PromReader {
	return &PromReader{
		Address:       address,
		StartTs:       start,
		EndTs:         end,
		DataStep:      dStep,
		MigrationStep: mStep,
		Expression:    expression,
		outer:         outer,
	}
}

func (r PromReader) Read(logger log.Logger) {
	// long term time duration
	// split to a ts series in order to be start & end of time range
	tss := timeRangeSplit(r.StartTs, r.EndTs, r.MigrationStep)
	orgnizeCh := make(chan chan PromReaderOutput, 10)
	go r.organizer(orgnizeCh)

	sema := semaphore.NewSemaphore(2)
	wg := sync.WaitGroup{}
	// multiple BE
	for i, _ := range tss {
		if i == len(tss)-1 {
			break
		}
		start := tss[i]
		end := tss[i+1]

		//level.Info(logger).Log("module", "prom_read", "msg", fmt.Sprintf("[%s] data fetch successful. step:%d", time.Unix(start, 0).Format("2006-01-02 15:04:05"), end-start))
		ch := make(chan PromReaderOutput, 20)
		orgnizeCh <- ch
		sema.Acquire()
		wg.Add(1)
		go func(ch chan PromReaderOutput, start, end int64) {
			defer sema.Release()
			defer wg.Done()
			r.readOneDurationData(logger, ch, start, end)
		}(ch, start, end)
	}
	wg.Wait()
	close(orgnizeCh)
}

func (r PromReader) organizer(tranChan chan chan PromReaderOutput) {
	for ch := range tranChan {
		for body := range ch {
			tmp := body
			r.outer <- &tmp
		}
	}
	close(r.outer)
}

func (r PromReader) readOneDurationData(logger log.Logger, resultChan chan PromReaderOutput, start, end int64) {
	wg := sync.WaitGroup{}
	for i, _ := range r.Address {
		wg.Add(1)
		go func(addr string, start, end, dStep int64) {
			defer wg.Done()
			data, err := QueryRange(addr, r.Expression, start, end, dStep)
			if err != nil {
				level.Error(logger).Log("msg", "query range failed.", "err", err.Error())
				return
			}

			allSeries := make([]*prompb.TimeSeries, len(data.Data.Result))
			for i, _ := range data.Data.Result {
				t := data.Data.Result[i].TranstoStdTimeSeries()
				allSeries[i] = t
			}

			level.Info(logger).Log(
				"module", "prom_reader",
				"msg", "read one time series success",
				"addr", addr,
				"start", time.Unix(start, 0).Format("2006-01-02 15:04:05"),
				"series_num", len(allSeries),
			)

			resultChan <- PromReaderOutput{
				Start:         start,
				End:           end,
				MigrationStep: r.MigrationStep,
				DataStep:      r.DataStep,
				TimeSeries:    &allSeries,
			}
		}(r.Address[i], start, end, r.DataStep)
	}
	wg.Wait()
	close(resultChan)
}
