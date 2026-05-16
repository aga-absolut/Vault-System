package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// registerHandler handles user registration form submission.
func (m *model) registerHandler() tea.Cmd {
	if m.form.focused == 0 {
		m.switchAuth()
		return nil
	}
	return func() tea.Msg {
		username := m.form.username.Value()
		password := m.form.password.Value()

		if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
			return errMsg("логин и пароль не могут быть пустыми")
		}

		err := m.client.Register(context.Background(), username, password)
		if err != nil {
			return errMsg(err.Error())
		}
		return successMsg("Регистрация прошла успешно!")
	}
}

// loginHandler handles user login form submission.
func (m *model) loginHandler() tea.Cmd {
	if m.form.focused == 0 {
		m.switchAuth()
		return nil
	}
	return func() tea.Msg {
		username := m.form.username.Value()
		password := m.form.password.Value()

		if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
			return errMsg("логин и пароль не могут быть пустыми")
		}

		err := m.client.Login(context.Background(), username, password)
		if err != nil {
			return errMsg(err.Error())
		}
		return successMsg("Авторизация прошла успешно!")
	}
}

// setDataHandler handles adding new user data.
func (m *model) setDataHandler() tea.Cmd {
	if m.form.focused < 2 {
		m.switchData()
		return nil
	}
	return func() tea.Msg {
		meta := m.form.meta.Value()
		recordType := m.form.recordType.Value()
		data := m.form.data.Value()

		if strings.TrimSpace(meta) == "" || strings.TrimSpace(recordType) == "" || strings.TrimSpace(data) == "" {
			return errMsg("данные не могут быть пустыми")
		}

		err := m.client.SetData(context.Background(), recordType, meta, []byte(data))
		if err != nil {
			return errMsg(err.Error())
		}
		return successMsg("Данные успешно добавлены!")
	}
}

// updateDataHandler handles updating existing user data.
func (m *model) updateDataHandler() tea.Cmd {
	if m.form.focused < 2 {
		m.switchData()
		return nil
	}
	return func() tea.Msg {
		meta := m.form.meta.Value()
		recordType := m.form.recordType.Value()
		data := m.form.data.Value()

		if strings.TrimSpace(meta) == "" || strings.TrimSpace(recordType) == "" || strings.TrimSpace(data) == "" {
			return errMsg("данные не могут быть пустыми")
		}

		err := m.client.UpdateData(context.Background(), recordType, meta, []byte(data))
		if err != nil {
			return errMsg(err.Error())
		}
		return successMsg("Данные успешно обновлены!")
	}
}

// getDataHandler handles retrieving user data by metadata.
func (m *model) getDataHandler() tea.Cmd {
	return func() tea.Msg {
		meta := m.form.meta.Value()

		if strings.TrimSpace(meta) == "" {
			return errMsg("метаданные не могут быть пустыми")
		}

		record, err := m.client.GetData(context.Background(), meta)
		if err != nil {
			return errMsg(err.Error())
		}
		return successMsg(fmt.Sprintf("Type: %s\ndata: %s", record.Type, string(record.Data)))
	}
}

// deleteDataHandler handles deleting user data.
func (m *model) deleteDataHandler() tea.Cmd {
	return func() tea.Msg {
		meta := m.form.meta.Value()

		if strings.TrimSpace(meta) == "" {
			return errMsg("метаданные не могут быть пустыми")
		}

		err := m.client.DeleteData(context.Background(), meta)
		if err != nil {
			return errMsg(err.Error())
		}
		return successMsg("Данные успешно удалены!")
	}
}

// getListDataHandler handles retrieving the list of stored metadata.
func (m *model) getListDataHandler() tea.Cmd {
	return func() tea.Msg {
		records, err := m.client.GetListMeta(context.Background())
		if err != nil {
			return errMsg(err.Error())
		}
		if len(records) == 0 {
			return successMsg("список записей пуст")
		}

		var sb strings.Builder
		for i, meta := range records {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, meta))
		}

		return successMsg(sb.String())
	}
}
