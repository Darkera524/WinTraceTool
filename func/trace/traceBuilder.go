package trace

import "fmt"

/*
	使用建造者设计模式来使构造与实现分离
 */

 type TraceInterface interface {
 	getData(string) error
 	formatData(string, string) error
 	sendData([]byte) error
 }

 type Trace struct {
 	Trace TraceInterface
 }

 func (t *Trace) CreateTrace(guid string) {
 	if t == nil {
 		fmt.Print("nil trace")
 		return
	}
	err := t.Trace.getData(guid)
	if err != nil {
		fmt.Println(err.Error())
	}
 }


