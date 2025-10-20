package domain

import "time"

// Microrregiao representa uma microrregião
type Microrregiao struct {
	ID           int          `bson:"id" json:"id"`
	Nome         string       `bson:"nome" json:"nome"`
	Mesorregiao  Mesorregiao  `bson:"mesorregiao" json:"mesorregiao"`
}

// Mesorregiao representa uma mesorregião
type Mesorregiao struct {
	ID   int         `bson:"id" json:"id"`
	Nome string      `bson:"nome" json:"nome"`
	UF   UFReference `bson:"UF" json:"UF"`
}

// UFReference é uma referência simplificada ao estado
type UFReference struct {
	ID     int    `bson:"id" json:"id"`
	Sigla  string `bson:"sigla" json:"sigla"`
	Nome   string `bson:"nome" json:"nome"`
	Regiao Regiao `bson:"regiao" json:"regiao"`
}

// RegiaoImediata representa uma região imediata (nova divisão do IBGE)
type RegiaoImediata struct {
	ID                    int                   `bson:"id" json:"id"`
	Nome                  string                `bson:"nome" json:"nome"`
	RegiaoIntermediaria   RegiaoIntermediaria   `bson:"regiao-intermediaria" json:"regiao-intermediaria"`
}

// RegiaoIntermediaria representa uma região intermediária
type RegiaoIntermediaria struct {
	ID   int         `bson:"id" json:"id"`
	Nome string      `bson:"nome" json:"nome"`
	UF   UFReference `bson:"UF" json:"UF"`
}

// Municipio representa um município brasileiro
type Municipio struct {
	ID              int             `bson:"id" json:"id"`
	Nome            string          `bson:"nome" json:"nome"`
	Microrregiao    Microrregiao    `bson:"microrregiao" json:"microrregiao"`
	RegiaoImediata  RegiaoImediata  `bson:"regiao-imediata" json:"regiao-imediata"`
	CreatedAt       time.Time       `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt       time.Time       `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

