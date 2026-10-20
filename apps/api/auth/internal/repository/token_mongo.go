package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/pur1/to_on/apps/api/auth/internal/model"
	"github.com/pur1/to_on/packages/utils/apperror"
)

type MongoTokenRepository struct {
	col *mongo.Collection
}

func NewMongoTokenRepository(db *mongo.Database) *MongoTokenRepository {
	return &MongoTokenRepository{col: db.Collection("refresh_tokens")}
}

func (r *MongoTokenRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "token_hash", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "family_id", Value: 1}}},
		{Keys: bson.D{{Key: "expires_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)},
	})
	return err
}

func (r *MongoTokenRepository) Create(ctx context.Context, t *model.RefreshToken) error {
	res, err := r.col.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	t.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *MongoTokenRepository) FindByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	var t model.RefreshToken
	err := r.col.FindOne(ctx, bson.D{{Key: "token_hash", Value: hash}}).Decode(&t)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperror.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *MongoTokenRepository) Revoke(ctx context.Context, id bson.ObjectID, at time.Time) (bool, error) {
	res, err := r.col.UpdateOne(ctx,
		bson.D{{Key: "_id", Value: id}, {Key: "revoked_at", Value: bson.D{{Key: "$exists", Value: false}}}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "revoked_at", Value: at}}}},
	)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount == 1, nil
}

func (r *MongoTokenRepository) RevokeFamily(ctx context.Context, familyID bson.ObjectID, at time.Time) error {
	_, err := r.col.UpdateMany(ctx,
		bson.D{{Key: "family_id", Value: familyID}, {Key: "revoked_at", Value: bson.D{{Key: "$exists", Value: false}}}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "revoked_at", Value: at}}}},
	)
	return err
}
