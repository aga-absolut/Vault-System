package errs

import "errors"

// Application error definitions.
var (
	ErrIncorrectLoginOrPassword = errors.New("неверный логин или пароль")
	ErrLoginAlreadyUsed         = errors.New("логин уже используется")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrRecordNotFound           = errors.New("данные не найдены")
	ErrInternal                 = errors.New("ошибка сервера")
	ErrTooShortPassword         = errors.New("слишком короткий пароль")
	ErrMetaAlreadyUsed          = errors.New("метаданные уже используются")
)