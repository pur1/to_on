package httpx

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/pur1/to_on/packages/utils/apperror"
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    apperror.Code `json:"code"`
	Message string        `json:"message"`
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		appErr = apperror.Wrap(err, apperror.CodeInternal, "internal error")
	}

	status := statusFor(appErr.Code)
	if status >= 500 {
		slog.ErrorContext(r.Context(), "request failed", "path", r.URL.Path, "err", err)
	}
	JSON(w, status, errorBody{Error: errorDetail{Code: appErr.Code, Message: appErr.Message}})
}

func statusFor(code apperror.Code) int {
	switch code {
	case apperror.CodeInvalid:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperror.CodeForbidden:
		return http.StatusForbidden
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case apperror.CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
