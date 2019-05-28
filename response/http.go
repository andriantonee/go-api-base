package response

type HttpError struct {
	Code   int    `json:"code"`
	Errors Errors `json:"errors"`
}

func NewHttpError(code int, message string) HttpError {
	err := NewErrors(message)
	return HttpError{
		Code:   code,
		Errors: err,
	}
}

func NewHttpErrorWithDetail(
	code int,
	message string,
	detail []string,
) HttpError {
	err := NewErrorsWithDetail(message, detail)
	return HttpError{
		Code:   code,
		Errors: err,
	}
}
