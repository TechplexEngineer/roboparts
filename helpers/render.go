package helpers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	tmpl, err := LoadBaseTemplates()
	if err != nil {
		return fmt.Errorf("failed to load base templates - %w", err)
	}
	tmpl.Funcs(template.FuncMap{
		// see https://echo.labstack.com/guide/routing/
		// Echo#Reverse(name string, params ...interface{}
		"pathFor": func(name string, params ...interface{}) string {
			uri := c.Echo().Reverse(name, params...)
			if len(uri) < 2 {
				//@todo it would be nice to scan all templates for URLs to see if any of the route names are missing
				log.Warnf("Unable to generate URL for '%s'", name)
			}
			return uri
		},
	})

	cwd, _ := os.Getwd()
	callingPkg := path.Dir(getFrame(2).File)

	file := strings.TrimPrefix(callingPkg, cwd+string(os.PathSeparator))
	file += string(os.PathSeparator) + name

	tmpl, err = tmpl.ParseFiles(file)
	if err != nil {
		log.Printf("failed to load %s templates - %s", file, err)
		return fmt.Errorf("failed to load %s templates - %w", file, err)
	}

	// Add global methods if data is a map
	//if viewContext, isMap := data.(map[string]interface{}); isMap {
	//	viewContext["reverse"] = c.Echo().Reverse
	//	viewContext["routes"] = c.Echo().Routes()
	//}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("error executing template %s - %s", file, err)
	}

	return err
}

// source: https://gist.github.com/changkun/161979dfedc0cd7d65dd06fe83b73cdc
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}
