package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// Cliente HTTP customizado baseado no Stack Overflow
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true, // Força fechamento de conexões
		},
	}

	// Teste da rota de rotação
	url := "http://localhost:8080/apikeys/rotate"
	jsonData := `{"keyId": "e163fa7d-0b09-4f08-826f-2b575fb415d0"}`
	
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(jsonData))
	if err != nil {
		fmt.Printf("Erro ao criar request: %v\n", err)
		return
	}
	
	// Configurações baseadas no Stack Overflow
	req.Close = true // Força fechamento da conexão
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	
	fmt.Printf("Enviando request para: %s\n", url)
	fmt.Printf("Headers: %v\n", req.Header)
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro na requisição: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Headers da resposta: %v\n", resp.Header)
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler resposta: %v\n", err)
		return
	}
	
	fmt.Printf("Resposta: %s\n", string(body))
}
