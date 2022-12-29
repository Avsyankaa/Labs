package handlers

import (
	"encoding/json"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/service"
	"io/ioutil"
	"net/http"
)

type AuthHandler struct {
	Service service.ServiceInterface
}

func NewAuthHandler(service service.ServiceInterface) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

func (auh *AuthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var user entities.User
	err = json.Unmarshal(data, &user)
	if err != nil || !user.Valid() {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := auh.Service.Authorize(&user)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		error := entities.Error{Message: err.Error()}
		errJson, _ := json.Marshal(error)
		writer.Write(errJson)
		return
	}
	if token == nil {
		writer.WriteHeader(http.StatusUnauthorized)
	}

	tokenJson, _ := json.Marshal(token)
	writer.Write(tokenJson)
	writer.WriteHeader(http.StatusOK)
}
