package text

import (
	"context"

	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/domain/client"
)

func (m *Model) createAction() error {
	item := client.NewTextItem(
		element.GetValue(m.elements, posCreateTitle),
		element.GetValue(m.elements, posCreateNotes),
	)

	return m.usecaseCreate.Do(context.TODO(), item)
}

func (m *Model) editAction() error {
	item := client.NewTextItem(
		element.GetValue(m.elements, posEditTitle),
		element.GetValue(m.elements, posEditNotes),
	)

	item.GUID = m.guid

	return m.usecaseEdit.Do(context.TODO(), item)
}

func (m *Model) deleteAction() error {
	return m.usecaseDelete.Do(context.TODO(), m.guid)
}
