package users

import "github.com/google/uuid"

type User struct {
	Id        uuid.UUID
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

func (u User) IsValid() bool {
	// TODO: Implement
	return false
}
