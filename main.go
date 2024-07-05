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

	resultCh := make(chan string)
	defer close(resultCh)

	go MakeRequest2("viacep", "https://viacep.com.br/ws/"+cep+"/json/", resultCh)
	go MakeRequest2("brasilapi", "https://brasilapi.com.br/api/cep/v1/"+cep, resultCh)

	select {
	case result := <-resultCh:
		fmt.Println(result)
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
