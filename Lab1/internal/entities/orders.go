package entities

type Order struct {
	Id     int64 `json:"id" db:"id"`
	UserId int64 `json:"user_id" db:"user_id"`
	Status bool  `json:"status" db:"status"`
}

type OrderItem struct {
	Id      int64 `json:"id" db:"id"`
	Count   int64 `json:"count" db:"count"`
	OrderId int64 `json:"order_id" db:"order_id"`
	GoodId  int64 `json:"good_id" db:"good_id"`
}

func (oi *OrderItem) Valid() bool {
	return oi.OrderId > 0 && oi.Count > 0 && oi.GoodId > 0
}

type OrderInfo struct {
}
