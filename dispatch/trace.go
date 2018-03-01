package dispatch

import (
	"github.com/Darkera524/WinTraceTool/func/trace"
	"github.com/Darkera524/WinTraceTool/g"
	"github.com/Darkera524/WinTraceTool/model"
	"time"
)


func TraceExec() {
	model.TracingList = []string{}
	for {
		for k, name := range g.ProviderMap.MetricMap {
			isTracing := false
			for _, i := range model.TracingList {
				if i == k {
					isTracing = true
					break
				}
			}
			if isTracing == true {
				continue
			}

			if name == "dns" {
				trace_ins := trace.TraceBuilder{}

				traceins := trace.Trace{&trace_ins}
				go traceins.CreateTrace(k)
			}
		}
		time.Sleep(time.Duration(60) * time.Second)
	}
}
