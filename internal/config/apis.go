package config

import (
	"fmt"
	"os"
	"time"
)

// ValidateExternalAPIsConfig valida se todas as ENVs obrigatórias estão configuradas
// DEVE ser chamada no startup (main.go) para falhar rápido
func ValidateExternalAPIsConfig() {
	missing := []string{}
	
	if os.Getenv("CEP_PRIMARY_URL") == "" {
		missing = append(missing, "CEP_PRIMARY_URL")
	}
	if os.Getenv("CEP_FALLBACK_URL") == "" {
		missing = append(missing, "CEP_FALLBACK_URL")
	}
	if os.Getenv("CNPJ_PRIMARY_URL") == "" {
		missing = append(missing, "CNPJ_PRIMARY_URL")
	}
	if os.Getenv("CNPJ_FALLBACK_URL") == "" {
		missing = append(missing, "CNPJ_FALLBACK_URL")
	}
	
	if len(missing) > 0 {
		fmt.Printf("\n🔴 ERRO DE CONFIGURAÇÃO: Variáveis obrigatórias não configuradas:\n")
		for _, env := range missing {
			fmt.Printf("   - %s\n", env)
		}
		fmt.Printf("\nConfigure no arquivo .env e reinicie a aplicação.\n")
		fmt.Printf("Exemplo: ver env.example\n\n")
		panic("Configuração de APIs externas incompleta!")
	}
	
	fmt.Printf("✅ [CONFIG] APIs externas configuradas corretamente\n")
}

// GetCEPPrimaryURL retorna URL do provider primário de CEP (OBRIGATÓRIO)
func GetCEPPrimaryURL() string {
	url := getenv("CEP_PRIMARY_URL", "")
	if url == "" {
		panic("🔴 CEP_PRIMARY_URL não configurada! Configure no .env antes de iniciar.")
	}
	return url
}

// GetCEPFallbackURL retorna URL do provider fallback de CEP (OBRIGATÓRIO)
func GetCEPFallbackURL() string {
	url := getenv("CEP_FALLBACK_URL", "")
	if url == "" {
		panic("🔴 CEP_FALLBACK_URL não configurada! Configure no .env antes de iniciar.")
	}
	return url
}

// GetCEPTimeout retorna timeout para requisições de CEP
func GetCEPTimeout() time.Duration {
	return parseDuration(getenv("CEP_TIMEOUT", "5s"))
}

// GetCNPJPrimaryURL retorna URL do provider primário de CNPJ (OBRIGATÓRIO)
func GetCNPJPrimaryURL() string {
	url := getenv("CNPJ_PRIMARY_URL", "")
	if url == "" {
		panic("🔴 CNPJ_PRIMARY_URL não configurada! Configure no .env antes de iniciar.")
	}
	return url
}

// GetCNPJFallbackURL retorna URL do provider fallback de CNPJ (OBRIGATÓRIO)
func GetCNPJFallbackURL() string {
	url := getenv("CNPJ_FALLBACK_URL", "")
	if url == "" {
		panic("🔴 CNPJ_FALLBACK_URL não configurada! Configure no .env antes de iniciar.")
	}
	return url
}

// GetCNPJTimeout retorna timeout para requisições de CNPJ
func GetCNPJTimeout() time.Duration {
	return parseDuration(getenv("CNPJ_TIMEOUT", "10s"))
}

// parseDuration converte string para time.Duration
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 5 * time.Second
	}
	return d
}

