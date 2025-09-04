package service

import (
	"context"

	"github.com/SussyaPusya/L0/internal/dto"
	"github.com/SussyaPusya/L0/pkg/logger"
	validate "github.com/SussyaPusya/L0/pkg/validator"
	"go.uber.org/zap"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *dto.Order) error

	GetOrder(ctx context.Context, orderID string) (*dto.Order, error)
}

type service struct {
	loger *logger.Logger
	repo  Repository
}

func NewService(repo Repository, l *logger.Logger) *service {
	return &service{repo: repo, loger: l}
}

func (s *service) CreateOrder(ctx context.Context, order *dto.Order) error {

	err := validate.Order(order)
	if err != nil {
		s.loger.Error("validation error", zap.Error(err))

		return err
	}

	if err = s.repo.CreateOrder(ctx, order); err != nil {
		// логи
		s.loger.Error("failed to create order", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) GetOrder(ctx context.Context, orderID string) (*dto.Order, error) {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		// логи
		s.loger.Error("failed to get order", zap.Error(err))
		return nil, err
	}
	return order, nil
}
