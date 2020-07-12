package flash

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/lexkong/log"
)

var Store = sessions.NewCookieStore([]byte("secret-password"))
var sessionName = "flash-session"

//GetCurrentUserName returns the username of the logged in user
func GetCurrentUserName(r *http.Request) string {
	session, err := Store.Get(r, "session")
	if err == nil {
		return session.Values["username"].(string)
	}
	return ""
}

func SetFlashMessage(w http.ResponseWriter, r *http.Request, name string, value string) {
	session, err := Store.Get(r, flashName)
	if err != nil {
		log.Warnf("[session] set flash message err: %v", err)
	}
	session.AddFlash(value, name)
	session.Save(r, w)
}

func GetFlashMessage(w http.ResponseWriter, r *http.Request, name string) string {
	session, err := Store.Get(r, flashName)
	if err != nil {
		log.Warnf("[session] get flash message err: %v", err)
		return ""
	}

	fm := session.Flashes(name)
	if fm == nil {
		fmt.Fprint(w, "No flash messages")
		return ""
	}
	session.Save(r, w)
	fmt.Fprintf(w, "%v", fm[0])

	return fm[0].(string)
}
