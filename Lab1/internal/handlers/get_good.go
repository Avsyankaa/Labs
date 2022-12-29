package handlers

import (
	"encoding/json"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/service"
	"net/http"
	"strconv"
)

type GetGoodHandler struct {
	Service service.ServiceInterface
}

func NewGetGoodHandler(service service.ServiceInterface) *GetGoodHandler {
	return &GetGoodHandler{
		Service: service,
	}
}

func (h *GetGoodHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	goodId, err := strconv.ParseInt(request.Form.Get("good_id"), 10, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	good, err := h.Service.GetGoodFull(goodId)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		error := entities.Error{Message: err.Error()}
		errJson, _ := json.Marshal(error)
		writer.Write(errJson)
		return
	}

	goodJson, _ := json.Marshal(good)
	writer.Write(goodJson)
	writer.WriteHeader(http.StatusOK)
}
