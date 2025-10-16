package domain

import "time"

type User struct {
	ID        string   `bson:"_id,omitempty" json:"id"`
	Email     string   `bson:"email" json:"email"`
	Password  string   `bson:"password" json:"-"`
	Roles     []string `bson:"roles" json:"roles"`
	Active    bool     `bson:"active" json:"active"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

