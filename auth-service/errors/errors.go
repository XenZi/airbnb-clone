package errors

type ErrorStruct struct {
	message string
	status  int
}

func NewError(message string, status int) *ErrorStruct {
	return &ErrorStruct{
		message: message,
		status:  status,
	}
}

func (e *ErrorStruct) GetErrorMessage() string {
	return e.message
}

func (e *ErrorStruct) GetErrorStatus() int {
	return e.status
}
