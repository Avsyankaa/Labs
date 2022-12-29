package handlers

import (
	"net/http"
	"strings"
)

func getToken(request *http.Request) string {
	token := request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") || strings.HasPrefix(token, "bearer ") {
		token = token[7:]
	}
	return token
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
