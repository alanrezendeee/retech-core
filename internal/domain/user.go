package domain

import "time"

// UserRole define os tipos de usuários do sistema
type UserRole string

const (
	RoleSuperAdmin UserRole = "SUPER_ADMIN" // Acesso total ao sistema
	RoleTenantUser UserRole = "TENANT_USER" // Acesso apenas aos próprios dados
)

// User representa um usuário do sistema (admin ou tenant user)
type User struct {
	ID        string     `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string     `bson:"email" json:"email"`                           // Email único
	Password  string     `bson:"password" json:"-"`                            // Hash bcrypt (nunca retornar em JSON)
	Name      string     `bson:"name" json:"name"`                             // Nome completo
	Role      UserRole   `bson:"role" json:"role"`                             // SUPER_ADMIN ou TENANT_USER
	TenantID  string     `bson:"tenantId,omitempty" json:"tenantId,omitempty"` // ID do tenant (null para SUPER_ADMIN)
	Active    bool       `bson:"active" json:"active"`                         // Usuário ativo?
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updatedAt" json:"updatedAt"`
	LastLogin *time.Time `bson:"lastLogin,omitempty" json:"lastLogin,omitempty"` // Último login
}

// IsSuperAdmin verifica se o usuário é super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// IsTenantUser verifica se o usuário é um usuário de tenant
func (u *User) IsTenantUser() bool {
	return u.Role == RoleTenantUser
}

// CanAccessTenant verifica se o usuário pode acessar dados do tenant
func (u *User) CanAccessTenant(tenantID string) bool {
	// Super admin pode acessar qualquer tenant
	if u.IsSuperAdmin() {
		return true
	}
	// Tenant user só pode acessar seu próprio tenant
	return u.TenantID == tenantID
}

// CreateUserRequest representa o payload para criar usuário
type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Name     string   `json:"name" binding:"required"`
	Role     UserRole `json:"role" binding:"required,oneof=SUPER_ADMIN TENANT_USER"`
	TenantID string   `json:"tenantId,omitempty"` // Obrigatório se role=TENANT_USER
}

// LoginRequest representa o payload de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // Segundos até expirar
	User         *User  `json:"user"`
}

// RefreshTokenRequest representa o payload para refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RegisterTenantRequest representa o payload de registro self-service
type RegisterTenantRequest struct {
	// Dados do tenant
	TenantName  string `json:"tenantName" binding:"required"`
	TenantEmail string `json:"tenantEmail" binding:"required,email"`
	Company     string `json:"company"`
	Purpose     string `json:"purpose"` // Para que vai usar a API

	// Dados do primeiro usuário
	UserName     string `json:"userName" binding:"required"`
	UserEmail    string `json:"userEmail" binding:"required,email"`
	UserPassword string `json:"userPassword" binding:"required,min=8"`
}

// RegisterTenantResponse representa a resposta do registro
type RegisterTenantResponse struct {
	Tenant       *Tenant `json:"tenant"`
	User         *User   `json:"user"`
	AccessToken  string  `json:"accessToken"`
	RefreshToken string  `json:"refreshToken"`
}
