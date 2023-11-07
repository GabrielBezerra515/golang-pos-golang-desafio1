package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type DolarReal struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		logger.Fatalf("http.NewRequestWithContext: Error in Construct Request [%v]", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Fatalf("http.DefaultClient.Do: Error to execute Request [%s]", err)
	}
	defer res.Body.Close()

	dr := DolarReal{}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Fatalf("ioutil.ReadAll: Error to read Request Body [%s]", err)
	}

	fmt.Printf("Status Code: %d\n", res.StatusCode)
	json.Unmarshal(resBody, &dr)

	if err := Database(ctx, &dr); err != nil {
		logger.Printf("Database: Error to insert Request to database [%s]", err)
	}

	drJson, err := json.Marshal(dr.USDBRL.Bid)
	if err != nil {
		logger.Fatalf("json.Marshal: Error to Marshal [%v]", err)
	}

	w.Write(drJson)
}
