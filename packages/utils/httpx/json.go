package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/pur1/to_on/packages/utils/apperror"
)

const maxBodyBytes = 1 << 20

func Decode(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return apperror.Wrap(err, apperror.CodeInvalid, "invalid request body")
	}
	return nil
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
