package client

import (
	"time"

	"github.com/google/uuid"
)

type Category int

const (
	CategoryPassword Category = iota
	CategoryText
	CategoryFile
	CategoryBankCard
)

func (c Category) String() string {
	switch c {
	case CategoryPassword:
		return "Пароль"
	case CategoryText:
		return "Текст"
	case CategoryFile:
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

func NewBankCardItem(title, number, cvv, expiry, note string) Item {
	return NewItem(
		CategoryBankCard,
		title,
		&BankCard{
			Number: number,
			CVV:    cvv,
			Expiry: expiry,
		},
		note,
	)
}

func NewFileItem(title, name string, data []byte, note string) Item {
	return NewItem(
		CategoryFile,
		title,
		&File{
			Name: name,
			Data: data,
		},
		note,
	)
}

type Password struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type File struct {
	Name string `json:"path"`
	Data []byte `json:"data"`
}

type BankCard struct {
	Number string `json:"number"`
	CVV    string `json:"cvv"`
	Expiry string `json:"exp"`
}
