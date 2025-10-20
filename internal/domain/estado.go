package domain

import "time"

// Regiao representa uma região geográfica do Brasil
type Regiao struct {
	ID    int    `bson:"id" json:"id"`
	Sigla string `bson:"sigla" json:"sigla"`
	Nome  string `bson:"nome" json:"nome"`
}

// Estado representa um estado brasileiro (UF)
type Estado struct {
	ID        int       `bson:"id" json:"id"`
	Sigla     string    `bson:"sigla" json:"sigla"`
	Nome      string    `bson:"nome" json:"nome"`
	Regiao    Regiao    `bson:"regiao" json:"regiao"`
	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
