package model

var(
	TracingList []string
)

type RequestModel struct {
	Hostname string
	Checksum string
}

type ProvidersResp struct {
	//Providers map[string][]string
	MetricMap map[string]string
	Checksum string
}

type HttpRequest struct {
	Provider string
}

type Pg_dns struct {
	querytime string
	eventCode string

}
