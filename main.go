package main

import (
	//"github.com/Darkera524/WinTraceTool/dispatch"
	"flag"
	"github.com/Darkera524/WinTraceTool/g"
	"github.com/Darkera524/WinTraceTool/http"
	"github.com/Darkera524/WinTraceTool/dispatch"
	"time"
)

func main(){
	cfg := flag.String("c", "cfg.json", "configuration file")

	g.InitLog()
	g.InitEncoding()
	g.ParseConfig(*cfg)
	g.InitRpcClients()

	go g.GetProviders(60)
	go g.GetWMIInfo(60)

	go http.Start()
	time.Sleep(time.Duration(10) * time.Second)
	go dispatch.TraceExec()
	go dispatch.WmiExec()

	select {}

}


