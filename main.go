package main

import (
	"fmt"
	"io"
	"log"
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

	resultChViaCEP2 := make(chan string)
	defer close(resultChViaCEP2)
	resultChBrasilAPI2 := make(chan string)
	defer close(resultChBrasilAPI2)

	go MakeRequest2("viacep", "https://viacep.com.br/ws/"+cep+"/json/", resultChViaCEP2)
	go MakeRequest2("brasilapi", "https://brasilapi.com.br/api/cep/v1/"+cep, resultChBrasilAPI2)

	select {
	case resultViaCEP := <-resultChViaCEP2:
		fmt.Println(resultViaCEP)
	case resultBrasilAPI := <-resultChBrasilAPI2:
		fmt.Println(resultBrasilAPI)
	case <-time.After(1 * time.Second):
		fmt.Println("Erro de timeout")
	}
}

func MakeRequest2(apiName string, url string, resultChannel chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf(err.Error())
	}
	resultChannel <- fmt.Sprintf("Resultado da %s: %s", apiName, string(body))
}
