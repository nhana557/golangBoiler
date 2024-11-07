package usecase

import (
	"boiler-go/entities"
	"context"
	"time"
)

type loginUsecase struct {
	userRepo 			entities.UserRepository
	contextTimeout 		time.Duration
}

func NewLoginUseCase(u entities.UserRepository, to time.Duration) entities.LoginUsecase{
	return &loginUsecase{
		userRepo: u,
		contextTimeout: to,
	}
}

func (login *loginUsecase) GetUser(c context.Context, username string, password string) (*entities.User, error) {

	ctx, cancel := context.WithTimeout(c, login.contextTimeout)
	defer cancel()
	

	res, err := login.userRepo.GetByCredential(ctx, username, password)
	if err != nil {
		return res, err
	}

	return res, nil
}