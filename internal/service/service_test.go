package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/SussyaPusya/L0/internal/dto"
	"github.com/SussyaPusya/L0/internal/service"
	"github.com/SussyaPusya/L0/pkg/logger"
	"github.com/SussyaPusya/L0/pkg/test"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder(T *testing.T) {

	repo := service.NewMockRepository(T)

	log := logger.NewMockLogger(T)

	srv := service.NewService(repo, log)

	tests := []*dto.Order{}

	for range 20 {
		tests = append(tests, test.GenerateOrder())
	}

	repo.EXPECT().CreateOrder(mock.Anything, mock.AnythingOfType("*dto.Order")).Return(nil)

	for i, tt := range tests {

		T.Run(fmt.Sprintf("Default data test : %d", i), func(t *testing.T) {
			if err := srv.CreateOrder(context.Background(), tt); err != nil {

				t.Error("Error:", err)
			}

		})
	}

}
