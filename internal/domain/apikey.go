package domain

import "time"

type APIKey struct {
	ID        string    `bson:"_id,omitempty"`
	KeyID     string    `bson:"keyId"`       // público
	KeyHash   string    `bson:"keyHash"`     // HMAC(keyId.keySecret)
	Scopes    []string  `bson:"scopes"`
	Roles     []string  `bson:"roles"`
	OwnerID   string    `bson:"ownerId"`
	ExpiresAt time.Time `bson:"expiresAt"`   // TTL index
	Revoked   bool      `bson:"revoked"`
	CreatedAt time.Time `bson:"createdAt"`
}

