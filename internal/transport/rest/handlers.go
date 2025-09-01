package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SussyaPusya/L0/internal/dto"
)

type Service interface {
	GetOrder(ctx context.Context, orderID string) (*dto.Order, error)
}

type Handlers struct {
	service Service
}

func NewHandlers(s Service) *Handlers {

	return &Handlers{service: s}

}

func (h *Handlers) GetOrder(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path

	orderID, _ := strings.CutPrefix(url, "/orders/")

	fmt.Println(orderID)
	if orderID == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
