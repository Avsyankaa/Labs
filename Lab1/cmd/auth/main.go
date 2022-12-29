package main

import (
	"github.com/avsyankaa/auth/internal/crypto"
	"github.com/avsyankaa/auth/internal/handlers"
	"github.com/avsyankaa/auth/internal/service"
	"github.com/avsyankaa/auth/internal/storage"
	"net/http"
)

func main() {
	authStorage := storage.NewStorage()
	if authStorage == nil {
		return
	}

	hashService := crypto.SHA256Hash{}

	service := service.NewService(authStorage, &hashService)
	if service == nil {
		return
	}

	addGoodHandler := handlers.NewAddGoodHandler(service)
	addUserHandler := handlers.NewAddUserHandler(service)
	getGoodsHandler := handlers.NewGetGoodsHandler(service)
	getGoodHandler := handlers.NewGetGoodHandler(service)
	addOrderHandler := handlers.NewAddOrderHandler(service)
	authHandler := handlers.NewAuthHandler(service)
	addGoodOrderHandler := handlers.NewAddGoodOrderHandler(service)
	changePasswordHandler := handlers.NewChangePasswordHandler(service)
	if addUserHandler == nil || addGoodHandler == nil || getGoodsHandler == nil {
		return
	}

	http.Handle("/add_user", addUserHandler)
	http.Handle("/add_good", addGoodHandler)
	http.Handle("/get_goods", getGoodsHandler)
	http.Handle("/get_good", getGoodHandler)
	http.Handle("/add_order", addOrderHandler)
	http.Handle("/authorize", authHandler)
	http.Handle("/add_order_good", addGoodOrderHandler)
	http.Handle("/change_password", changePasswordHandler)

	http.ListenAndServe(":8080", nil)
}
