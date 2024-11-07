package entities

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	Fullname            string             `bson:"fullname" json:"fullname" validate:"required"`
	Username            string             `bson:"username" json:"username" validate:"required"`
	Email               string             `bson:"email" json:"email" validate:"required,email"`
	Password        	string             `bson:"password" json:"password"`
	IsActive            bool               `bson:"isActive" json:"isActive"`
	CreatedAt           primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedBy           []UpdatedBy        `bson:"updatedBy" json:"updatedBy"`
	IsDeleted           bool               `bson:"isDeleted" json:"isDeleted"`
}

type UserRepository interface {
	InsertOne(ctx context.Context, user *User) (*User, error)
	FindOne(ctx context.Context, id string) (*User, error)
	UpdateOne(ctx context.Context, user *User, id string) (*User, error)
	GetByCredential(ctx context.Context, username string, password string) (*User, error)
}

type UserUsecase interface {
	InsertOne(ctx context.Context, user *User) (*User, error)
	FindOne(ctx context.Context, id string) (*User, error)
	UpdateOne(ctx context.Context, user *User, id string) (*User, error)
}
