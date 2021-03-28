package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"html/template"
	"log"
	"math"
	"reflect"
	"strings"
	"time"
)

type FormField struct {
	Name      string
	Data      interface{}
	InputType string
}

// for development we want to load the templates for each request
// for production we should cache the templates
func LoadBaseTemplates() (*template.Template, error) {

	funcMap := template.FuncMap{
		"now": time.Now,
		"toHTML": func(s string) (template.HTML, error) {
			return template.HTML(s), nil
		},
		"toJS": func(s string) template.JS {
			return template.JS(s)
		},
		// https://github.com/gohugoio/hugo/blob/edc5c4741caaee36ba4d42b5947c195a3e02e6aa/tpl/encoding/encoding.go#L60
		"jsonify": func(args ...interface{}) (template.HTML, error) {
			var (
				b   []byte
				err error
			)

			switch len(args) {
			case 0:
				return "", nil
			case 1:
				b, err = json.MarshalIndent(args[0], "", "    ")
			case 2:
				var opts map[string]string

				opts, err = cast.ToStringMapStringE(args[0])
				if err != nil {
					break
				}

				b, err = json.MarshalIndent(args[1], opts["prefix"], opts["indent"])
			default:
				err = errors.New("too many arguments to jsonify")
			}

			if err != nil {
				return "", err
			}

			return template.HTML(b), nil
		},

		// based loosely on:
		// https://github.com/gohugoio/hugo/blob/ce96895debb67df20ae24fb5f0f04b98a30cc6cc/tpl/math/math.go#L129
		"round3": func(x interface{}) (float64, error) {
			xf, err := cast.ToFloat64E(x)
			if err != nil {
				return 0, errors.New("Round operator can't be used with non-float value")
			}
			xf *= 1000

			return math.Round(xf) / 1000.0, nil
		},
		// based loosely on
		//https://github.com/gohugoio/hugo/blob/8a26ab0bc5dd9fa34e1362681fc08b0e522cd4ea/tpl/strings/strings.go#L398
		"trim": func(s interface{}) (string, error) {
			ss, err := cast.ToStringE(s)
			if err != nil {
				return "", err
			}

			return strings.TrimSpace(ss), nil
		},
		"member": func(s interface{}, col string) (string, error) {

			marshal, err := json.Marshal(s)
			if err != nil {
				return "", err
			}
			var data map[string]interface{}
			err = json.Unmarshal(marshal, &data)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%v", data[col]), nil
		},
		"members": func(s interface{}) (map[string]FormField, error) {

			marshal, err := json.Marshal(s)
			if err != nil {
				return nil, err
			}
			var conversion map[string]interface{}
			err = json.Unmarshal(marshal, &conversion)
			if err != nil {
				return nil, err
			}
			data := map[string]FormField{}
			for name, field := range conversion {

				fie := FormField{
					Name: "name",
					Data: field,
				}

				f, ok := reflect.TypeOf(s).FieldByName(name)
				if ok {
					// if the field has a gorm tag
					if len(f.Tag.Get("form_type")) > 0 {
						fie.InputType = f.Tag.Get("form_type")
					}
					log.Printf("%s Tag: %s", name, string(f.Tag))
				}
				data[name] = fie
			}

			return data, nil
		},
	}

	tmpl := template.New("_layout.html")

	tmpl = tmpl.Funcs(funcMap)

	partialGlob := "partials/*.html"
	var err error
	tmpl, err = tmpl.ParseGlob(partialGlob)
	if err != nil {
		return nil, fmt.Errorf("unable to parse glob (%s): %w", partialGlob, err)
	}

	return tmpl, nil
}
