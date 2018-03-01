package trace

import (
	"github.com/Darkera524/WinTraceTool/g"
)

/*
	使用建造者设计模式来使构造与实现分离
 */

 type TraceInterface interface {
 	getData(string) error
 	formatData(string) error
 	sendData([]byte) error
 }

 type Trace struct {
 	Trace TraceInterface
 }

 func (t *Trace) CreateTrace(guid string) {
 	if t == nil {
 		g.Logger().Print("nil trace")
 		return
	}
	err := t.Trace.getData(guid)
	if err != nil {
		//g.Logger().Println(err.Error())
	}
 }


