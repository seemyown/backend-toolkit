package db

import (
	"errors"
	"github.com/lib/pq"
	"github.com/seemyown/backend-toolkit/btools/exc"
)

func MapPGError(err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return exc.NotFoundError("not_found", "Ресурс не найден")
		case "23505":
			return exc.ConflictError("already_exists", "Запись уже существует")
		case "23502":
			// not_null_violation
			return exc.ValidationError("not_null_violation", "", "Обязательное поле не может быть пустым")
		case "23514":
			// check_violation
			return exc.ValidationError("check_violation", "", "Нарушение ограничения")
		case "22001":
			// string_data_right_truncation
			return exc.ValidationError("string_too_long", "", "Строка слишком длинная")
		case "22P02":
			// invalid_text_representation
			return exc.ValidationError("invalid_format", "", "Неверный формат данных")
		case "22007":
			// invalid_datetime_format
			return exc.ValidationError("invalid_date_format", "", "Неверный формат даты")
		case "42P02":
			// undefined_parameter
			return exc.BadRequestError("undefined_parameter", "Передан неизвестный параметр")
		default:
			return exc.RepositoryError(err.Error())
		}
	}
	return exc.RepositoryError(err.Error())
}
