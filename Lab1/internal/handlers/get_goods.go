package handlers

import (
	"encoding/json"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/service"
	"net/http"
)

type GetGoodsHandler struct {
	Service service.ServiceInterface
}

func NewGetGoodsHandler(service service.ServiceInterface) *GetGoodsHandler {
	return &GetGoodsHandler{
		Service: service,
	}
}

func (h *GetGoodsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	goods, err := h.Service.GetGoodsBrief()
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		error := entities.Error{Message: err.Error()}
		errJson, _ := json.Marshal(error)
		writer.Write(errJson)
		return
	}

	writer.WriteHeader(http.StatusOK)
	goodsJson, _ := json.Marshal(goods)
	writer.Write(goodsJson)
	return
}
