package domain

import (
	"fmt"
)

type BaseHttpResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (b BaseHttpResponse) ShowData() {
	fmt.Println(b.Data)
}

func (b BaseHttpResponse) ShowStatus() {
	fmt.Println(b.Status)
}
