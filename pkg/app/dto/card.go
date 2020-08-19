package dto

import "github.com/DBoyara/Netology-Go-11/pkg/card"

type CardDTO struct {
	Id     int64       `json:"id"`
	UserId card.UserID `json:"userId"`
	Number int64       `json:"number"`
	Type   string      `json:"type"`
	Issuer string      `json:"issuer"`
}

type CardErrDTO struct {
	Err string `json:"error"`
}
