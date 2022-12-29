package handlers

import (
	"encoding/json"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/service"
	"io/ioutil"
	"net/http"
)

type AddUserHandler struct {
	Service service.ServiceInterface
}

func NewAddUserHandler(service service.ServiceInterface) *AddUserHandler {
	return &AddUserHandler{
		Service: service,
	}
}

func (h *AddUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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
	
	if user.Admin == true {
		token := getToken(request)
		currUser, _ := h.Service.GetUserByToken(token)
		if currUser == nil || !currUser.Admin {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	if err = h.Service.AddUser(&user); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json")
		error := entities.Error{Message: err.Error()}
		errJson, _ := json.Marshal(error)
		writer.Write(errJson)
		return
	}

	id := entities.Id{
		Id: user.Id,
	}
	idJson, _ := json.Marshal(id)
	writer.Write(idJson)
	writer.WriteHeader(http.StatusOK)
}
