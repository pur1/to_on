package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/pur1/to_on/apps/api/auth/internal/model"
	"github.com/pur1/to_on/apps/api/auth/internal/port"
	"github.com/pur1/to_on/packages/utils/apperror"
	"github.com/pur1/to_on/packages/utils/password"
)

type RegisterInput struct {
	Email    string
	Username string
	Password string
}

type LoginInput struct {
	Identifier string // email or username
	Password   string
}

// ClientInfo is request metadata recorded on refresh tokens.
type ClientInfo struct {
	UserAgent string
	IP        string
}

// Session is what a successful register/login/refresh returns.
type Session struct {
	User         *model.User
	AccessToken  string
	RefreshToken string
	RefreshTTL   time.Duration
}

type AuthService struct {
	users      port.UserRepository
	tokens     port.TokenRepository
	issuer     *TokenIssuer
	refreshTTL time.Duration
	now        func() time.Time
}

func NewAuthService(users port.UserRepository, tokens port.TokenRepository, issuer *TokenIssuer, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		users:      users,
		tokens:     tokens,
		issuer:     issuer,
		refreshTTL: refreshTTL,
		now:        time.Now,
	}
}

var errInvalidCredentials = apperror.New(apperror.CodeUnauthorized, "invalid credentials")
var errInvalidSession = apperror.New(apperror.CodeUnauthorized, "invalid session")

// dummyHash is compared against when the user does not exist so login timing does not reveal accounts.
var dummyHash = sync.OnceValue(func() string {
	h, _ := password.Hash("dummy-password-for-timing")
	return h
})

func (s *AuthService) Register(ctx context.Context, in RegisterInput, client ClientInfo) (*Session, error) {
	in.Email = normalizeEmail(in.Email)
	in.Username = normalizeUsername(in.Username)
	if err := validateRegister(in); err != nil {
		return nil, err
	}

	hash, err := password.Hash(in.Password)
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "hash password")
	}

	now := s.now()
	u := &model.User{
		Email:        in.Email,
		Username:     in.Username,
		PasswordHash: hash,
		Role:         model.RoleReader,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(ctx, u); err != nil {
		if errors.Is(err, apperror.ErrDuplicate) {
			return nil, apperror.New(apperror.CodeConflict, "email or username already taken")
		}
		return nil, apperror.Wrap(err, apperror.CodeInternal, "create user")
	}

	return s.startSession(ctx, u, bson.NewObjectID(), client)
}

func (s *AuthService) Login(ctx context.Context, in LoginInput, client ClientInfo) (*Session, error) {
	id := strings.TrimSpace(in.Identifier)
	if id == "" || in.Password == "" {
		return nil, apperror.New(apperror.CodeInvalid, "identifier and password are required")
	}

	var u *model.User
	var err error
	if strings.Contains(id, "@") {
		u, err = s.users.FindByEmail(ctx, normalizeEmail(id))
	} else {
		u, err = s.users.FindByUsername(ctx, normalizeUsername(id))
	}
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "find user")
	}

	if u == nil || u.PasswordHash == "" {
		password.Compare(dummyHash(), in.Password)
		return nil, errInvalidCredentials
	}
	if !password.Compare(u.PasswordHash, in.Password) {
		return nil, errInvalidCredentials
	}

	return s.startSession(ctx, u, bson.NewObjectID(), client)
}

// Refresh rotates the refresh token. Reusing an already-rotated token revokes the whole family.
func (s *AuthService) Refresh(ctx context.Context, rawRefresh string, client ClientInfo) (*Session, error) {
	if rawRefresh == "" {
		return nil, errInvalidSession
	}
	t, err := s.tokens.FindByHash(ctx, hashToken(rawRefresh))
	if errors.Is(err, apperror.ErrNotFound) {
		return nil, errInvalidSession
	}
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "find refresh token")
	}

	now := s.now()
	if t.RevokedAt != nil {
		if err := s.tokens.RevokeFamily(ctx, t.FamilyID, now); err != nil {
			return nil, apperror.Wrap(err, apperror.CodeInternal, "revoke token family")
		}
		return nil, errInvalidSession
	}
	if !now.Before(t.ExpiresAt) {
		return nil, errInvalidSession
	}

	revoked, err := s.tokens.Revoke(ctx, t.ID, now)
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "revoke refresh token")
	}
	if !revoked {
		// Lost a race with another refresh using the same token: treat as reuse.
		if err := s.tokens.RevokeFamily(ctx, t.FamilyID, now); err != nil {
			return nil, apperror.Wrap(err, apperror.CodeInternal, "revoke token family")
		}
		return nil, errInvalidSession
	}

	u, err := s.users.FindByID(ctx, t.UserID)
	if errors.Is(err, apperror.ErrNotFound) {
		return nil, errInvalidSession
	}
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "find user")
	}

	return s.startSession(ctx, u, t.FamilyID, client)
}

// Logout revokes every token in the session's family. Unknown tokens are ignored.
func (s *AuthService) Logout(ctx context.Context, rawRefresh string) error {
	if rawRefresh == "" {
		return nil
	}
	t, err := s.tokens.FindByHash(ctx, hashToken(rawRefresh))
	if errors.Is(err, apperror.ErrNotFound) {
		return nil
	}
	if err != nil {
		return apperror.Wrap(err, apperror.CodeInternal, "find refresh token")
	}
	if err := s.tokens.RevokeFamily(ctx, t.FamilyID, s.now()); err != nil {
		return apperror.Wrap(err, apperror.CodeInternal, "revoke token family")
	}
	return nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*model.User, error) {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errInvalidSession
	}
	u, err := s.users.FindByID(ctx, id)
	if errors.Is(err, apperror.ErrNotFound) {
		return nil, apperror.New(apperror.CodeNotFound, "user not found")
	}
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "find user")
	}
	return u, nil
}

// VerifyAccess returns the user ID for a valid access token.
func (s *AuthService) VerifyAccess(raw string) (string, error) {
	sub, err := s.issuer.ParseAccess(raw)
	if err != nil {
		return "", apperror.Wrap(err, apperror.CodeUnauthorized, "invalid access token")
	}
	return sub, nil
}

func (s *AuthService) startSession(ctx context.Context, u *model.User, familyID bson.ObjectID, client ClientInfo) (*Session, error) {
	now := s.now()
	access, err := s.issuer.IssueAccess(u.ID.Hex(), now)
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "issue access token")
	}
	raw, hash, err := newRefreshToken()
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "generate refresh token")
	}
	err = s.tokens.Create(ctx, &model.RefreshToken{
		UserID:    u.ID,
		TokenHash: hash,
		FamilyID:  familyID,
		ExpiresAt: now.Add(s.refreshTTL),
		UserAgent: client.UserAgent,
		IP:        client.IP,
		CreatedAt: now,
	})
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "store refresh token")
	}
	return &Session{User: u, AccessToken: access, RefreshToken: raw, RefreshTTL: s.refreshTTL}, nil
}
