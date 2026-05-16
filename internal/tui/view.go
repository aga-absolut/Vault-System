package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Width(70).Render("=== GOPHKEEPER TUI ===") + "\n\n")

	if m.message != "" {
		s.WriteString(m.message + "\n\n")
	}

	switch m.state {
	case stateRegister:
		s.WriteString(titleStyle.Render("=== РЕГИСТРАЦИЯ ===") + "\n\n")

		s.WriteString("Логин:    " + m.form.username.View() + "\n")
		s.WriteString("Пароль:   " + m.form.password.View() + "\n\n")

		if m.message != "" {
			s.WriteString(m.message + "\n\n")
		}

		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).
			Render("Enter — дальше | Tab — переключить поле | Esc — назад"))
		return s.String()

	case stateLogin:
		s.WriteString(titleStyle.Render("=== АВТОРИЗАЦИЯ ===") + "\n\n")

		s.WriteString("Логин:    " + m.form.username.View() + "\n")
		s.WriteString("Пароль:   " + m.form.password.View() + "\n\n")

		if m.message != "" {
			s.WriteString(m.message + "\n\n")
		}

		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).
			Render("Enter — дальше | Tab — переключить поле | Esc — назад"))
		return s.String()

	case stateGetListData:
		s.WriteString(titleStyle.Render("=== СПИСОК МЕТАДАННЫХ ===") + "\n\n")

		if m.message != "" {
			s.WriteString(m.message + "\n\n")
		}

		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).
			Render("Esc — назад"))
		return s.String()

	case stateGetData:
		s.WriteString(titleStyle.Render("=== ВЫВОД ДАННЫХ ===") + "\n\n")

		s.WriteString("Метаданные:    " + m.form.meta.View() + "\n\n")

		if m.message != "" {
			s.WriteString(m.message + "\n\n")
		}

		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).
			Render("Enter — добавить | Tab — переключить поле | Esc — назад"))
		return s.String()

	case stateSetData:
		s.WriteString(titleStyle.Render("=== ДОБАВЛЕНИЕ ДАННЫХ ===") + "\n\n")

		s.WriteString("Метаданные:    " + m.form.meta.View() + "\n")
		s.WriteString("Тип:   " + m.form.recordType.View() + "\n")
		s.WriteString("Данные:   " + m.form.data.View() + "\n\n")

		if m.message != "" {
			s.WriteString(m.message + "\n\n")
		}

		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).
			Render("Enter — дальше | Tab — переключить поле | Esc — назад"))
		return s.String()

	case stateUpdateData:
		s.WriteString(titleStyle.Render("=== ОБНОВЛЕНИЕ ДАННЫХ ===") + "\n\n")

		s.WriteString("Метаданные:    " + m.form.meta.View() + "\n")
		s.WriteString("Тип:   " + m.form.recordType.View() + "\n")
		s.WriteString("Данные:   " + m.form.data.View() + "\n\n")

		if m.message != "" {
			s.WriteString(m.message + "\n\n")
		}

		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).
			Render("Enter — дальше | Tab — переключить поле | Esc — назад"))
		return s.String()

	case stateDeleteData:
		s.WriteString(titleStyle.Render("=== УДАЛЕНИЕ ДАННЫХ ===") + "\n\n")

		s.WriteString("метаданные:    " + m.form.meta.View() + "\n\n")

		if m.message != "" {
			s.WriteString(m.message + "\n\n")
		}

		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).
			Render("Enter — дальше | Tab — переключить поле | Esc — назад"))
		return s.String()
	}

	// Меню
	instruction := "↑↓ или j/k — перемещение | Enter — выполнить | ctrl+c — выход"
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c9c9")).Render(instruction) + "\n\n")

	for i, item := range m.menuItems {
		if i == m.cursor {
			s.WriteString(selectedStyle.Render(" → "+item) + "\n")
		} else {
			s.WriteString(itemStyle.Render("   "+item) + "\n")
		}
	}

	return s.String()
}
