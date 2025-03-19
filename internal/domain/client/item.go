package client

import (
	"time"

	"github.com/google/uuid"
)

// Category тип под категорию.
type Category int

const (
	// CategoryPassword пароль.
	CategoryPassword Category = iota
	// CategoryText текст.
	CategoryText
	// CategoryFile файл (бинарные данные).
	CategoryFile
	// CategoryBankCard банковская карта.
	CategoryBankCard
)

// String возвращает строковое представление категории.
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

// RawItem описывает элемент как он приходит с сервера.
type RawItem struct {
	GUID      uuid.UUID
	Category  Category
	Title     string
	Value     *[]byte
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Item содержит элемент после разбора.
// В Value будет лежать структура в зависимости от категории Category - Password, File, BankCard.
type Item struct {
	GUID      uuid.UUID
	Category  Category
	Title     string
	Value     interface{}
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewItem создает новый элемент.
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

// NewPasswordItem создает новый элемент категории CategoryPassword.
func NewPasswordItem(title, login, password, note string) Item {
	return NewItem(
		CategoryPassword,
		title,
		Password{
			Login:    login,
			Password: password,
		},
		note,
	)
}

// NewTextItem создает новый элемент категории CategoryText.
func NewTextItem(title, note string) Item {
	return NewItem(
		CategoryText,
		title,
		nil,
		note,
	)
}

// NewBankCardItem создает новый элемент категории CategoryBankCard.
func NewBankCardItem(title, number, cvv, expiry, note string) Item {
	return NewItem(
		CategoryBankCard,
		title,
		BankCard{
			Number: number,
			CVV:    cvv,
			Expiry: expiry,
		},
		note,
	)
}

// NewFileItem создает новый элемент категории CategoryFile.
func NewFileItem(title, name string, data []byte, note string) Item {
	return NewItem(
		CategoryFile,
		title,
		File{
			Name: name,
			Data: data,
		},
		note,
	)
}

// Password описывает пароль.
type Password struct {
	// Login логин.
	Login string `json:"login"`
	// Password пароль.
	Password string `json:"password"`
}

// File файл.
type File struct {
	// Name название файла.
	Name string `json:"path"`
	// Data контент файла.
	Data []byte `json:"data"`
}

// BankCard банковская карта.
type BankCard struct {
	// Number номер карты.
	Number string `json:"number"`
	// CVV код.
	CVV string `json:"cvv"`
	// Expiry действует до.
	Expiry string `json:"exp"`
}
