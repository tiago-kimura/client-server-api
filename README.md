# Go Webserver Challenge

In this challenge, we will apply what we have learned about HTTP web servers, contexts, databases, and file manipulation with Go.

### Prerequisites

- [Golang v1.22+](https://golang.org/) 

## Deliverables

You will need to deliver two Go systems:
- `client.go`
- `server.go`

## Requirements

The requirements to fulfill this challenge are:

1. **HTTP Request:**
   - `client.go` should make an HTTP request to `server.go` requesting the dollar exchange rate.

2. **Consume API:**
   - `server.go` should consume the API containing the USD to BRL exchange rate at the address: [Awesome API](https://economia.awesomeapi.com.br/json/last/USD-BRL).
   - `server.go` should return the result in JSON format to the client.

3. **Database Logging:**
   - Using the `context` package, `server.go` should log each received exchange rate in a SQLite database.
   - The maximum timeout for calling the exchange rate API should be 200ms.
   - The maximum timeout for persisting the data in the database should be 10ms.

4. **Client Handling:**
   - `client.go` needs to receive from `server.go` only the current exchange rate value (the `bid` field from the JSON).
   - Using the `context` package, `client.go` should have a maximum timeout of 300ms to receive the result from `server.go`.

5. **Error Logging:**
   - All 3 contexts should return an error in the logs if the execution time is insufficient.

6. **File Saving:**
   - `client.go` should save the current exchange rate in a file named `cotacao.txt` in the format: `Dollar: {value}`.

7. **Server Endpoint:**
   - The necessary endpoint generated by `server.go` for this challenge will be `/cotacao`.
   - The port to be used by the HTTP server will be `8080`.

### Running the server

Inside the "server" folder, type in the terminal:

`go run server.go`

### Running the client

Inside the "client" folder, type in the terminal:

`go run client.go`