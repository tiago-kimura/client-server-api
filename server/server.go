package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	//_ "github.com/mattn/go-sqlite3"
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

type Response struct {
	Bid string `json:"bid"`
}

//TODO: Ainda não está funcionando o timeout e precisa arrumar o banco. Além de fazer o lado do client

func writeResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}

var database *sql.DB

// func NewDatabase() {
// 	dbPath := "sqlite3://user:password@tcp(localhost:3306)/data/exchange_rates.db"
// 	var err error
// 	database, err = sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Panic("Error connecting to database:", err)
// 	}
// }

func main() {
	//	NewDatabase()
	http.HandleFunc("/dollarExchangeRate", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request iniciada")
	ctx := r.Context()
	defer log.Println("Request finalizada")

	select {
	case <-ctx.Done():
		log.Println("Request cancelada pelo cliente")
	case <-time.After(1 * time.Second):

		ctxEx, cancelEx := context.WithTimeout(ctx, 200*time.Millisecond)
		defer cancelEx()
		exchange, err := fetchExchangeRate(ctxEx, "USD-BRL")
		if err != nil {
			writeResponse(w, http.StatusBadRequest, []byte(err.Error()))
		}

		// ctxDB, cancelDB := context.WithTimeout(context.Background(), 20*time.Millisecond)
		// defer cancelDB()
		// err = saveExchangeRate(ctxDB, database, exchange)
		// if err != nil {
		// 	writeResponse(w, http.StatusBadRequest, []byte(err.Error()))
		// 	return
		// }

		var response Response
		response.Bid = exchange.USDBRL.Bid
		body, _ := json.Marshal(response)
		writeResponse(w, http.StatusOK, body)
	}
}

func fetchExchangeRate(_ context.Context, typeConvertion string) (ExchangeRate, error) {
	var exchange ExchangeRate
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/" + typeConvertion)
	if err != nil {
		log.Println(err.Error())
		return exchange, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err.Error())
		return exchange, err
	}

	err = json.Unmarshal(body, &exchange)
	if err != nil {
		log.Println(err.Error())
		return exchange, err
	}
	time.Sleep(3 * time.Second)
	return exchange, nil
}

// func saveExchangeRate(ctx context.Context, db *sql.DB, exchange ExchangeRate) error {
// 	query := "INSERT INTO exchange_rates (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
// 	stmt, err := db.PrepareContext(ctx, query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.ExecContext(ctx, exchange.Code, exchange.Codein, exchange.Name, exchange.High, exchange.Low, exchange.VarBid, exchange.PctChange, exchange.Bid, exchange.Ask, exchange.Timestamp, exchange.CreateDate)
// 	return err
// }
