package usecase

import (
	"boiler-go/entities"
	"boiler-go/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userRepo       entities.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(ur entities.UserRepository, to time.Duration) entities.UserUsecase {
	return &userUsecase{
		userRepo:       ur,
		contextTimeout: to,
	}
}

func (user *userUsecase) InsertOne(c context.Context, m *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(c, user.contextTimeout)
	defer cancel()

	hashedPassword, err := utils.HashPassword(m.Password)
	if err != nil {
		return nil, err
	}
	m.Password = string(hashedPassword)

	m.ID = primitive.NewObjectID()
	m.IsActive = true
	m.IsDeleted = false
	m.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	m.UpdatedBy = []entities.UpdatedBy{}

	res, err := user.userRepo.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (user *userUsecase) FindOne(c context.Context, id string) (*entities.User, error) {

	ctx, cancel := context.WithTimeout(c, user.contextTimeout)
	defer cancel()

	res, err := user.userRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (user *userUsecase) UpdateOne(c context.Context, m *entities.User, id string) (*entities.User, error) {

	ctx, cancel := context.WithTimeout(c, user.contextTimeout)
	defer cancel()

	res, err := user.userRepo.UpdateOne(ctx, m, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
