package validate

import (
	"github.com/SussyaPusya/L0/internal/dto"
	"github.com/go-playground/validator/v10"
)

func Order(order *dto.Order) error {
	valid := validator.New()

	return valid.Struct(order)
}
