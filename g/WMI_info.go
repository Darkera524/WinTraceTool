package g

import (
	"github.com/Darkera524/WinTraceTool/model"
	"time"
)

var (
	WMI_info_list *model.WmiResp
)

func GetWMIInfo(interval int){
	for {
		//var resp model.ProvidersResp
		var checksum string = "nil"
		hostname, err := Hostname()
		if err != nil {
			Logger().Println(err.Error())
		}

		req := model.RequestModel{
			Hostname: hostname,
			Checksum: checksum,
		}

		WMI_info_list = nil
		err = HbsClient.Call("Wmi.Infos", req, &WMI_info_list)
		if err != nil {
			Logger().Println(err.Error())
		}
		//fmt.Println(ProviderMap)

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
