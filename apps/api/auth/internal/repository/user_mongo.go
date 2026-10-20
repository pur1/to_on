package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/pur1/to_on/apps/api/auth/internal/model"
	"github.com/pur1/to_on/packages/utils/apperror"
)

type MongoUserRepository struct {
	col *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{col: db.Collection("users")}
}

func (r *MongoUserRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	return err
}

func (r *MongoUserRepository) Create(ctx context.Context, u *model.User) error {
	res, err := r.col.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return apperror.ErrDuplicate
	}
	if err != nil {
		return err
	}
	u.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	return r.findOne(ctx, bson.D{{Key: "_id", Value: id}})
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.findOne(ctx, bson.D{{Key: "email", Value: email}})
}

func (r *MongoUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.findOne(ctx, bson.D{{Key: "username", Value: username}})
}

func (r *MongoUserRepository) findOne(ctx context.Context, filter bson.D) (*model.User, error) {
	var u model.User
	err := r.col.FindOne(ctx, filter).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperror.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
