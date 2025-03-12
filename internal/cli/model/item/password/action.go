package password

import (
	"context"

	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/domain/client"
)

func (m *Model) createAction() error {
	item := client.NewPasswordItem(
		element.GetValue(m.elements, posCreateTitle),
		element.GetValue(m.elements, posCreateLogin),
		element.GetValue(m.elements, posCreatePassword),
		element.GetValue(m.elements, posCreateNotes),
	)

	return m.usecaseCreate.Do(context.TODO(), item)
}

func (m *Model) editAction() error {
	item := client.NewPasswordItem(
		element.GetValue(m.elements, posEditTitle),
		element.GetValue(m.elements, posEditLogin),
		element.GetValue(m.elements, posEditPassword),
		element.GetValue(m.elements, posEditNotes),
	)
	item.GUID = m.guid
	item.CreatedAt = m.item.CreatedAt
	item.UpdatedAt = m.item.UpdatedAt

	return m.usecaseEdit.Do(context.TODO(), item)
}

func (m *Model) deleteAction() error {
	return m.usecaseDelete.Do(context.TODO(), m.guid)
}
