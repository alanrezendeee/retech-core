package domain

import "time"

// ArtigoPenal representa um artigo do Código Penal ou legislação especial
type ArtigoPenal struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	Codigo        string    `json:"codigo" bson:"codigo"`         // "121", "121.1", "121.1.I", "121.1.I.a"
	Artigo        int       `json:"artigo" bson:"artigo"`         // 121
	Paragrafo     *int      `json:"paragrafo,omitempty" bson:"paragrafo,omitempty"`
	Inciso        *string   `json:"inciso,omitempty" bson:"inciso,omitempty"`   // "I", "II", "III"
	Alinea        *string   `json:"alinea,omitempty" bson:"alinea,omitempty"`   // "a", "b", "c"
	Descricao     string    `json:"descricao" bson:"descricao"`   // "Homicídio simples"
	TextoCompleto string    `json:"textoCompleto" bson:"textoCompleto"` // Texto completo do artigo
	Tipo          string    `json:"tipo" bson:"tipo"`              // "crime", "contravencao"
	Legislacao    string    `json:"legislacao" bson:"legislacao"` // "CP", "LCP", "Lei 11.343/2006"
	LegislacaoNome string   `json:"legislacaoNome" bson:"legislacaoNome"` // "Código Penal", "Lei de Contravenções Penais"
	PenaMin       string    `json:"penaMin,omitempty" bson:"penaMin,omitempty"` // "Reclusão, de 6 a 20 anos"
	PenaMax       string    `json:"penaMax,omitempty" bson:"penaMax,omitempty"` // "e multa"
	CodigoFormatado string  `json:"codigoFormatado" bson:"codigoFormatado"` // "Art. 121, § 1º, I, a) do CP"
	// Para busca e autocomplete
	Busca         string    `json:"-" bson:"busca"` // Texto normalizado para busca (lowercase, sem acentos)
	// Documentação e rastreabilidade
	Fonte         string    `json:"fonte" bson:"fonte"` // URL ou referência da fonte oficial
	DataAtualizacao string  `json:"dataAtualizacao" bson:"dataAtualizacao"` // Data da última atualização da fonte oficial
	HashConteudo  string    `json:"hashConteudo,omitempty" bson:"hashConteudo,omitempty"` // SHA256 para detectar alterações
	IdUnico       string    `json:"idUnico" bson:"idUnico"` // Identificador único: "LEGISLACAO:CODIGO" (ex: "CP:121", "Lei 11.343/2006:33")
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}

// PenalResponse representa a resposta da API para autocomplete
type PenalResponse struct {
	Codigo          string `json:"codigo"`
	CodigoFormatado string `json:"codigoFormatado"`
	Descricao       string `json:"descricao"`
	Tipo            string `json:"tipo"`
	Legislacao      string `json:"legislacao"`
	LegislacaoNome  string `json:"legislacaoNome"`
	IdUnico         string `json:"idUnico"` // Identificador único: "LEGISLACAO:CODIGO" para diferenciar artigos com mesmo código em legislações diferentes
}

