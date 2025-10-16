package domain

import "time"

type RefreshToken struct {
	ID        string    `bson:"_id"`         // jti
	UserID    string    `bson:"userId"`
	ExpiresAt time.Time `bson:"expiresAt"`   // TTL index
	Revoked   bool      `bson:"revoked"`
	CreatedAt time.Time `bson:"createdAt"`
	ParentJTI string    `bson:"parentJti,omitempty"` // para rotação segura
}

