package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/pur1/to_on/apps/api/auth/internal/service"
	"github.com/pur1/to_on/packages/utils/apperror"
	"github.com/pur1/to_on/packages/utils/httpx"
)

type ctxKey struct{}

func RequireAuth(svc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || raw == "" {
				httpx.Error(w, r, apperror.New(apperror.CodeUnauthorized, "missing bearer token"))
				return
			}
			userID, err := svc.VerifyAccess(raw)
			if err != nil {
				httpx.Error(w, r, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, userID)))
		})
	}
}

// UserIDFrom returns the authenticated user ID set by RequireAuth, or "".
func UserIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}
