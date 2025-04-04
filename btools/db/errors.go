package db

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/lib/pq"
)

// RepositoryError представляет ошибку уровня хранилища (репозитория).
// Она содержит собственный код, сообщение для пользователя и оригинальную ошибку.
type RepositoryError struct {
	Code              int               // внутренний код ошибки (не SQLSTATE)
	LocalizedMessages map[string]string // понятное сообщение об ошибке
	Err               error             // исходная ошибка, возвращённая драйвером или другим уровнем
}

// Error удовлетворяет интерфейсу error, возвращая сообщение об ошибке.
func (e *RepositoryError) Error() string {
	return e.MessageFor("en")
}

func (e *RepositoryError) MessageFor(lang string) string {
	if msg, ok := e.LocalizedMessages[lang]; ok {
		return msg
	}
	return e.LocalizedMessages["ru"] // fallback
}

// Unwrap позволяет извлечь вложенную оригинальную ошибку (для errors.Is/As).
func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// Набор констант с внутренними кодами ошибок для RepositoryError.
const (
	ErrCodeNotFound                  = 10  // Например, запись не найдена (sql.ErrNoRows)
	ErrCodeUniqueViolation           = 11  // Нарушение уникального ограничения (дубликат)
	ErrCodeForeignKeyViolation       = 12  // Нарушение ограничения внешнего ключа
	ErrCodeNotNullViolation          = 13  // Нарушение ограничения NOT NULL
	ErrCodeCheckViolation            = 14  // Нарушение CHECK ограничения
	ErrCodeExclusionViolation        = 15  // Нарушение исключающего ограничения
	ErrCodeRestrictViolation         = 16  // Нарушение ограничения RESTRICT (FOREIGN KEY)
	ErrCodeStringTooLong             = 21  // Превышение допустимой длины строки
	ErrCodeNumericOutOfRange         = 22  // Числовое значение вне допустимого диапазона
	ErrCodeInvalidTextRepresentation = 23  // Неверный формат входных данных
	ErrCodeDeadlockDetected          = 31  // Обнаружен дедлок
	ErrCodeSerializationFailure      = 32  // Ошибка сериализации транзакции
	ErrCodeUnhandled                 = 999 // Неизвестная/необработанная ошибка
)

var MapToHttpError = map[int]int{
	ErrCodeNotFound:                  http.StatusNotFound,
	ErrCodeUniqueViolation:           http.StatusConflict,
	ErrCodeForeignKeyViolation:       http.StatusConflict,
	ErrCodeNotNullViolation:          http.StatusNotAcceptable,
	ErrCodeCheckViolation:            http.StatusNotAcceptable,
	ErrCodeExclusionViolation:        http.StatusConflict,
	ErrCodeRestrictViolation:         http.StatusNotAcceptable,
	ErrCodeStringTooLong:             http.StatusRequestEntityTooLarge,
	ErrCodeNumericOutOfRange:         http.StatusRequestedRangeNotSatisfiable,
	ErrCodeInvalidTextRepresentation: http.StatusUnprocessableEntity,
	ErrCodeDeadlockDetected:          http.StatusGatewayTimeout,
	ErrCodeSerializationFailure:      http.StatusInternalServerError,
	ErrCodeUnhandled:                 http.StatusInternalServerError,
}

func localized(messageRu, messageEn string) map[string]string {
	return map[string]string{
		"ru": messageRu,
		"en": messageEn,
	}
}

func WrapError(err error) *RepositoryError {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return &RepositoryError{
			Code:              ErrCodeNotFound,
			LocalizedMessages: localized("Запись не найдена", "Record not found"),
			Err:               err,
		}
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		code := ErrCodeUnhandled
		messages := localized("Необработанная ошибка", "Unhandled database error")

		switch pqErr.Code.Name() {
		case "unique_violation":
			code = ErrCodeUniqueViolation
			messages = localized("Нарушение уникального ограничения", "Unique constraint violated")
		case "foreign_key_violation":
			code = ErrCodeForeignKeyViolation
			messages = localized("Нарушение внешнего ключа", "Foreign key constraint violated")
		case "not_null_violation":
			code = ErrCodeNotNullViolation
			messages = localized("NULL в поле, где он запрещён", "NULL value where NOT NULL is required")
		case "check_violation":
			code = ErrCodeCheckViolation
			messages = localized("Нарушение CHECK ограничения", "CHECK constraint violated")
		case "exclusion_violation":
			code = ErrCodeExclusionViolation
			messages = localized("Нарушение EXCLUDE ограничения", "Exclusion constraint violated")
		case "restrict_violation":
			code = ErrCodeRestrictViolation
			messages = localized("Удаление запрещено (RESTRICT)", "Delete restricted due to FK")
		case "deadlock_detected":
			code = ErrCodeDeadlockDetected
			messages = localized("Обнаружен дедлок", "Deadlock detected")
		case "serialization_failure":
			code = ErrCodeSerializationFailure
			messages = localized("Ошибка сериализации транзакции", "Transaction serialization failure")
		case "string_data_right_truncation":
			code = ErrCodeStringTooLong
			messages = localized("Строка слишком длинная", "String too long for field")
		case "numeric_value_out_of_range":
			code = ErrCodeNumericOutOfRange
			messages = localized("Число вне допустимого диапазона", "Numeric value out of range")
		case "invalid_text_representation":
			code = ErrCodeInvalidTextRepresentation
			messages = localized("Неверный формат входных данных", "Invalid input format")
		}

		return &RepositoryError{
			Code:              code,
			LocalizedMessages: messages,
			Err:               err,
		}
	}

	return &RepositoryError{
		Code:              ErrCodeUnhandled,
		LocalizedMessages: localized("Необработанная ошибка", "Unhandled database error"),
		Err:               err,
	}
}
