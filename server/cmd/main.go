package main

import (
	"context"
	"fmt"
	"net/http"
	"url-shortener/internal/handlers"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5435/postgres?sslmode=disable")
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	defer conn.Close(context.Background())

	conn.Exec(context.Background(), "CREATE DATABASE urlshortener")
	h := &handlers.Handlers{Conn: conn}

	http.HandleFunc("/shorten", h.CreateUrlShort)
	http.HandleFunc("/stats/", h.GetStat)
	http.HandleFunc("/delete/", h.UrlDelete)
	http.HandleFunc("/", h.RedirectUrl)

	if err := http.ListenAndServe(":8060", nil); err != nil {
		fmt.Println("Error open serv")
	}
}
