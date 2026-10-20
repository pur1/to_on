package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RefreshToken struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"user_id"`
	TokenHash string        `bson:"token_hash"`
	FamilyID  bson.ObjectID `bson:"family_id"` // shared across rotations of one login session
	ExpiresAt time.Time     `bson:"expires_at"`
	RevokedAt *time.Time    `bson:"revoked_at,omitempty"`
	UserAgent string        `bson:"user_agent,omitempty"`
	IP        string        `bson:"ip,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
}
