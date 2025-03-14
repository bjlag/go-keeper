package file

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/domain/client"
)

var (
	errNoFileSelected = errors.New("no file selected")
	errFileNotExist   = errors.New("file does not exist")
)

func (m *Model) createAction() error {
	if m.selectedFile == "" {
		return errNoFileSelected
	}

	data, err := os.ReadFile(m.selectedFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errFileNotExist
		}
		return err
	}

	item := client.NewFileItem(
		element.GetValue(m.elements, posCreateTitle),
		filepath.Base(m.selectedFile),
		data,
		element.GetValue(m.elements, posCreateNotes),
	)

	return m.usecaseCreate.Do(context.TODO(), item)
}

func (m *Model) editAction() error {
	item := client.NewFileItem(
		element.GetValue(m.elements, posEditTitle),
		m.selectedFile,
		m.fileData,
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
