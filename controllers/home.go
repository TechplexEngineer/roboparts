package controllers

import (
	"github.com/techplexengineer/gorm-roboparts/helpers"
	"net/http"
)

func Home(res http.ResponseWriter, req *http.Request) {

	data := map[string]interface{}{}

	helpers.DoRender(res, req, "home.html", data)
}
