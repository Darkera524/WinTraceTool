package dispatch

import (
	"github.com/Darkera524/WinTraceTool/func/trace"
)

func TraceExec() {


	trace_ins := trace.TraceBuilder{}

	traceins := trace.Trace{&trace_ins}
	traceins.CreateTrace()


}
