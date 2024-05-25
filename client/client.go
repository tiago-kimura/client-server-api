package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Bid string `json:"bid"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	dollarBid, err := fetchExchangeRate(ctx)
	if err != nil {
		println("error: ", err.Error())
		return
	}
	err = WriteFile(dollarBid)
	if err != nil {
		println("error: ", err.Error())
		return
	}

	println("operation finished with success!")
}

func fetchExchangeRate(ctx context.Context) (Response, error) {
	var response Response
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		return response, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return response, err
	}
	if resp.StatusCode != http.StatusOK {
		return response, errors.New(resp.Status + " - " + string(body))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err.Error())
		return response, err
	}

	return response, nil

}

func WriteFile(dollarBid Response) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("Dólar:" + dollarBid.Bid)
	if err != nil {
		return err
	}
	return nil
}
