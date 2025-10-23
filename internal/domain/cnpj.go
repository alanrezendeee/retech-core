package domain

import "time"

// CNPJ representa os dados de uma empresa
type CNPJ struct {
	CNPJ                string               `json:"cnpj" bson:"cnpj"`
	RazaoSocial         string               `json:"razaoSocial" bson:"razaoSocial"`
	NomeFantasia        string               `json:"nomeFantasia,omitempty" bson:"nomeFantasia,omitempty"`
	Situacao            string               `json:"situacao" bson:"situacao"`
	DataSituacao        string               `json:"dataSituacao,omitempty" bson:"dataSituacao,omitempty"`
	DataAbertura        string               `json:"dataAbertura,omitempty" bson:"dataAbertura,omitempty"`
	Porte               string               `json:"porte,omitempty" bson:"porte,omitempty"`
	NaturezaJuridica    string               `json:"naturezaJuridica,omitempty" bson:"naturezaJuridica,omitempty"`
	CapitalSocial       float64              `json:"capitalSocial,omitempty" bson:"capitalSocial,omitempty"`
	Endereco            CNPJEndereco         `json:"endereco" bson:"endereco"`
	Telefones           []string             `json:"telefones,omitempty" bson:"telefones,omitempty"`
	Email               string               `json:"email,omitempty" bson:"email,omitempty"`
	AtividadePrincipal  CNPJAtividade        `json:"atividadePrincipal,omitempty" bson:"atividadePrincipal,omitempty"`
	AtividadesSecundarias []CNPJAtividade    `json:"atividadesSecundarias,omitempty" bson:"atividadesSecundarias,omitempty"`
	QSA                 []CNPJSocio          `json:"qsa,omitempty" bson:"qsa,omitempty"`
	Source              string               `json:"source" bson:"source"` // brasilapi, receitaws, cache
	CachedAt            time.Time            `json:"cachedAt,omitempty" bson:"cachedAt,omitempty"`
}

// CNPJEndereco representa o endereço da empresa
type CNPJEndereco struct {
	Logradouro  string `json:"logradouro,omitempty" bson:"logradouro,omitempty"`
	Numero      string `json:"numero,omitempty" bson:"numero,omitempty"`
	Complemento string `json:"complemento,omitempty" bson:"complemento,omitempty"`
	Bairro      string `json:"bairro,omitempty" bson:"bairro,omitempty"`
	CEP         string `json:"cep,omitempty" bson:"cep,omitempty"`
	Municipio   string `json:"municipio,omitempty" bson:"municipio,omitempty"`
	UF          string `json:"uf,omitempty" bson:"uf,omitempty"`
}

// CNPJAtividade representa uma atividade econômica (CNAE)
type CNPJAtividade struct {
	Codigo    string `json:"codigo" bson:"codigo"`
	Descricao string `json:"descricao" bson:"descricao"`
}

// CNPJSocio representa um sócio ou administrador
type CNPJSocio struct {
	Nome         string `json:"nome" bson:"nome"`
	Qualificacao string `json:"qualificacao,omitempty" bson:"qualificacao,omitempty"`
}

// ValidateCNPJ valida o formato e dígito verificador de um CNPJ
func ValidateCNPJ(cnpj string) bool {
	// Remove caracteres não numéricos
	cleaned := ""
	for _, char := range cnpj {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}

	// CNPJ deve ter 14 dígitos
	if len(cleaned) != 14 {
		return false
	}

	// CNPJ não pode ser sequência de números iguais
	allSame := true
	for i := 1; i < len(cleaned); i++ {
		if cleaned[i] != cleaned[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	// Validar primeiro dígito verificador
	sum := 0
	multipliers1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	for i := 0; i < 12; i++ {
		sum += int(cleaned[i]-'0') * multipliers1[i]
	}
	remainder := sum % 11
	digit1 := 0
	if remainder >= 2 {
		digit1 = 11 - remainder
	}
	if int(cleaned[12]-'0') != digit1 {
		return false
	}

	// Validar segundo dígito verificador
	sum = 0
	multipliers2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	for i := 0; i < 13; i++ {
		sum += int(cleaned[i]-'0') * multipliers2[i]
	}
	remainder = sum % 11
	digit2 := 0
	if remainder >= 2 {
		digit2 = 11 - remainder
	}
	if int(cleaned[13]-'0') != digit2 {
		return false
	}

	return true
}

// NormalizeCNPJ remove formatação de um CNPJ
func NormalizeCNPJ(cnpj string) string {
	cleaned := ""
	for _, char := range cnpj {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}
	return cleaned
}

