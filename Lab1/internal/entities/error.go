package entities

type Error struct {
	Message string `json:"error"`
}

type Id struct {
	Id int64 `json:"id"`
}
