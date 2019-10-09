package model

type User struct {
	ID uint // the request ID
	UserID   string
	Username string
	Exp      int
	Tier     int
	Role     int
	Scopes   string
}
