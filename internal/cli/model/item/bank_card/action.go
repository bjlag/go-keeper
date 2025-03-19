package bank_card

import (
	"context"

	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/domain/client"
)

func (m *Model) createAction() error {
	item := client.NewBankCardItem(
		element.GetValue(m.elements, posCreateTitle),
		element.GetValue(m.elements, posCreateNumber),
		element.GetValue(m.elements, posCreateCVV),
		element.GetValue(m.elements, posCreateExpiry),
		element.GetValue(m.elements, posCreateNotes),
	)

	return m.usecaseCreate.Do(context.TODO(), item)
}

func (m *Model) editAction() error {
	item := client.NewBankCardItem(
		element.GetValue(m.elements, posEditTitle),
		element.GetValue(m.elements, posEditNumber),
		element.GetValue(m.elements, posEditCVV),
		element.GetValue(m.elements, posEditExpiry),
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
