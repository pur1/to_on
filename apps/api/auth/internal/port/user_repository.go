package port

import (
	"context"

	"github.com/pur1/to_on/apps/api/auth/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error // sets u.ID; ErrDuplicate on email/username clash
	FindByID(ctx context.Context, id bson.ObjectID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}
