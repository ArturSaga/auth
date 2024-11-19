package model

// User - модель для редиса
type User struct {
	ID              int64  `redis:"Id"`
	Name            string `redis:"Name"`
	Email           string `redis:"Email"`
	Password        string `redis:"Password"`
	PasswordConfirm string `redis:""`
	Role            string `redis:"Role"`
	CreatedAt       int64  `redis:"CreatedAt"`
	UpdatedAt       int64  `redis:"UpdatedAt"`
}
