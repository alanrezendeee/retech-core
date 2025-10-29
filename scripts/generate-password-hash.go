package main

import (
	"fmt"
	"os"
	
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run generate-password-hash.go <senha>")
		os.Exit(1)
	}
	
	password := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Erro: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Senha: %s\n", password)
	fmt.Printf("Hash:  %s\n", string(hash))
}

