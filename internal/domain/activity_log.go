package domain

import (
	"time"
)

// TimeNow retorna o tempo atual em UTC
func TimeNow() time.Time {
	return time.Now().UTC()
}

// ActivityLog representa uma ação realizada no sistema (auditoria)
type ActivityLog struct {
	ID        string                 `bson:"_id,omitempty" json:"id"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Type      string                 `bson:"type" json:"type"` // tenant.created, apikey.revoked, etc
	Actor     Actor                  `bson:"actor" json:"actor"`
	Resource  Resource               `bson:"resource" json:"resource"`
	Action    string                 `bson:"action" json:"action"` // create, update, delete, login, etc
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IP        string                 `bson:"ip,omitempty" json:"ip,omitempty"`
	UserAgent string                 `bson:"userAgent,omitempty" json:"userAgent,omitempty"`
}

// Actor representa quem realizou a ação
type Actor struct {
	UserID string `bson:"userId" json:"userId"`
	Email  string `bson:"email" json:"email"`
	Name   string `bson:"name" json:"name"`
	Role   string `bson:"role" json:"role"` // SUPER_ADMIN, TENANT_USER
}

// Resource representa o recurso afetado
type Resource struct {
	Type string `bson:"type" json:"type"` // tenant, apikey, settings, user
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

// ActivityType constantes para tipos de atividades
const (
	// Tenant events
	ActivityTypeTenantCreated     = "tenant.created"
	ActivityTypeTenantUpdated     = "tenant.updated"
	ActivityTypeTenantDeleted     = "tenant.deleted"
	ActivityTypeTenantActivated   = "tenant.activated"
	ActivityTypeTenantDeactivated = "tenant.deactivated"

	// API Key events
	ActivityTypeAPIKeyCreated = "apikey.created"
	ActivityTypeAPIKeyRevoked = "apikey.revoked"
	ActivityTypeAPIKeyRotated = "apikey.rotated"

	// Settings events
	ActivityTypeSettingsUpdated = "settings.updated"

	// User events
	ActivityTypeUserCreated = "user.created"
	ActivityTypeUserLogin   = "user.login"
	ActivityTypeUserLogout  = "user.logout"

	// System events
	ActivityTypeSystemStartup  = "system.startup"
	ActivityTypeSystemShutdown = "system.shutdown"
)

// ResourceType constantes para tipos de recursos
const (
	ResourceTypeTenant   = "tenant"
	ResourceTypeAPIKey   = "apikey"
	ResourceTypeSettings = "settings"
	ResourceTypeUser     = "user"
	ResourceTypeSystem   = "system"
)

// Action constantes para ações
const (
	ActionCreate     = "create"
	ActionUpdate     = "update"
	ActionDelete     = "delete"
	ActionActivate   = "activate"
	ActionDeactivate = "deactivate"
	ActionRevoke     = "revoke"
	ActionRotate     = "rotate"
	ActionLogin      = "login"
	ActionLogout     = "logout"
	ActionStartup    = "startup"
	ActionShutdown   = "shutdown"
)

