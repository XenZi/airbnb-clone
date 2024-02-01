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

type BaseMessageResponse struct {
	Message string `json:"message"`
}

type AccommodationDTO struct {
	Id               string   `json:"id"`
	UserId           string   `json:"userId" `
	UserName         string   `json:"username" `
	Email            string   `json:"email" bson:"email"`
	Name             string   `json:"name" `
	Address          string   `json:"address" `
	City             string   `json:"city" `
	Country          string   `json:"country" `
	Conveniences     []string `json:"conveniences" `
	MinNumOfVisitors int      `json:"minNumOfVisitors" `
	MaxNumOfVisitors int      `json:"maxNumOfVisitors" `
	ImageIds         []string `json:"imageIds"`
	Rating           float32  `json:"rating"`
}
