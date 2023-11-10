package domains

type BaseHttpResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type BaseErrorHttpResponse struct {
	Status int    `json:"status"`
	Path   string `json:"path"`
	Time   string `json:"time"`
	Error  string `json:"error"`
}
