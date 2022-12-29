package service

import (
	"errors"
	"github.com/avsyankaa/auth/internal/crypto"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/avsyankaa/auth/internal/storage"
)

type ServiceInterface interface {
	AddUser(*entities.User) error
	ChangePassword(*entities.UserPasswordChange) error
	Authorize(*entities.User) (*entities.AccessToken, error)
	GetUserByToken(string) (*entities.User, error)

	AddGood(*entities.User, *entities.Good) error
	AddOrder(*entities.User) (int64, error)
	GetGoodsBrief() ([]entities.GoodBrief, error)
	GetGoodFull(int64) (*entities.Good, error)
	AddGoodToOrder(*entities.User, *entities.OrderItem) error
}

type Service struct {
	Storage     storage.StorageInterface
	HashService crypto.HashInterface
}

func NewService(userStorage storage.StorageInterface, hashInterface crypto.HashInterface) *Service {
	return &Service{
		Storage:     userStorage,
		HashService: hashInterface,
	}
}

func (s *Service) AddGoodToOrder(user *entities.User, orderItem *entities.OrderItem) error {
	userId, err := s.Storage.GetOrderUser(orderItem.OrderId)
	if err != nil {
		return err
	}
	if userId != user.Id {
		return errors.New("No access to this action")
	}
	return s.Storage.AddGoodToOrder(orderItem)
}

func (s *Service) AddOrder(user *entities.User) (int64, error) {
	order := entities.Order{
		UserId: user.Id,
		Status: false,
	}
	err := s.Storage.AddOrder(&order)
	if err != nil {
		return 0, err
	}
	return order.Id, nil
}

func (s *Service) GetUserByToken(token string) (*entities.User, error) {
	return s.Storage.GetUserByToken(token)
}

func (as *Service) GetGoodFull(goodId int64) (*entities.Good, error) {
	// Here we will check for access
	return as.Storage.GetGoodFull(goodId)
}

func (as *Service) AddGood(user *entities.User, good *entities.Good) error {
	//if user.Admin == false {
	//return errors.New("No access to this action")
	//	}
	// Here we will check for access
	return as.Storage.AddGood(good)
}

func (as *Service) GetGoodsBrief() ([]entities.GoodBrief, error) {
	// Here we will check for access
	return as.Storage.GetGoodsBrief()
}

func (as *Service) AddUser(user *entities.User) error {
	user.Password = as.HashService.GenerateHash(user.Password)
	return as.Storage.AddUser(user)
}

func (as *Service) ChangePassword(user *entities.UserPasswordChange) error {
	user.Password = as.HashService.GenerateHash(user.Password)
	user.NewPassword = as.HashService.GenerateHash(user.NewPassword)
	return as.Storage.ChangePassword(user)
}

func (as *Service) Authorize(user *entities.User) (*entities.AccessToken, error) {
	user.Password = as.HashService.GenerateHash(user.Password)
	return as.Storage.Authorize(user)
}
