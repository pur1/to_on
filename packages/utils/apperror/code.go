package apperror

type Code string

const (
	CodeInvalid      Code = "invalid"
	CodeNotFound     Code = "not_found"
	CodeUnauthorized Code = "unauthorized"
	CodeConflict     Code = "conflict"
	CodeForbidden    Code = "forbidden"
	CodeInternal     Code = "internal"
)
