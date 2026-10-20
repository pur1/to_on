package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

	"github.com/pur1/to_on/apps/api/auth/internal/model"
	"github.com/pur1/to_on/apps/api/auth/internal/service"
	"github.com/pur1/to_on/packages/utils/httpx"
)

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/api/v1/auth"
)

type AuthHandler struct {
	svc          *service.AuthService
	cookieSecure bool
}

func NewAuthHandler(svc *service.AuthService, cookieSecure bool) *AuthHandler {
	return &AuthHandler{svc: svc, cookieSecure: cookieSecure}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitBy(10, time.Minute, keyByClientIP))
		r.Post("/register", h.register)
		r.Post("/login", h.login)
	})
	r.Post("/refresh", h.refresh)
	r.Post("/logout", h.logout)
	r.With(RequireAuth(h.svc)).Get("/me", h.me)
	return r
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type userResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	Role          model.Role `json:"role"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
}

type sessionResponse struct {
	User        userResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}

func toUserResponse(u *model.User) userResponse {
	role := u.Role
	if role == "" {
		role = model.RoleReader
	}
	return userResponse{
		ID:            u.ID.Hex(),
		Email:         u.Email,
		Username:      u.Username,
		Role:          role,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
	}
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httpx.Decode(w, r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	sess, err := h.svc.Register(r.Context(), service.RegisterInput(req), clientInfo(r))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	h.writeSession(w, http.StatusCreated, sess)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.Decode(w, r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	sess, err := h.svc.Login(r.Context(), service.LoginInput(req), clientInfo(r))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	h.writeSession(w, http.StatusOK, sess)
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	sess, err := h.svc.Refresh(r.Context(), readRefreshCookie(r), clientInfo(r))
	if err != nil {
		h.clearRefreshCookie(w)
		httpx.Error(w, r, err)
		return
	}
	h.writeSession(w, http.StatusOK, sess)
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Logout(r.Context(), readRefreshCookie(r)); err != nil {
		httpx.Error(w, r, err)
		return
	}
	h.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) me(w http.ResponseWriter, r *http.Request) {
	u, err := h.svc.Me(r.Context(), UserIDFrom(r.Context()))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]userResponse{"user": toUserResponse(u)})
}

func (h *AuthHandler) writeSession(w http.ResponseWriter, status int, sess *service.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    sess.RefreshToken,
		Path:     refreshCookiePath,
		MaxAge:   int(sess.RefreshTTL.Seconds()),
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
	httpx.JSON(w, status, sessionResponse{User: toUserResponse(sess.User), AccessToken: sess.AccessToken})
}

func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

func readRefreshCookie(r *http.Request) string {
	c, err := r.Cookie(refreshCookieName)
	if err != nil {
		return ""
	}
	return c.Value
}

func clientInfo(r *http.Request) service.ClientInfo {
	return service.ClientInfo{UserAgent: r.UserAgent(), IP: middleware.GetClientIP(r.Context())}
}

func keyByClientIP(r *http.Request) (string, error) {
	return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
}
