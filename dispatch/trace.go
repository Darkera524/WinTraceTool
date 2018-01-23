package dispatch

import (
	"github.com/Darkera524/WinTraceTool/func/trace"
)

func TraceExec() {
	dns := trace.DNSBuilder{}

	traceins := trace.Trace{&dns}
	traceins.CreateTrace()


}
