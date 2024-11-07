package usecase

import (
	"boiler-go/entities"
	"context"
	"time"

	"github.com/spf13/viper"
)

type JwtUsecase struct {
	UserRepo       entities.UserRepository
	ContextTimeout time.Duration
	Config         *viper.Viper
}

func NewJwtUsecase(u entities.UserRepository, to time.Duration, config *viper.Viper) entities.JwtUsecase {
	return &JwtUsecase{
		UserRepo:       u,
		ContextTimeout: to,
		Config:         config,
	}
}

func (h *JwtUsecase) getOneUser(c context.Context, id string) (*entities.User, error) {

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	res, err := h.UserRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}