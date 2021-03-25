package helpers

import (
	"log"
	"net/http"
)

func DoRender(res http.ResponseWriter, req *http.Request, template string, data interface{}) {

	session, err := GetSession(req)
	if err != nil {
		panic(err) //@todo
	}

	currentUser := ""
	if session.IsLoggedIn(req) {

		currentUser, err = session.GetCurrentUser(req)
		if err != nil {
			log.Printf("unable to get current user: %s", err)
		}
	}

	log.Printf("%s '%s' - %s", req.Method, req.URL.String(), currentUser)

	tmpl, err := LoadBaseTemplates()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Panicf("unable to LoadBaseTemplates: %s", err)
		return
	}

	file := "controllers/" + template
	tmpl, err = tmpl.ParseFiles(file)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Panicf("unable to parse file (%s): %s", file, err)
		return
	}

	//if er != nil {
	//	use := UserSafeError{}
	//	if errors.As(err, &use) {
	//		//@todo ideally this shows as a flash message, gets injected into the template
	//		res.WriteHeader(http.StatusInternalServerError)
	//		log.Panicf("unable to execute route function: %s", err)
	//	} else {
	//		res.WriteHeader(http.StatusInternalServerError)
	//		log.Panicf("unable to execute route function: %s", err)
	//	}
	//	return
	//}

	//res.WriteHeader(http.StatusOK)
	err = tmpl.Execute(res, data)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Panicf("unable to execute template: %s", err)
		return
	}
}
