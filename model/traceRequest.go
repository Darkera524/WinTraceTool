package model

type RequestModel struct {
	Hostname string
	Checksum string
}

type ProvidersResp struct {
	Providers map[string][]string
	MetricMap map[string]string
	Checksum string
}

type HttpRequest struct {
	Provider string
}
