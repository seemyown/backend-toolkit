package ctxbinding

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

func assignField(field reflect.Value, val any) error {
	typ := field.Type()

	parser, ok := getParser(typ)
	if !ok {
		return fmt.Errorf("no parser registered for type %s", typ.String())
	}

	parsed, err := parser(val)
	if err != nil {
		return err
	}

	field.Set(reflect.ValueOf(parsed))
	return nil
}

// Bind - binding user struct from fiber.Ctx
func Bind(ctx *fiber.Ctx, target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("bind expects a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		var raw string
		var found bool
		var defValue string

		switch {
		case fieldType.Tag.Get("ctx") != "":
			val := ctx.Locals(fieldType.Tag.Get("ctx"))
			if val == nil {
				raw = ""
			} else {
				raw = fmt.Sprintf("%v", val)
			}
			found = raw != ""

		case fieldType.Tag.Get("path") != "":
			raw = ctx.Params(fieldType.Tag.Get("path"))
			found = raw != ""

		case fieldType.Tag.Get("query") != "":
			raw = ctx.Query(fieldType.Tag.Get("query"))
			found = raw != ""
		}

		defValue = getDefaultValue(fieldType)

		// Get default if exists
		if !found {
			if defValue != "" {
				raw = defValue
				found = true
			}
		}

		// Parsing value
		if found {
			parser, ok := getParser(field.Type())
			if !ok {
				return fmt.Errorf("no parser for type %s", field.Type())
			}
			parsed, err := parser(raw)
			if err != nil {
				// If parse had error, check default value
				if defValue != "" {
					if parsed, err = parser(defValue); err != nil {
						return fmt.Errorf("field %s: %w", fieldType.Name, err)
					}
				} else {
					return fmt.Errorf("field %s: %w", fieldType.Name, err)
				}
			}
			field.Set(reflect.ValueOf(parsed))
		}

		// Проверка на required
		if fieldType.Tag.Get("required") == "true" && isZero(field) {
			return NewBindError(fmt.Sprintf("missing required field: %s", fieldType.Name), fieldType.Name)
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

func getDefaultValue(f reflect.StructField) string {
	return f.Tag.Get("default")
}
