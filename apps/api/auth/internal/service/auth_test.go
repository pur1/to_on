package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/pur1/to_on/apps/api/auth/internal/model"
	"github.com/pur1/to_on/packages/utils/apperror"
)

type fakeUsers struct {
	mu   sync.Mutex
	byID map[bson.ObjectID]*model.User
}

func (f *fakeUsers) Create(_ context.Context, u *model.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, existing := range f.byID {
		if existing.Email == u.Email || existing.Username == u.Username {
			return apperror.ErrDuplicate
		}
	}
	u.ID = bson.NewObjectID()
	cp := *u
	f.byID[u.ID] = &cp
	return nil
}

func (f *fakeUsers) FindByID(_ context.Context, id bson.ObjectID) (*model.User, error) {
	return f.find(func(u *model.User) bool { return u.ID == id })
}

func (f *fakeUsers) FindByEmail(_ context.Context, email string) (*model.User, error) {
	return f.find(func(u *model.User) bool { return u.Email == email })
}

func (f *fakeUsers) FindByUsername(_ context.Context, username string) (*model.User, error) {
	return f.find(func(u *model.User) bool { return u.Username == username })
}

func (f *fakeUsers) find(match func(*model.User) bool) (*model.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.byID {
		if match(u) {
			cp := *u
			return &cp, nil
		}
	}
	return nil, apperror.ErrNotFound
}

type fakeTokens struct {
	mu   sync.Mutex
	byID map[bson.ObjectID]*model.RefreshToken
}

func (f *fakeTokens) Create(_ context.Context, t *model.RefreshToken) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	t.ID = bson.NewObjectID()
	cp := *t
	f.byID[t.ID] = &cp
	return nil
}

func (f *fakeTokens) FindByHash(_ context.Context, hash string) (*model.RefreshToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, t := range f.byID {
		if t.TokenHash == hash {
			cp := *t
			return &cp, nil
		}
	}
	return nil, apperror.ErrNotFound
}

func (f *fakeTokens) Revoke(_ context.Context, id bson.ObjectID, at time.Time) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	t, ok := f.byID[id]
	if !ok || t.RevokedAt != nil {
		return false, nil
	}
	t.RevokedAt = &at
	return true, nil
}

func (f *fakeTokens) RevokeFamily(_ context.Context, familyID bson.ObjectID, at time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, t := range f.byID {
		if t.FamilyID == familyID && t.RevokedAt == nil {
			t.RevokedAt = &at
		}
	}
	return nil
}

func newTestService() *AuthService {
	return NewAuthService(
		&fakeUsers{byID: map[bson.ObjectID]*model.User{}},
		&fakeTokens{byID: map[bson.ObjectID]*model.RefreshToken{}},
		NewTokenIssuer([]byte("test-secret-test-secret-test-secret"), 15*time.Minute),
		24*time.Hour,
	)
}

func assertCode(t *testing.T, err error, want apperror.Code) {
	t.Helper()
	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Code != want {
		t.Fatalf("got err %v, want code %s", err, want)
	}
}

var validInput = RegisterInput{Email: "Alice@Example.com", Username: "Alice_1", Password: "password123"}

func TestRegister(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	sess, err := svc.Register(ctx, validInput, ClientInfo{})
	if err != nil {
		t.Fatal(err)
	}
	if sess.User.Email != "alice@example.com" || sess.User.Username != "alice_1" {
		t.Errorf("not normalized: %+v", sess.User)
	}
	if sess.AccessToken == "" || sess.RefreshToken == "" {
		t.Error("missing tokens")
	}
	if sub, err := svc.VerifyAccess(sess.AccessToken); err != nil || sub != sess.User.ID.Hex() {
		t.Errorf("access token sub=%q err=%v", sub, err)
	}

	_, err = svc.Register(ctx, validInput, ClientInfo{})
	assertCode(t, err, apperror.CodeConflict)
}

func TestRegisterValidation(t *testing.T) {
	cases := map[string]RegisterInput{
		"bad email":      {Email: "nope", Username: "alice", Password: "password123"},
		"short username": {Email: "a@b.co", Username: "al", Password: "password123"},
		"bad username":   {Email: "a@b.co", Username: "al ice", Password: "password123"},
		"short password": {Email: "a@b.co", Username: "alice", Password: "short"},
	}
	svc := newTestService()
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := svc.Register(context.Background(), in, ClientInfo{})
			assertCode(t, err, apperror.CodeInvalid)
		})
	}
}

func TestLogin(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, validInput, ClientInfo{}); err != nil {
		t.Fatal(err)
	}

	for _, id := range []string{"alice@example.com", "ALICE_1"} {
		if _, err := svc.Login(ctx, LoginInput{Identifier: id, Password: "password123"}, ClientInfo{}); err != nil {
			t.Errorf("login with %q: %v", id, err)
		}
	}

	_, err := svc.Login(ctx, LoginInput{Identifier: "alice_1", Password: "wrong-password"}, ClientInfo{})
	assertCode(t, err, apperror.CodeUnauthorized)
	_, err = svc.Login(ctx, LoginInput{Identifier: "nobody", Password: "password123"}, ClientInfo{})
	assertCode(t, err, apperror.CodeUnauthorized)
}

func TestRefreshRotationAndReuse(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	first, err := svc.Register(ctx, validInput, ClientInfo{})
	if err != nil {
		t.Fatal(err)
	}

	second, err := svc.Refresh(ctx, first.RefreshToken, ClientInfo{})
	if err != nil {
		t.Fatal(err)
	}
	if second.RefreshToken == first.RefreshToken {
		t.Fatal("refresh token not rotated")
	}

	// Replaying the rotated token is treated as theft: it fails and kills the family.
	_, err = svc.Refresh(ctx, first.RefreshToken, ClientInfo{})
	assertCode(t, err, apperror.CodeUnauthorized)
	_, err = svc.Refresh(ctx, second.RefreshToken, ClientInfo{})
	assertCode(t, err, apperror.CodeUnauthorized)
}

func TestRefreshExpired(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	sess, err := svc.Register(ctx, validInput, ClientInfo{})
	if err != nil {
		t.Fatal(err)
	}
	svc.now = func() time.Time { return time.Now().Add(25 * time.Hour) }
	_, err = svc.Refresh(ctx, sess.RefreshToken, ClientInfo{})
	assertCode(t, err, apperror.CodeUnauthorized)
}

func TestLogout(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	sess, err := svc.Register(ctx, validInput, ClientInfo{})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Logout(ctx, sess.RefreshToken); err != nil {
		t.Fatal(err)
	}
	_, err = svc.Refresh(ctx, sess.RefreshToken, ClientInfo{})
	assertCode(t, err, apperror.CodeUnauthorized)

	if err := svc.Logout(ctx, "unknown"); err != nil {
		t.Errorf("logout with unknown token: %v", err)
	}
}

func TestMe(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	sess, err := svc.Register(ctx, validInput, ClientInfo{})
	if err != nil {
		t.Fatal(err)
	}
	u, err := svc.Me(ctx, sess.User.ID.Hex())
	if err != nil || u.ID != sess.User.ID {
		t.Fatalf("me: %+v %v", u, err)
	}
}
