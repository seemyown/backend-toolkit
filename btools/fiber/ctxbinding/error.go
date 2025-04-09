package ctxbinding

type BindError struct {
	Msg   string
	Field string
}

func (e *BindError) Error() string {
	return e.Field
}

func NewBindError(msg string, field string) *BindError {
	return &BindError{
		Msg:   msg,
		Field: field,
	}
}
