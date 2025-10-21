package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("🌱 Seed de Dados de Teste - Retech Core")
	fmt.Println("========================================")
	fmt.Println()

	// Conectar ao MongoDB
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "retech_core")

	fmt.Printf("Conectando ao MongoDB: %s/%s\n", mongoURI, mongoDB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Erro ao conectar MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(mongoDB)

	// 1. Criar Super Admin
	fmt.Println("\n1️⃣  Criando Super Admin...")
	if err := createSuperAdmin(db); err != nil {
		log.Printf("   ⚠️  Super admin: %v\n", err)
	} else {
		fmt.Println("   ✅ Super admin criado")
	}

	// 2. Criar Tenant de Teste
	fmt.Println("\n2️⃣  Criando Tenant de Teste...")
	tenantID, err := createTestTenant(db)
	if err != nil {
		log.Printf("   ⚠️  Tenant: %v\n", err)
	} else {
		fmt.Printf("   ✅ Tenant criado: %s\n", tenantID)
	}

	// 3. Criar Usuário Tenant
	fmt.Println("\n3️⃣  Criando Usuário Tenant...")
	if err := createTenantUser(db, tenantID); err != nil {
		log.Printf("   ⚠️  Tenant user: %v\n", err)
	} else {
		fmt.Println("   ✅ Tenant user criado")
	}

	// 4. Criar API Keys de Teste
	fmt.Println("\n4️⃣  Criando API Keys de Teste...")
	apiKey, err := createTestAPIKey(db, tenantID)
	if err != nil {
		log.Printf("   ⚠️  API Key: %v\n", err)
	} else {
		fmt.Printf("   ✅ API Key criada: %s\n", apiKey)
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ Seed completo!")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
	fmt.Println("📋 CREDENCIAIS CRIADAS:")
	fmt.Println()
	fmt.Println("🔐 SUPER ADMIN (você):")
	fmt.Println("   Email: admin@theretech.com.br")
	fmt.Println("   Senha: admin12345678")
	fmt.Println("   Login: http://localhost:3000/admin/login")
	fmt.Println()
	fmt.Println("👨‍💻 TENANT USER (teste):")
	fmt.Println("   Email: dev@teste.com")
	fmt.Println("   Senha: teste12345678")
	fmt.Println("   Login: http://localhost:3000/painel/login")
	fmt.Println()
	fmt.Println("🔑 API KEY (teste):")
	fmt.Printf("   Key: %s\n", apiKey)
	fmt.Println("   Uso: curl http://localhost:8080/geo/ufs -H 'x-api-key: " + apiKey + "'")
	fmt.Println()
}

func createSuperAdmin(db *mongo.Database) error {
	coll := db.Collection("users")
	
	// Verificar se já existe
	count, _ := coll.CountDocuments(context.Background(), bson.M{"email": "admin@theretech.com.br"})
	if count > 0 {
		return fmt.Errorf("já existe")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin12345678"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	admin := bson.M{
		"email":     "admin@theretech.com.br",
		"password":  string(hashedPassword),
		"name":      "Alan Leite (Super Admin)",
		"role":      "SUPER_ADMIN",
		"active":    true,
		"createdAt": now,
		"updatedAt": now,
	}

	_, err = coll.InsertOne(context.Background(), admin)
	return err
}

func createTestTenant(db *mongo.Database) (string, error) {
	coll := db.Collection("tenants")
	
	tenantID := "tenant-teste-" + time.Now().Format("20060102")
	
	// Verificar se já existe
	count, _ := coll.CountDocuments(context.Background(), bson.M{"tenantId": tenantID})
	if count > 0 {
		return tenantID, fmt.Errorf("já existe")
	}

	now := time.Now()
	tenant := bson.M{
		"tenantId":  tenantID,
		"name":      "Empresa Teste LTDA",
		"email":     "contato@teste.com",
		"active":    true,
		"createdAt": now,
		"updatedAt": now,
	}

	_, err := coll.InsertOne(context.Background(), tenant)
	return tenantID, err
}

func createTenantUser(db *mongo.Database, tenantID string) error {
	coll := db.Collection("users")
	
	// Verificar se já existe
	count, _ := coll.CountDocuments(context.Background(), bson.M{"email": "dev@teste.com"})
	if count > 0 {
		return fmt.Errorf("já existe")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("teste12345678"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	user := bson.M{
		"email":     "dev@teste.com",
		"password":  string(hashedPassword),
		"name":      "Desenvolvedor Teste",
		"role":      "TENANT_USER",
		"tenantId":  tenantID,
		"active":    true,
		"createdAt": now,
		"updatedAt": now,
	}

	_, err = coll.InsertOne(context.Background(), user)
	return err
}

func createTestAPIKey(db *mongo.Database, tenantID string) (string, error) {
	coll := db.Collection("api_keys")
	
	keyID := "rtc_test_" + time.Now().Format("20060102150405")
	
	now := time.Now()
	expiresAt := now.Add(365 * 24 * time.Hour) // 1 ano

	apiKey := bson.M{
		"keyId":     keyID,
		"keyHash":   "hash_" + keyID, // Simplificado para teste
		"scopes":    []string{"geo:read"},
		"ownerId":   tenantID,
		"expiresAt": expiresAt,
		"revoked":   false,
		"createdAt": now,
	}

	_, err := coll.InsertOne(context.Background(), apiKey)
	return keyID, err
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}


