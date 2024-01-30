package domain

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

type BaseMessageResponse struct {
	Message string `json:"message"`
}
