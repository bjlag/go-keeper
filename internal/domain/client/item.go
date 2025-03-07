package client

import (
	"time"

	"github.com/google/uuid"
)

type Category int

const (
	CategoryPassword Category = iota
	CategoryText
	CategoryBlob
	CategoryBankCard
)

func (c Category) String() string {
	switch c {
	case CategoryPassword:
		return "Пароль"
	case CategoryText:
		return "Текст"
	case CategoryBlob:
		return "Файл"
	case CategoryBankCard:
		return "Банковская карта"
	}
	return ""
}

type RawItem struct {
	GUID      uuid.UUID
	Category  Category
	Title     string
	Value     *[]byte
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Item struct {
	GUID      uuid.UUID
	Category  Category
	Title     string
	Value     interface{}
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Password struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Blob struct {
	Data string `json:"data"`
}

type BankCard struct {
	Number string    `json:"number"`
	CVV    string    `json:"cvv"`
	Expiry time.Time `json:"exp"`
}
