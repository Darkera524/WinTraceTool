package g

import (
	"github.com/Darkera524/WinTraceTool/model"
	"time"
)

var (
	WMI_statistic_list []string
)

func GetWMIStatistics(interval int){
	for {
		//var resp model.ProvidersResp
		//var checksum string = "nil"
		hostname, err := Hostname()
		if err != nil {
			Logger().Println(err.Error())
		}

		req := model.RequestModel{
			Hostname: hostname,
		}

		WMI_statistic_list = nil
		err = HbsClient.Call("Wmi.statistics", req, &WMI_statistic_list)
		if err != nil {
			Logger().Println(err.Error())
		}
		//fmt.Println(ProviderMap)

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
