package main

import (
	//"github.com/Darkera524/WinTraceTool/dispatch"
	"flag"
	"github.com/Darkera524/WinTraceTool/g"
	"github.com/Darkera524/WinTraceTool/http"
	"github.com/Darkera524/WinTraceTool/dispatch"
)

func main(){
	cfg := flag.String("c", "cfg.json", "configuration file")

	g.ParseConfig(*cfg)
	g.InitRpcClients()
	go http.Start()
	go g.GetProviders(60)

	go dispatch.TraceExec()
}


