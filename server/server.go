package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ExchangeRate struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
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
}

func main() {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		println(err.Error())
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		println(err.Error())
	}
	var exchange ExchangeRate
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		println(err.Error())
	}
	println(exchange.USDBRL.Bid)
	fmt.Println(exchange)
}
