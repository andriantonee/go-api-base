package response

type Errors struct {
	Message string   `json:"message"`
	Detail  []string `json:"detail,omitempty"`
}

func NewErrors(message string) Errors {
	return Errors{
		Message: message,
	}
}

func NewErrorsWithDetail(message string, detail []string) Errors {
	return Errors{
		Message: message,
		Detail:  detail,
	}
}
