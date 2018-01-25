package g

import (
	"github.com/Darkera524/WinTraceTool/model"
	"time"
)

var (
	ProviderMap *model.ProvidersResp
)

func GetProviders(interval int){
	for {
		//var resp model.ProvidersResp
		var checksum string = "nil"
		hostname, err := Hostname()
		if err != nil {

		}

		req := model.RequestModel{
			Hostname: hostname,
			Checksum: checksum,
		}

		err = HbsClient.Call("Trace.Providers", req, &ProviderMap)
		if err != nil {

		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

}
