package password

import (
	"context"
	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/domain/client"
)

func (m *Model) saveAction() error {
	item := client.Item{
		GUID:     uuid.New(),
		Category: client.CategoryPassword,
		Title:    element.GetValue(m.elements, posCreateTitle),
		Value: client.Password{
			Login:    element.GetValue(m.elements, posCreateLogin),
			Password: element.GetValue(m.elements, posCreatePassword),
		},
		Notes: element.GetValue(m.elements, posCreateNotes),
	}

	return m.usecaseCreate.Do(context.TODO(), item)
}

func (m *Model) editAction() error {
	i := client.Item{
		GUID:     m.guid,
		Category: m.category,
		Title:    element.GetValue(m.elements, posEditTitle),
		Value: client.Password{
			Login:    element.GetValue(m.elements, posEditLogin),
			Password: element.GetValue(m.elements, posEditPassword),
		},
		Notes: element.GetValue(m.elements, posEditNotes),
	}

	return m.usecaseEdit.Do(context.TODO(), i)
}

func (m *Model) deleteAction() error {
	return m.usecaseDelete.Do(context.TODO(), m.guid)
}
