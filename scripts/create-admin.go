package main

import (
	"bufio"
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
	fmt.Println("🔐 Criar Super Admin - Retech Core")
	fmt.Println("===================================")
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

	// Verificar se já existe super admin
	coll := db.Collection("users")
	count, _ := coll.CountDocuments(ctx, bson.M{"role": "SUPER_ADMIN"})

	if count > 0 {
		fmt.Printf("\n⚠️  Já existe %d super admin(s) no banco.\n", count)
		fmt.Println("Deseja criar outro? (s/N):")
		var resp string
		fmt.Scanln(&resp)
		if resp != "s" && resp != "S" {
			fmt.Println("Cancelado.")
			return
		}
	}

	// Coletar dados usando bufio para ler linhas completas
	reader := bufio.NewReader(os.Stdin)
	var email, name, password string

	fmt.Print("\n📧 Email do admin: ")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)
	if email == "" {
		email = "admin@theretech.com.br"
		fmt.Printf("   Usando padrão: %s\n", email)
	}

	fmt.Print("👤 Nome completo: ")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Super Admin"
		fmt.Printf("   Usando padrão: %s\n", name)
	}

	fmt.Print("🔑 Senha (min 8 caracteres): ")
	password, _ = reader.ReadString('\n')
	password = strings.TrimSpace(password)
	if password == "" {
		password = "admin123456"
		fmt.Printf("   Usando padrão: %s\n", password)
	}

	if len(password) < 8 {
		log.Fatal("❌ Senha deve ter no mínimo 8 caracteres")
	}

	// Hash da senha
	fmt.Println("\n🔒 Gerando hash bcrypt...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Erro ao gerar hash:", err)
	}

	// Criar documento
	now := time.Now()
	admin := bson.M{
		"email":     email,
		"password":  string(hashedPassword),
		"name":      name,
		"role":      "SUPER_ADMIN",
		"active":    true,
		"createdAt": now,
		"updatedAt": now,
	}

	// Inserir no banco
	fmt.Println("💾 Inserindo no banco...")
	result, err := coll.InsertOne(context.Background(), admin)
	if err != nil {
		log.Fatal("❌ Erro ao inserir:", err)
	}

	fmt.Println()
	fmt.Println("✅ Super Admin criado com sucesso!")
	fmt.Println()
	fmt.Println("📋 Detalhes:")
	fmt.Printf("   ID: %v\n", result.InsertedID)
	fmt.Printf("   Email: %s\n", email)
	fmt.Printf("   Nome: %s\n", name)
	fmt.Printf("   Role: SUPER_ADMIN\n")
	fmt.Println()
	fmt.Println("🚀 Agora você pode fazer login em:")
	fmt.Println("   http://localhost:3000/admin/login")
	fmt.Println()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
