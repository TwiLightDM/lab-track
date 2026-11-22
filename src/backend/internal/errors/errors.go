package errs

import "errors"

var (
	ErrInvalidValue     = errors.New("не валидные данные")
	ErrEntityNotFound   = errors.New("сущность не найдена")
	ErrNoFieldsProvided = errors.New("не передано ни одно поле")

	ErrInvalidCredentials = errors.New("не валидные данные для входа")
	ErrForbidden          = errors.New("отказано в доступе")
	ErrUnknownRole        = errors.New("неизвестная роль")
	ErrInvalidJSON        = errors.New("невалидный JSON")

	ErrDuplicatedKey      = errors.New("такая сущность уже существует")
	ErrForeignKeyViolated = errors.New("ошибка внешнего ключа")

	ErrMissingToken = errors.New("отсутствует токен")
	ErrInvalidToken = errors.New("проблема с токеном")
	ErrTokenExpired = errors.New("время жизни токена исчерпано")
)
