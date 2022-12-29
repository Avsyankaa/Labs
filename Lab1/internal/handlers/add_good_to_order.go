package handlers

import (
	"encoding/json"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/service"
	"io/ioutil"
	"net/http"
)

type AddGoodOrderHandler struct {
	Service service.ServiceInterface
}

func NewAddGoodOrderHandler(service service.ServiceInterface) *AddGoodOrderHandler {
	return &AddGoodOrderHandler{
		Service: service,
	}
}

func (h *AddGoodOrderHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	token := getToken(request)
	user, _ := h.Service.GetUserByToken(token)
	if user == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var orderItem entities.OrderItem
	err = json.Unmarshal(data, &orderItem)
	if err != nil || !orderItem.Valid() {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if err = h.Service.AddGoodToOrder(user, &orderItem); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		error := entities.Error{Message: err.Error()}
		errJson, _ := json.Marshal(error)
		writer.Write(errJson)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
