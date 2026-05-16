package tui

func (m *model) switchAuth() {
	if m.form.focused == 0 {
		m.form.username.Blur()
		m.form.focused = 1
		m.form.password.Focus()
	} else {
		m.form.password.Blur()
		m.form.focused = 0
		m.form.username.Focus()
	}
}

func (m *model) switchData() {
	switch m.form.focused {
	case 0:
		m.form.meta.Blur()
		m.form.focused = 1
		m.form.recordType.Focus()
	case 1:
		m.form.recordType.Blur()
		m.form.focused = 2
		m.form.data.Focus()
	case 2:
		m.form.data.Blur()
		m.form.focused = 0
		m.form.meta.Focus()
	}
}

func (m *model) resetAuthForm() {
	m.form.username.SetValue("")
	m.form.password.SetValue("")
	m.form.focused = 0
	m.form.username.Focus()
	m.form.password.Blur()
}

func (m *model) resetData() {
	m.form.meta.SetValue("")
	m.form.recordType.SetValue("")
	m.form.data.SetValue("")
	m.form.focused = 0
	m.form.meta.Focus()
}
