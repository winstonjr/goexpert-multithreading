package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <cep>")
		os.Exit(1)
	}
	cep := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a channel to receive the results
	resultChViaCEP := make(chan string)
	defer close(resultChViaCEP)
	resultChBrasilAPI := make(chan string)
	defer close(resultChBrasilAPI)

	go func() {
		result, err := MakeRequest(ctx, "https://viacep.com.br/ws/"+cep+"/json/")
		if err != nil {
			fmt.Println("Erro ao buscar na viacep:", err)
			return
		}
		resultChViaCEP <- result
	}()

	go func() {
		result, err := MakeRequest(ctx, "https://brasilapi.com.br/api/cep/v1/"+cep)
		if err != nil {
			fmt.Println("Erro ao buscar na brasilapi:", err)
			return
		}
		resultChBrasilAPI <- result
	}()

	var result string
	select {
	case resultViaCEP := <-resultChViaCEP:
		result = resultViaCEP
	case resultBrasilAPI := <-resultChBrasilAPI:
		result = resultBrasilAPI
	case <-ctx.Done():
		result = "Erro ao buscar cep"
	}

	fmt.Println(result)

	cancel()
}

func MakeRequest(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
