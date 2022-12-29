package storage

import (
	"errors"
	"fmt"
	"github.com/avsyankaa/auth/internal/entities"
	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type StorageInterface interface {
	AddUser(*entities.User) error
	Authorize(*entities.User) (*entities.AccessToken, error)
	GetUserByToken(string) (*entities.User, error)
	ChangePassword(*entities.UserPasswordChange) error

	AddGood(*entities.Good) error
	AddOrder(*entities.Order) error
	GetGoodsBrief() ([]entities.GoodBrief, error)
	GetGoodFull(int64) (*entities.Good, error)
	AddGoodToOrder(*entities.OrderItem) error
	GetOrderUser(int64) (int64, error)
}

type Storage struct {
	Session *dbr.Session
}

func (us *Storage) GetUserByToken(token string) (*entities.User, error) {
	tx, err := us.Session.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var user entities.User
	err = tx.Select("*").
		From("users").
		Where(dbr.Eq("token", token)).
		LoadOne(&user)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return &user, nil
}

func (us *Storage) ChangePassword(user *entities.UserPasswordChange) error {
	tx, err := us.Session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.Update("users").
		Set("password_hash", user.NewPassword).
		Set("token", nil).
		Where(
			dbr.And(
				dbr.Eq("password_hash", user.Password),
				dbr.Eq("login", user.Login),
			),
		).
		Exec()
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (us *Storage) AddUser(user *entities.User) error {
	tx, err := us.Session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertInto("users").
		Pair("login", user.Login).
		Pair("password_hash", user.Password).
		Returning("id").
		Load(&user.Id)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (us *Storage) Authorize(user *entities.User) (*entities.AccessToken, error) {
	tx, err := us.Session.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var ids []int64
	_, err = tx.Select("id").
		From("users").
		Where(
			dbr.And(
				dbr.Eq("login", user.Login),
				dbr.Eq("password_hash", user.Password),
			),
		).Load(&ids)

	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	token := uuid.New().String()
	_, err = tx.Update("users").
		Set("token", token).
		Where(
			dbr.And(
				dbr.Eq("login", user.Login),
				dbr.Eq("password_hash", user.Password),
			),
		).Exec()
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return &entities.AccessToken{Token: token, Type: entities.BearerType}, nil
}

func (us *Storage) AddGood(good *entities.Good) error {
	tx, err := us.Session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	if good.BriefInfo.Id != 0 {
		_, err = tx.Update("goods").
			Set("name", good.BriefInfo.Name).
			Set("cost", good.BriefInfo.Cost).
			Set("description", good.Description).
			Set("count", good.Count).
			Exec()
		if err != nil {
			return err
		}
	} else {
		err = tx.InsertInto("goods").
			Pair("name", good.BriefInfo.Name).
			Pair("cost", good.BriefInfo.Cost).
			Pair("description", good.Description).
			Pair("count", good.Count).
			Returning("id").
			Load(&good.BriefInfo.Id)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (us *Storage) GetGoodsBrief() ([]entities.GoodBrief, error) {
	tx, err := us.Session.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var goods []entities.GoodBrief
	_, err = tx.Select("id", "name", "cost").
		From("goods").
		Load(&goods)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return goods, nil
}

func (us *Storage) GetGoodFull(id int64) (*entities.Good, error) {
	tx, err := us.Session.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var good entities.Good
	err = tx.Select("id", "name", "cost", "count", "description").
		From("goods").
		Where(dbr.Eq("id", id)).
		LoadOne(&good)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return &good, nil
}

func (us *Storage) AddOrder(order *entities.Order) error {
	tx, err := us.Session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertInto("orders").
		Pair("status", false).
		Pair("user_id", order.UserId).
		Returning("id").
		Load(&order.Id)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (us *Storage) GetOrderUser(orderID int64) (int64, error) {
	tx, err := us.Session.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.RollbackUnlessCommitted()

	var userID int64
	err = tx.Select("user_id").
		From("orders").
		Where(dbr.Eq("id", orderID)).
		LoadOne(&userID)
	if err != nil {
		return 0, err
	}

	tx.Commit()
	return userID, nil
}

func (us *Storage) AddGoodToOrder(order *entities.OrderItem) error {
	tx, err := us.Session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	var count int64
	err = tx.Select("count").
		From("goods").
		Where(dbr.Eq("id", order.GoodId)).
		LoadOne(&count)
	if err != nil {
		return err
	}

	if count < order.Count {
		return errors.New("Not enough goods")
	}

	err = tx.InsertInto("order_list").
		Pair("order_id", order.OrderId).
		Pair("good_id", order.GoodId).
		Pair("count", order.Count).
		Returning("id").
		Load(&order.Id)
	if err != nil {
		return err
	}

	_, err = tx.Update("goods").
		Set("count", count-order.Count).
		Where(dbr.Eq("id", order.GoodId)).
		Exec()
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (us *Storage) PayOrder(orderId int64) error {
	tx, err := us.Session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.Update("orders").
		Set("status", true).
		Where(dbr.Eq("id", orderId)).
		Exec()
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (us *Storage) GetOrder(orderId int64) error {
	tx, err := us.Session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	var order *entities.Order
	err = tx.Select("*").
		Where(dbr.Eq("order_id", orderId)).
		LoadOne(&order)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func NewStorage() *Storage {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		"localhost", "5432", "postgres", "InternetShop", "123", "disable")

	conn, err := dbr.Open("postgres", dsn, nil)
	if err != nil {
		return nil
	}
	conn.SetMaxOpenConns(10)
	return &Storage{
		Session: conn.NewSession(nil),
	}
}
