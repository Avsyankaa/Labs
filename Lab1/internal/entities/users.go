package entities

const (
	BearerType = "Bearer"
)

type User struct {
	Id       int64  `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password_hash"`
	Admin    bool   `json:"admin" db:"admin"`
}

type UserPasswordChange struct {
	Login       string `json:"login" db:"login"`
	Password    string `json:"password" db:"password_hash"`
	NewPassword string `json:"new_password"`
}

type AccessToken struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

func (u *UserPasswordChange) Valid() bool {
	return u.NewPassword != "" && u.Login != "" && u.Password != ""
}

func (u *User) Valid() bool {
	return u.Password != "" && u.Login != ""
}
