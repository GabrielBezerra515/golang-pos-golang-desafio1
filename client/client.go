package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		logger.Fatalf("http.NewRequestWithContext: Error in Construct Request [%v]", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Fatalf("http.DefaultClient.Do: Error to execute Request [%s]", err)
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Fatalf("ioutil.ReadAll: Error to read Request Body [%s]", err)
	}

	f, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("os.OpenFile: Error to open File [%s]", err)
	}

	if _, err := f.WriteString(fmt.Sprintf("Dólar: {%s}\n", string(resBody))); err != nil {
		logger.Fatalf("f.WriteString: Error to append dolar info to File [%s]", err)
	}

	logger.Printf("Succes get and write dolar price.")
}
