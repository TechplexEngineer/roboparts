package helpers

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"net/http"
)

type Session struct {
	sc *securecookie.SecureCookie
}

//securecookie.GenerateRandomKey(64),
//securecookie.GenerateRandomKey(32)

func NewSess(hashKey, blockKey []byte) *Session {
	return &Session{
		sc: securecookie.New(hashKey, blockKey),
	}
}

// name of the cookie to set when a user is logged in
const cookieName = "session"

// structure for data to store in the cookie
type sessionData map[string]string

//type sessionData struct {
//	userUUID string
//}

func (sc *Session) StartSession(userid string, response http.ResponseWriter) {
	sessionData := sessionData{}
	sessionData["userUUID"] = userid

	if encoded, err := sc.sc.Encode(cookieName, sessionData); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func (sc *Session) ClearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func (sc *Session) GetCurrentUser(req *http.Request) (string, error) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("unable to get %s cookie - %w", cookieName, err)
	}

	cookieValue := make(map[string]string)
	err = sc.sc.Decode(cookieName, cookie.Value, &cookieValue)
	if err != nil {
		return "", fmt.Errorf("unable to decode %s cookie - %w", cookieName, err)
	}
	return cookieValue["userUUID"], nil
}

func (sc *Session) IsLoggedIn(request *http.Request) bool {
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return false
	}

	cookieValue := sessionData{}
	err = sc.sc.Decode(cookieName, cookie.Value, &cookieValue)
	if err != nil {
		return false
	}

	return true
}

func (sc Session) SetFlashMessage() {

}
