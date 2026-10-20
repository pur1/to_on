package port

import (
	"context"
	"time"

	"github.com/pur1/to_on/apps/api/auth/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TokenRepository interface {
	Create(ctx context.Context, t *model.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*model.RefreshToken, error)
	// Revoke marks the token revoked only if it is still active; returns false if it was already revoked.
	Revoke(ctx context.Context, id bson.ObjectID, at time.Time) (bool, error)
	RevokeFamily(ctx context.Context, familyID bson.ObjectID, at time.Time) error
}
