package service

import (
	"context"

	"github.com/SussyaPusya/L0/internal/dto"
	validate "github.com/SussyaPusya/L0/pkg/validator"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *dto.Order) error

	GetOrder(ctx context.Context, orderID string) (*dto.Order, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateOrder(ctx context.Context, order *dto.Order) error {

	err := validate.Order(order)
	if err != nil {
		//логи

		return err
	}

	if err = s.repo.CreateOrder(ctx, order); err != nil {
		// логи
		return err
	}

	return nil
}

func (s *service) GetOrder(ctx context.Context, orderID string) (*dto.Order, error) {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		// логи
		return nil, err
	}
	return order, nil
}
