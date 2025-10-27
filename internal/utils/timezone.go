package utils

import "time"

// GetBrasiliaTime retorna o horário atual no timezone de Brasília (America/Sao_Paulo)
func GetBrasiliaTime() time.Time {
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		// Fallback para UTC-3 se LoadLocation falhar
		return time.Now().Add(-3 * time.Hour)
	}
	return time.Now().In(location)
}

// GetTodayBrasilia retorna a data de hoje no formato YYYY-MM-DD no timezone de Brasília
func GetTodayBrasilia() string {
	return GetBrasiliaTime().Format("2006-01-02")
}

// GetStartOfMonthBrasilia retorna YYYY-MM do mês atual no timezone de Brasília
func GetStartOfMonthBrasilia() string {
	return GetBrasiliaTime().Format("2006-01")
}

