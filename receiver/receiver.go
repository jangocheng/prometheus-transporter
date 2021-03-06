package receiver

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"

	"prometheus-transporter/model"
)

func Start(logger *model.Logger, conf *model.Config) http.Server {
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
		return
	})

	http.HandleFunc("/api/v1/receive", func(w http.ResponseWriter, r *http.Request) {
		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errW := fmt.Errorf("read http body failed. %w", err)
			logger.Error("module", "receiver",
				"msg", errW)
			http.Error(w, errW.Error(), http.StatusInternalServerError)
			return
		}

		reqBuf, err := snappy.Decode(nil, compressed)
		if err != nil {
			errW := fmt.Errorf("snappy decode failed. %w", err)
			logger.Error("module", "receiver",
				"msg", errW)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req prompb.WriteRequest
		if err := proto.Unmarshal(reqBuf, &req); err != nil {
			errW := fmt.Errorf("proto unmarshal failed. %w", err)
			http.Error(w, errW.Error(), http.StatusBadRequest)
			return
		}

		for _, ts := range req.Timeseries {
			//m := make(prom_model.Metric, len(ts.Labels))
			fmt.Println(ts.Labels, len(ts.Samples))
			/*
				for i, l := range ts.Labels {
					m[prom_model.LabelName(l.Name)] = prom_model.LabelValue(l.Value)
					fmt.Println(i, l)
				}
			*/

			/*
				for _, s := range ts.Samples {
					fmt.Printf("  %f %d\n", s.Value, s.Timestamp)
				}
			*/
		}
	})

	s := &http.Server{
		Addr:           conf.HTTP,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info(
			"module", "receiver",
			"msg", fmt.Sprintf("starting http server, listening on:%s", conf.HTTP),
		)
		s.ListenAndServe()
		logger.Info(
			"module", "receiver",
			"msg", fmt.Sprintf("http server shutting down"),
		)
	}()

	return s
}

func handle() {}
