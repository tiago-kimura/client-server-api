package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExchangeRate struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	Id         string `json:"id"`
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

type Response struct {
	Bid string `json:"bid"`
}

func main() {
	NewDatabase()
	http.HandleFunc("/cotacao", handler)
	http.HandleFunc("/cotacoes", handlerGet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func NewDatabase() {
	db, err := sql.Open("sqlite3", "./currencies.db")
	if err != nil {
		log.Fatalf("Erro ao abrir banco de dados: %v\n", err)
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS exchange_rates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varBid TEXT,
		pctChange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	);`)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v\n", err)
	}
}

func handlerGet(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./currencies.db")
	if err != nil {
		log.Fatalf("Error to open data base: %v\n", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM exchange_rates ")
	if err != nil {
		log.Fatalf("Error to execute SQL: %v\n", err)
	}
	defer rows.Close()

	var rates []ExchangeRate
	for rows.Next() {
		var rate ExchangeRate
		err := rows.Scan(&rate.USDBRL.Id, &rate.USDBRL.Code, &rate.USDBRL.Codein, &rate.USDBRL.Name, &rate.USDBRL.High, &rate.USDBRL.Low,
			&rate.USDBRL.VarBid, &rate.USDBRL.PctChange, &rate.USDBRL.Bid, &rate.USDBRL.Ask, &rate.USDBRL.Timestamp, &rate.USDBRL.CreateDate)
		if err != nil {
			log.Fatalf("Error read result: %v\n", err)
		}
		rates = append(rates, rate)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Error to iterate result: %v\n", err)
	}
	body, _ := json.Marshal(rates)
	writeResponse(w, http.StatusOK, body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request started")
	ctx := r.Context()
	defer log.Println("Request finished")

	timeoutCtxEx, cancelEx := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancelEx()
	exchange, err := fetchExchangeRate(timeoutCtxEx, "USD-BRL")
	if err != nil {
		writeResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	timeoutCtxSv, cancelSv := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancelSv()
	err = saveExchangeRate(timeoutCtxSv, exchange)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	var response Response
	response.Bid = exchange.USDBRL.Bid
	body, _ := json.Marshal(response)
	writeResponse(w, http.StatusOK, body)

}

func fetchExchangeRate(ctx context.Context, typeConversion string) (ExchangeRate, error) {
	var exchange ExchangeRate
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/"+typeConversion, nil)
	if err != nil {
		log.Println(err.Error())
		return exchange, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return exchange, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return exchange, err
	}

	err = json.Unmarshal(body, &exchange)
	if err != nil {
		log.Println(err.Error())
		return exchange, err
	}
	return exchange, nil
}

func saveExchangeRate(ctx context.Context, exchange ExchangeRate) error {
	db, err := sql.Open("sqlite3", "./currencies.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query := "INSERT INTO exchange_rates (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err = db.ExecContext(ctx, query, exchange.USDBRL.Code, exchange.USDBRL.Codein, exchange.USDBRL.Name, exchange.USDBRL.High, exchange.USDBRL.Low, exchange.USDBRL.VarBid, exchange.USDBRL.PctChange, exchange.USDBRL.Bid, exchange.USDBRL.Ask, exchange.USDBRL.Timestamp, exchange.USDBRL.CreateDate)
	if err != nil {
		return err
	}
	return nil

}

func writeResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}
