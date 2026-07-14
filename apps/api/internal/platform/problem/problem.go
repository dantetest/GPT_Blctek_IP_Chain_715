package problem

import "net/http"

type Problem struct {
	Status  int
	Code    string
	Message string
}

func (p Problem) Error() string { return p.Code + ": " + p.Message }

func InvalidArgument(code, message string) Problem {
	return Problem{Status: http.StatusBadRequest, Code: code, Message: message}
}

func Conflict(code, message string) Problem {
	return Problem{Status: http.StatusConflict, Code: code, Message: message}
}

func Internal() Problem {
	return Problem{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "an internal error occurred"}
}
