package dto

type CardDTO struct {
	Id     int64  `json:"id"`
	UserId int64  `json:"userId"`
	Number int64  `json:"number"`
	Type   string `json:"type"`
	Issuer string `json:"issuer"`
}

type CardErrDTO struct {
	Err string `json:"error"`
}

