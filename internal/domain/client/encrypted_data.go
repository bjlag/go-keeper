package client

import (
	"encoding/json"
	"fmt"
)

// EncryptedData описывает зашифрованные данные, которые приходят с сервера.
type EncryptedData struct {
	// Title название элемента.
	Title string `json:"title"`
	// Category категория элемента: пароль, текст, файл и пр.
	Category Category `json:"category_id"`
	// Value у каждой категории свой набор данных, после расшифровки хранится в байтах, по сути это JSON.
	Value *[]byte `json:"value,omitempty"`
	// Notes какие-то заметки, которые можно указать у элемента.
	Notes string `json:"notes"`
}

func (d *EncryptedData) UnmarshalJSON(data []byte) error {
	type Alias EncryptedData

	alias := &struct {
		*Alias
		Value *json.RawMessage `json:"value,omitempty"`
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, alias); err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	if alias.Value != nil {
		value := []byte(*alias.Value)
		d.Value = &value
	}

	return nil
}

func (d EncryptedData) MarshalJSON() ([]byte, error) {
	type Alias EncryptedData

	alias := struct {
		Alias
		Value json.RawMessage `json:"value,omitempty"`
	}{
		Alias: Alias(d),
	}

	if d.Value != nil {
		alias.Value = *d.Value
	}

	return json.Marshal(alias)
}
