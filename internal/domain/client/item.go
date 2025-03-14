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

func NewItem(category Category, title string, value interface{}, note string) Item {
	now := time.Now()
	return Item{
		GUID:      uuid.New(),
		Category:  category,
		Title:     title,
		Value:     value,
		Notes:     note,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewPasswordItem(title, login, password, note string) Item {
	return NewItem(
		CategoryPassword,
		title,
		&Password{
			Login:    login,
			Password: password,
		},
		note,
	)
}

func NewTextItem(title, note string) Item {
	return NewItem(
		CategoryText,
		title,
		nil,
		note,
	)
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
