package handlers

import (
	"encoding/json"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/service"
	"io/ioutil"
	"net/http"
)

type ChangePasswordHandler struct {
	Service service.ServiceInterface
}

func NewChangePasswordHandler(service service.ServiceInterface) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		Service: service,
	}
}

func (h *ChangePasswordHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var changePassword entities.UserPasswordChange
	err = json.Unmarshal(data, &changePassword)
	if err != nil || !changePassword.Valid() {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if err = h.Service.ChangePassword(&changePassword); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		error := entities.Error{Message: err.Error()}
		errJson, _ := json.Marshal(error)
		writer.Write(errJson)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
