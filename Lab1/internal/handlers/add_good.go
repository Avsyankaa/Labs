package handlers

import (
	"encoding/json"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/service"
	"io/ioutil"
	"net/http"
)

type AddGoodHandler struct {
	Service service.ServiceInterface
}

func NewAddGoodHandler(service service.ServiceInterface) *AddGoodHandler {
	return &AddGoodHandler{
		Service: service,
	}
}

func (h *AddGoodHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

	var good entities.Good
	err = json.Unmarshal(data, &good)
	if err != nil || !good.Valid() {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if err = h.Service.AddGood(user, &good); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		error := entities.Error{Message: err.Error()}
		errJson, _ := json.Marshal(error)
		writer.Write(errJson)
		return
	}

	id := entities.Id{
		Id: good.BriefInfo.Id,
	}
	idJson, _ := json.Marshal(id)
	writer.Write(idJson)
	writer.WriteHeader(http.StatusOK)
}
