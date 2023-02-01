package helpers

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func RenderGoTemplate(t string, data any) (string, error) {
	tpl, err := template.
		New("template").
		Funcs(TemplateFuncs()).
		Parse(t)
	if err != nil {
		return "", err
	}

	var content bytes.Buffer
	if err := tpl.Execute(&content, data); err != nil {
		return "", err
	}
	return content.String(), nil
}

func TemplateFuncs() map[string]any {
	return map[string]any{
		"hasPrefix": strings.HasPrefix,
		"toLower":   strings.ToLower,
		"slugify":   Slugify,
		"renderProperty": func(key string, value any) string {
			return strings.TrimSuffix(SerializeToHCL(key, value), "\n")
		},
		"renderOptionalProperty": func(key string, value any) string {
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.String && v.IsZero() {
				return ""
			}
			if v.Kind() == reflect.Slice && v.IsNil() {
				return ""
			}

			return strings.TrimSuffix(SerializeToHCL(key, value), "\n")
		},
	}
}
