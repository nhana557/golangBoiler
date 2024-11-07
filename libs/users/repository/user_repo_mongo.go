package repository

import (
	"boiler-go/database/mongo"
	"boiler-go/entities"
	"boiler-go/utils"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoRepository struct {
	DB 				mongo.Database
	Collection  	mongo.Collection
}

const collectionName = "users"

func NewMongoRepository(DB mongo.Database) entities.UserRepository {
	return &mongoRepository{DB, DB.Collection(collectionName)}
}

func (m *mongoRepository) InsertOne(ctx context.Context, user *entities.User) (*entities.User, error) {
	var (
		err error
	)
	_, err = m.Collection.InsertOne(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *mongoRepository) FindOne(ctx context.Context, id string) (*entities.User, error) {
	var (
		user entities.User
		err  error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &user, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (m *mongoRepository) UpdateOne(ctx context.Context, user *entities.User, id string) (*entities.User, error) {
	var (
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	password, err := utils.HashPassword(user.Password)
	if err != nil {
		return user, err
	}
	// Set the update timestamp
	currentTime := time.Now()
	fullname := user.Fullname
	description := "updated user data"

	// Create an entry for the UpdatedBy field
	updatedByEntry := entities.UpdatedBy{
		UserId:      &idHex,
		Name:        &fullname,
		Date:        &currentTime,
		Description: &description,
	}

	// Define the filter and update operations
	filter := bson.M{"_id": idHex}
	update := bson.M{
		"$set": bson.M{
			"fullname":  user.Fullname,
			"username":  user.Username,
			"password":  password,
		},
		"$push": bson.M{
			"updatedBy": updatedByEntry, 
		},
	}
	_, err = m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return user, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *mongoRepository) GetByCredential(ctx context.Context, username string, password string) (*entities.User, error) {
	var (
		user entities.User
		err  error
	)
	
	filter := bson.M{"username": username}
	err = m.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		
		return nil, err 
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("incorrect password")
	}


	return &user, nil
}
