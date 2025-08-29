package service

import (
	"fmt"

	"github.com/SussyaPusya/L0/internal/dto"
)

type Repository interface {
}

type Service struct {
	repo Repository
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateOrder(order *dto.Order) error {
	fmt.Println(order)
	return nil
}
