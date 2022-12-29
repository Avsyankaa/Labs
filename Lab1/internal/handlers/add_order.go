package handlers

import (
	"encoding/json"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/service"
	"net/http"
)

type AddOrderHandler struct {
	Service service.ServiceInterface
}

func NewAddOrderHandler(service service.ServiceInterface) *AddOrderHandler {
	return &AddOrderHandler{
		Service: service,
	}
}

func (h *AddOrderHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	token := getToken(request)
	user, _ := h.Service.GetUserByToken(token)
	if user == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := h.Service.AddOrder(user)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		error := entities.Error{Message: err.Error()}
		errJson, _ := json.Marshal(error)
		writer.Write(errJson)
		return
	}

	idJson, _ := json.Marshal(id)
	writer.Write(idJson)
	writer.WriteHeader(http.StatusOK)
}
