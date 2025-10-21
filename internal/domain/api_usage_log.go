package domain

import "time"

// APIUsageLog representa um log de uso da API
type APIUsageLog struct {
	ID           string    `bson:"_id,omitempty" json:"id,omitempty"`
	APIKey       string    `bson:"apiKey" json:"apiKey"`             // API Key usada
	TenantID     string    `bson:"tenantId" json:"tenantId"`         // Tenant que fez a request
	Endpoint     string    `bson:"endpoint" json:"endpoint"`         // Endpoint acessado
	Method       string    `bson:"method" json:"method"`             // HTTP method
	StatusCode   int       `bson:"statusCode" json:"statusCode"`     // Status code da resposta
	ResponseTime int64     `bson:"responseTime" json:"responseTime"` // Tempo de resposta em ms
	IPAddress    string    `bson:"ipAddress" json:"ipAddress"`       // IP do cliente
	UserAgent    string    `bson:"userAgent" json:"userAgent"`       // User agent
	Timestamp    time.Time `bson:"timestamp" json:"timestamp"`       // Quando aconteceu
	Date         string    `bson:"date" json:"date"`                 // YYYY-MM-DD para agregações
	Hour         int       `bson:"hour" json:"hour"`                 // Hora (0-23) para agregações
}

// UsageStats estatísticas de uso agregadas
type UsageStats struct {
	TotalRequests  int64              `json:"totalRequests"`
	TotalToday     int64              `json:"totalToday"`
	TotalThisMonth int64              `json:"totalThisMonth"`
	AvgResponseTime float64           `json:"avgResponseTime"` // ms
	TopEndpoints   []EndpointStat     `json:"topEndpoints"`
	ErrorRate      float64            `json:"errorRate"` // Percentual
	RequestsByDay  []RequestsByDay    `json:"requestsByDay"`
}

// EndpointStat estatística por endpoint
type EndpointStat struct {
	Endpoint string `json:"endpoint"`
	Count    int64  `json:"count"`
}

// RequestsByDay requests agrupados por dia
type RequestsByDay struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

