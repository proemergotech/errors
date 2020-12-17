package errors

const (
	ErrCode          = "code"
	ErrDetails       = "details"
	ErrPublicDetails = "public_details"
	ErrHTTPCode      = "http_code"
)

type HTTPError struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string                   `json:"message,omitempty" validate:"omitempty"`
	Code    string                   `json:"code" validate:"required"`
	Details []map[string]interface{} `json:"details,omitempty" validate:"omitempty"`
}

func ToHTTPError(err error) (*HTTPError, int) {
	httpCode := HTTPCode(err)
	if httpCode == 0 {
		httpCode = 500
	}

	return &HTTPError{
		Error: Error{
			Message: err.Error(),
			Code:    Code(err),
			Details: Details(err),
		},
	}, httpCode
}

func ToPublicHTTPError(err error) (*HTTPError, int) {
	httpCode := HTTPCode(err)
	if httpCode == 0 {
		httpCode = 500
	}

	code := Code(err)
	if httpCode >= 500 {
		code = "ERR_INTERNAL"
	}

	return &HTTPError{
		Error: Error{
			Code:    code,
			Details: PublicDetails(err),
		},
	}, httpCode
}

func Code(err error) string {
	field := Field(err, ErrCode)
	if field == nil {
		return ""
	}

	return field.(string)
}

func HTTPCode(err error) int {
	field := Field(err, ErrHTTPCode)
	if field == nil {
		return 0
	}

	return field.(int)
}

func Details(err error) []map[string]interface{} {
	field := Field(err, ErrDetails)
	if field == nil {
		return []map[string]interface{}{}
	}

	return field.([]map[string]interface{})
}

func PublicDetails(err error) []map[string]interface{} {
	field := Field(err, ErrPublicDetails)
	if field == nil {
		return []map[string]interface{}{}
	}

	return field.([]map[string]interface{})
}
