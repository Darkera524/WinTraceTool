package wmi_collector

import (
	"github.com/Darkera524/WinTraceTool/g"
	"github.com/Darkera524/WinTraceTool/model"
	"time"
)

/*
	使用建造者设计模式来使构造与实现分离
 */

type WmiInterface interface {
	getData(string) error
	formatData(model.Wmi_Base) error
	sendData(string) error
}

type Wmi struct {
	Wmi WmiInterface
}

func (t *Wmi) CreateWmi(wmiClass string) {
	if t == nil {
		g.Logger().Print("nil trace")
		return
	}
	for {
		stoped := false
		err := t.Wmi.getData(wmiClass)
		if err != nil {
			//g.Logger().Println(err.Error())
		}
		time.Sleep(time.Duration(60) * time.Second)
		for i,serverClass := range g.WMI_info_list.WmiList{
			if serverClass == wmiClass{
				break
			}
			if i == len(g.WMI_info_list.WmiList)-1 {
				stoped = true
			}
		}
		if stoped{
			g.Logger().Println("wmi class " + wmiClass + "trace stopped")
			break
		}
	}
}


