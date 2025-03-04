package list

import "github.com/bjlag/go-keeper/internal/cli/element"

type GetAllDataMessage struct{}

type OpenCategoryListMessage struct{}

type OpenPasswordListMessage struct {
	Category element.Item
}
