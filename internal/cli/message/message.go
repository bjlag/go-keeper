package message

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/domain/client"
)

type (
	// OpenLoginMsg сообщение указывает, что надо открыть модель для аутентификации.
	OpenLoginMsg struct{}

	// OpenRegisterMsg сообщение указывает, что надо открыть модель для регистрации.
	OpenRegisterMsg struct {
		// LoginModel содержит ссылку на модель аутентификации, чтобы на нее можно было вернуться
		// в случае отказа от регистрации.
		LoginModel tea.Model
	}

	// SuccessLoginMsg вспомогательное сообщение, что какое-то событие выполнилось успешно, например, аутентификация.
	SuccessLoginMsg struct {
		Email    string
		Password string
	}

	// OpenCategoriesMsg сообщение используется для открытия списка категорий.
	OpenCategoriesMsg struct{}

	// OpenItemsMsg сообщение используется для открытия списка элементов определенной категории.
	OpenItemsMsg struct {
		// Category содержит категорию элементов, которые надо вывести.
		Category client.Category
	}

	// BackMsg используется для возврата в предыдущую модель, например, из элемента в список элементов.
	BackMsg struct {
		// State состояние модели, в которое надо вернуться.
		State int
		// Item каким элементом надо заменить данные в модели, в которую возвращаемся.
		Item *client.Item
	}

	// OpenItemMsg содержит информацию, какой элемент надо открыть.
	OpenItemMsg struct {
		// BackModel модель, в которую надо вернуться, в случае отмены.
		BackModel tea.Model
		// BackState состояние, в которое надо вернуть модель, в случае отмены.
		BackState int
		// Item каким элементом надо заменить данные в модели, в которую возвращаемся.
		Item *client.Item
	}
)
