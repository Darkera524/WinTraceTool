package dispatch

import (
	"github.com/Darkera524/WinTraceTool/model"
	"github.com/Darkera524/WinTraceTool/g"
	"time"
	"github.com/Darkera524/WinTraceTool/func/wmi_collector"
)

func WmiExec() {
	model.WmiInfoList = []string{}
	for {
		for _, wmi_server := range g.WMI_info_list.WmiList {
			isTracing := false
			for _, wmi_running := range model.WmiInfoList {
				if wmi_server == wmi_running {
					isTracing = true
					break
				}
			}
			if isTracing == true {
				continue
			}


			wmi_ins := wmi_collector.WmiInfos{}

			wmiins := wmi_collector.Wmi{&wmi_ins}
			go wmiins.CreateWmi(wmi_server)

		}
		time.Sleep(time.Duration(60) * time.Second)
	}
}