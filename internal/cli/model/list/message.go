package list

import (
	"github.com/bjlag/go-keeper/internal/domain/client"
)

type GetAllDataMessage struct{}

type OpenCategoryListMessage struct{}

type OpenItemListMessage struct {
	Category client.Category
}
