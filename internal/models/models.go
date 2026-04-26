package models

import "github.com/jackc/pgx/v5"

type Link struct {
	OriginUrl string `json:"origin_url"`
	ShortUrl  string `json:"short_url"`
	Id        int    `json:"id"`
	Clicks    int    `json:"clicks"`
}

type Connect struct {
	conn *pgx.Conn
}
