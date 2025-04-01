package user

import (
	"time"
)

// Models taken from https://docs.powens.com/api-reference/user-connections/users#data-model

type User struct {
	id     int
	signin time.Time
}

type UsersList struct {
	users []User
}
