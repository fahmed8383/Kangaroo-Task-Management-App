package routes

import (
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/auth"
	"github.com/fahmed8383/SchedulingApp/libraries/setup"
)

// LogOut is responsible for removing the cookies for the user and the session from the database
func LogOut(w http.ResponseWriter, r *http.Request, ess *setup.Essentials, secrets *setup.Secrets) {

	// make sure the request is a post request
	if r.Method != "POST" {
		ess.Log.Error("method not POST request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// check to make sure the logged in cookies are present.
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		ess.Log.Error("jwt cookie missing ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	sessIDCookie, err := r.Cookie("sessionid")
	if err != nil {
		ess.Log.Error("sessionid cookie missing ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	usernameCookie, err := r.Cookie("username")
	if err != nil {
		ess.Log.Error("username cookie missing ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// check to make sure the jwt signature is valid and has not been tampered with
	// we do not care if token is expired or not since we would want to log out regardless
	// if it is not a jwt we issued, do no log out
	_, _, err = auth.GetToken(jwtCookie.Value, secrets.Jwt)
	if err != nil {
		ess.Log.Error("invalid jwt ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// delete the login entry for the session from the database
	sql := `DELETE FROM app.login WHERE username = $1 AND sessionid = $2;`
	_, err = ess.PG.Exec(sql, usernameCookie.Value, sessIDCookie.Value)
	if err != nil {
		ess.Log.Error("unable to delete login for user session ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// delete jwt cookie for the user
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// delete session id cookie for the user
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionid",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// delete session id cookie for the user
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Write([]byte(`{"msg":"success"}`))
	return
}
