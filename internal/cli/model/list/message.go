package list

import (
	"github.com/bjlag/go-keeper/internal/domain/client"
)

type GetDataMsg struct{}

type OpenCategoriesMsg struct{}

type OpenItemsMsg struct {
	Category client.Category
}
