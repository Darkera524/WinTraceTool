package model

type Wmi_Base interface {

}

type Win32_PerfRawData_perfProc_Process struct {
	Name string
	IDProcess int64
	PercentProcessorTime int64
}
