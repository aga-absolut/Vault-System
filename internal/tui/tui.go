package tui

import (
	"github.com/aga-absolut/Vault-System/internal/client"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#3784ff")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#59b4ff")).Bold(true)
	itemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#59b4ff"))
)

type state string

const (
	stateMenu        state = "menu"
	stateRegister    state = "register"
	stateLogin       state = "login"
	stateGetData     state = "get data"
	stateSetData     state = "set data"
	stateGetListData state = "get list data"
	stateDeleteData  state = "delete data"
	stateUpdateData  state = "update data"
)

type (
	errMsg     string
	successMsg string
)

type registerForm struct {
	username   textinput.Model
	password   textinput.Model
	meta       textinput.Model
	recordType textinput.Model
	data       textinput.Model
	focused    int
}

type model struct {
	client    *client.Client
	menuItems []string
	cursor    int
	state     state

	form    registerForm
	message string
}

func InitialModel(client *client.Client) model {
	username := textinput.New()
	username.Placeholder = "логин"
	username.Focus()

	password := textinput.New()
	password.Placeholder = "пароль"
	password.EchoMode = textinput.EchoPassword

	meta := textinput.New()
	meta.Placeholder = "метаданные"
	meta.Focus()

	recordType := textinput.New()
	recordType.Placeholder = "тип данных"

	data := textinput.New()
	data.Placeholder = "данные"

	return model{
		client: client,
		state:  stateMenu,
		menuItems: []string{
			"1. 🔑 Регистрация",
			"2. 🔓 Вход в аккаунт",
			"3. 📋 Просмотреть все записи",
			"4. 🔍 Найти запись",
			"5. ➕ Добавить новую запись",
			"6. ✏️  Отредактировать запись",
			"7.  🗑 Удалить запись",
			"8. ❌ Выйти из программы",
		},
		form: registerForm{
			username:   username,
			password:   password,
			meta:       meta,
			recordType: recordType,
			data:       data,
			focused:    0,
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			switch m.state {
			case stateRegister, stateLogin:
				m.state = stateMenu
				m.message = ""
				m.resetAuthForm()
				return m, nil
			case stateGetData, stateDeleteData, stateSetData, stateUpdateData:
				m.state = stateMenu
				m.message = ""
				m.resetData()
				return m, nil
			case stateGetListData:
				m.state = stateMenu
				m.message = ""
				return m, nil
			}

		case "enter", " ":
			switch m.state {
			case stateMenu:
				return m, m.executeAction(m.cursor)
			case stateRegister:
				return m, m.registerHandler()
			case stateLogin:
				return m, m.loginHandler()
			case stateGetData:
				return m, m.getDataHandler()
			case stateDeleteData:
				return m, m.deleteDataHandler()
			case stateGetListData:
				return m, m.getListDataHandler()
			case stateSetData:
				return m, m.setDataHandler()
			case stateUpdateData:
				return m, m.updateDataHandler()
			}

		case "up":
			if m.state == stateMenu && m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.state == stateMenu && m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}

		case "tab", "shift+tab":
			if m.state == stateRegister || m.state == stateLogin {
				m.switchAuth()
				return m, nil
			}
			if m.state == stateSetData || m.state == stateUpdateData {
				m.switchData()
				return m, nil
			}
		}

	case errMsg:
		m.message = errorStyle.Render("ОШИБКА: " + string(msg))
		m.state = stateMenu
		m.resetAuthForm()
		m.resetData()
		return m, nil

	case successMsg:
		m.message = successStyle.Render(string(msg))
		m.state = stateMenu
		m.resetAuthForm()
		m.resetData()
		return m, nil
	}

	// === ОБНОВЛЕНИЕ ФОРМЫ РЕГИСТРАЦИИ ===
	switch m.state {
	case stateRegister, stateLogin:
		switch m.form.focused {
		case 0:
			m.form.username, cmd = m.form.username.Update(msg)
		case 1:
			m.form.password, cmd = m.form.password.Update(msg)
		}

	case stateGetData, stateDeleteData:
		m.form.meta, cmd = m.form.meta.Update(msg)

	case stateSetData, stateUpdateData:
		switch m.form.focused {
		case 0:
			m.form.meta, cmd = m.form.meta.Update(msg)
		case 1:
			m.form.recordType, cmd = m.form.recordType.Update(msg)
		case 2:
			m.form.data, cmd = m.form.data.Update(msg)
		}
	}

	return m, cmd
}

func (m *model) executeAction(index int) tea.Cmd {
	switch index {
	case 0:
		m.state = stateRegister
		m.message = ""
		m.resetAuthForm()
		return nil
	case 1:
		m.state = stateLogin
		m.message = ""
		m.resetAuthForm()
		return nil
	case 2:
		m.state = stateGetListData
		m.message = ""
		return m.getListDataHandler()
	case 3:
		m.state = stateGetData
	case 4:
		m.state = stateSetData
	case 5:
		m.state = stateUpdateData
	case 6:
		m.state = stateDeleteData
	case 7:
		return tea.Quit
	default:
		return nil
	}
	m.message = ""
	m.resetData()
	return nil
}
