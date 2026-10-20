package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string

const (
	RoleReader  Role = "reader"
	RoleCreator Role = "creator"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	Email         string        `bson:"email"`
	Username      string        `bson:"username"`
	PasswordHash  string        `bson:"password_hash,omitempty"` // empty for OAuth-only users
	Role          Role          `bson:"role,omitempty"`
	EmailVerified bool          `bson:"email_verified"`
	CreatedAt     time.Time     `bson:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at"`
}
