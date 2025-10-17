package domain

import "time"

type Tenant struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	TenantID  string    `bson:"tenantId" json:"tenantId"`         // ID único do tenant (ex: "tenant-1")
	Name      string    `bson:"name" json:"name"`                 // Nome do cliente
	Email     string    `bson:"email" json:"email"`               // Email de contato
	Active    bool      `bson:"active" json:"active"`             // Tenant ativo?
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

