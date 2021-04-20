package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"html/template"
	"log"
	"math"
	"reflect"
	"strings"
	"time"
)

// for development we want to load the templates for each request
// for production we should cache the templates
func LoadBaseTemplates(c echo.Context) (*template.Template, error) {

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
		//"member": func(s interface{}, col string) (string, error) {
		//
		//	marshal, err := json.Marshal(s)
		//	if err != nil {
		//		return "", err
		//	}
		//	var data map[string]interface{}
		//	err = json.Unmarshal(marshal, &data)
		//	if err != nil {
		//		return "", err
		//	}
		//
		//	return fmt.Sprintf("%v", data[col]), nil
		//},

		"members": getMembers,
		"getColumns": func(s interface{}) ([]FormField, error) {

			// if there is no data, this might error... @todo
			v := reflect.ValueOf(s)
			//log.Printf("%#v", v.Index(0).Interface())
			members, err := getMembers(v.Index(0).Interface())
			if err != nil {
				return nil, err
			}
			return members, nil
			return nil, nil
		},

		// see https://echo.labstack.com/guide/routing/
		// Echo#Reverse(name string, params ...interface{}
		"pathFor": func(name string, params ...interface{}) string {
			uri := c.Echo().Reverse(name, params...)
			if len(uri) < 2 {
				//@todo it would be nice to scan all templates for URLs to see if any of the route names are missing
				log.Printf("Unable to generate URL for '%s'", name)
			}
			return uri
		},
		"getFlash": func() []interface{} {
			return GetFlashMessages(c)
		},
		"isLoggedIn": func() bool {
			return IsLoggedIn(c)
		},
		"getCurrentUser": func() (string, error) {
			return GetCurrentUser(c)
		},

		// args is a list of key, value pairs. Keys must be strings
		"createMap": func(args ...interface{}) (map[string]interface{}, error) {

			if len(args)%2 != 0 {
				return nil, fmt.Errorf("number of args must be a multiple of 2")
			}

			data := map[string]interface{}{}
			for i := 0; i < len(args); i += 2 {
				data[args[i].(string)] = args[i+1]
			}

			return data, nil
		},
		// check if the current field is valid
		"isValid": func(fieldName string, valData map[string]string) (string, error) {
			message, err := getValMessage(fieldName, valData)
			if err != nil {
				// the callers are currently logging
				return "", err
			}
			log.Printf("isvalid %s %s - %#v", fieldName, message, valData)
			if len(message) > 0 {
				// is-invalid is a boostrap class that triggers the validation to show
				return "is-invalid", nil
			}
			return "", nil

		},
		// get any validation message
		"getValMessage": getValMessage,
	}

	tmpl := template.New("_layout.html")

	tmpl = tmpl.Funcs(funcMap)

	partialGlob := "partials/*.html"
	tmpl, err := tmpl.ParseGlob(partialGlob)
	if err != nil {
		log.Printf("Error parsing partial templates - %s", err)
		return nil, fmt.Errorf("unable to parse partial templates (%s): %w", partialGlob, err)
	}

	autoformGlob := "partials/autoform/*.html"
	tmpl, err = tmpl.ParseGlob(autoformGlob)
	if err != nil {
		log.Printf("Error parsing autoform templates - %s", err)
		return nil, fmt.Errorf("unable to parse autoform templates (%s): %w", autoformGlob, err)
	}

	return tmpl, nil
}

func getValMessage(fieldName string, valData map[string]string) (string, error) {
	if valData == nil {
		// this is expected when there is no validation information, such as when
		// the field is first shown to the user before it is submitted
		log.Printf("validation data is nil for field '%s'", fieldName)
		return "", nil
	}

	s, ok := valData[fieldName]
	if !ok {
		// this is expected if there is nothing wrong with the fields value per the validation rules
		log.Printf("no entry for field '%s'", fieldName)
		return "", nil
	}
	return s, nil

}

type FormField struct {
	Name       string
	Data       interface{}
	InputType  string
	InputClass string
	Options    []FormFieldOption
}

type FormFieldOption struct {
	Label string
	Value string // this could be interface{} but will be rendered to string by the template anyway
}

func getMembers(s interface{}) ([]FormField, error) {

	data := []FormField{}

	//https://stackoverflow.com/a/57073506/429544
	v := reflect.ValueOf(s)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		//fmt.Printf("Field: %12s\t Type: %15s\t Value: %v\n", typeOfS.Field(i).Name, v.Field(i).Type(), v.Field(i).Interface())
		//fmt.Printf("Field: %-12s\t Type: %-15s\n", typeOfS.Field(i).Name, v.Field(i).Type())

		uiTag := typeOfS.Field(i).Tag.Get("ui")
		if uiTag == "-" {
			//log.Printf("skipping %s b/c %s", typeOfS.Field(i).Name, uiTag)
			continue
		}

		ff := FormField{
			Name:       typeOfS.Field(i).Name,
			Data:       v.Field(i).Interface(),
			InputType:  "text",
			InputClass: "form-control",
		}

		switch v.Field(i).Type().Kind() {

		case reflect.Bool:
			ff.InputType = "checkbox"
			ff.InputClass = "form-check-input"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64:
			ff.InputType = "number"
		case reflect.Slice, reflect.Array:
			if v.Field(i).Type().String() != "uuid.UUID" {
				//log.Printf("Skipping Slice or array for field - %s", typeOfS.Field(i).Name)
				continue
			}
		case reflect.String:
			ff.InputType = "text"
		case reflect.Struct:

			_, ok := v.Field(i).Interface().(fmt.Stringer)
			if !ok && v.Field(i).Type().String() != "gorm.DeletedAt" {
				//fmt.Printf("Recursing into %s\n", typeOfS.Field(i).Name)
				members, err := getMembers(v.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				//fmt.Printf("found members: %#v\n", members)
				for _, m := range members {
					data = append(data, m)
				}
				continue
			}
			fmt.Printf("found string method for field %s, assuming text\n", typeOfS.Field(i).Name)
			ff.InputType = "text"

		case reflect.Interface, reflect.Ptr:
			//fallthrough
		case reflect.Map:
			//fallthrough
		default:
			log.Printf("unsupported field with type %s", v.Field(i).Type().Kind())
			ff.InputType = "text"
		}

		if uiTag == "textarea" {
			ff.InputType = "textarea"
		}

		data = append(data, ff)

	}

	//v := reflect.ValueOf(s)
	//
	//for i := 0; i < v.NumField(); i++ {
	//	f := v.Field(i)
	//	t := reflect.TypeOf(f)
	//	log.Printf("Field: %s - %s", f.Type(), t.Name())
	//}
	return data, nil

	//// covert the thing (s) into a key:value map
	//marshal, err := json.Marshal(s)
	//if err != nil {
	//	return nil, err
	//}
	//
	////log.Printf("JSON: %s", marshal)
	//
	//var conversion map[string]interface{}
	//err = json.Unmarshal(marshal, &conversion)
	//if err != nil {
	//	return nil, err
	//}
	//
	////log.Printf("MAP: %#v", conversion)
	//
	//data := map[string]FormField{}
	//for name, field := range conversion {
	//	log.Printf("== Field: %s", name)
	//	fie := FormField{
	//		Name:      name,
	//		Data:      field,
	//		InputType: "text",
	//	}
	//
	//	f, ok := reflect.TypeOf(s).FieldByName(name)
	//	if !ok {
	//		log.Printf("Field: %s is not ok", name)
	//	} else {
	//		log.Printf("Field: %s is %s", name, f.Type.Name())
	//		if f.Type.Name() == "bool" {
	//			log.Printf("field is bool! %s", name)
	//		}
	//
	//		// if the field has a gorm tag
	//		if len(f.Tag.Get("form_type")) > 0 {
	//			fie.InputType = f.Tag.Get("form_type")
	//		}
	//		//log.Printf("%s Tag: %s", name, string(f.Tag))
	//
	//		if f.Tag.Get("ui") == "-" {
	//			fie.InputType = "none"
	//		}
	//	}
	//	data[name] = fie
	//}
	//
	//return data, nil
}

//func getColumns(s interface{}) {
//	v := reflect.ValueOf(s)
//}
