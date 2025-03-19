package sync

import "context"

func (m *Model) syncAction() error {
	item, err := m.usecaseSync.Do(context.TODO(), m.item.GUID)
	if err != nil {
		return err
	}

	m.item = *item

	return nil
}
