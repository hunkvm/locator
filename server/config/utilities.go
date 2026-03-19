package config

import (
	"reflect"
	"strings"
	"unicode"
)

func ReflectToMap(v any) map[string]any {
	tmp := convertValue(reflect.ValueOf(v))
	return tmp.(map[string]any)
}

func KebabKeys(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))

	for k, v := range m {
		out[kebabCase(k)] = kebabValue(v)
	}

	return out
}

func convertValue(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		out := make(map[string]any)

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue
			}
			key := field.Name
			out[key] = convertValue(v.Field(i))
		}

		return out

	case reflect.Slice, reflect.Array:
		out := make([]any, v.Len())

		for i := 0; i < v.Len(); i++ {
			out[i] = convertValue(v.Index(i))
		}

		return out

	case reflect.Map:
		out := make(map[string]any)

		for _, key := range v.MapKeys() {
			if key.Kind() != reflect.String {
				continue
			}

			out[key.String()] = convertValue(v.MapIndex(key))
		}

		return out

	default:
		return v.Interface()
	}
}

func kebabCase(s string) string {
	if strings.TrimSpace(s) == "" {
		return s
	}

	var b strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		// ignore _
		if r == ' ' || r == '-' {
			if b.Len() > 0 {
				b.WriteRune('-')
			}
			continue
		}

		if unicode.IsUpper(r) {
			hasPrev := i > 0
			hasNext := i+1 < len(runes)

			prev := rune(0)
			next := rune(0)

			if hasPrev {
				prev = runes[i-1]
			}
			if hasNext {
				next = runes[i+1]
			}

			if b.Len() > 0 && ((unicode.IsLower(prev) || unicode.IsDigit(prev)) ||
				(unicode.IsUpper(prev) && hasNext && unicode.IsLower(next))) {
				b.WriteRune('-')
			}

			b.WriteRune(unicode.ToLower(r))
			continue
		}

		b.WriteRune(unicode.ToLower(r))
	}

	return strings.Trim(b.String(), "-")
}

func kebabValue(v any) any {
	switch x := v.(type) {
	case map[string]any:
		return KebabKeys(x)

	case []any:
		out := make([]any, len(x))
		for i, item := range x {
			out[i] = kebabValue(item)
		}
		return out

	default:
		return v
	}
}
