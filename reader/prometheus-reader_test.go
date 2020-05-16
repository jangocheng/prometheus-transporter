package reader

import (
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
)

func TestReader(t *testing.T) {
	resultCh := make(chan *PromReaderOutput, 10)
	reader := NewPromReader(
		[]string{"http://localhost:9090"},
		time.Now().Unix()-3600*24*30*3,
		time.Now().Unix()-3600*24*6,
		15,
		600,
		`{__name__=~'.+'}`,
		resultCh,
	)

	go reader.Read(log.NewLogfmtLogger(os.Stdout))
	for timeSeries := range resultCh {
		//fmt.Println(len(*timeSeries.TimeSeries))
		_ = timeSeries
	}
}
