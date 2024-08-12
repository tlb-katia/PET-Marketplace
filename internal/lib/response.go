package lib

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func RespError(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func RespOK() Response {
	return Response{
		Status: StatusOK,
	}
}
