package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
	"github.com/fahmed8383/SchedulingApp/libraries/auth"
	"github.com/fahmed8383/SchedulingApp/libraries/setup"
)

// Login is responsible for ensuring that the validation code submitted is correct and setting the user status to verified
func Login(w http.ResponseWriter, r *http.Request, ess *setup.Essentials, secrets *setup.Secrets) {

	// make sure the request is a post request
	if r.Method != "POST" {
		ess.Log.Error("method not POST request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ess.Log.Error("cannot read request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// returns the parsed request
	dataStruct, err := api.ParseLoginInfo(body)
	if err != nil {
		ess.Log.Error("cannot parse request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// queries for the password of the specific user
	var password string
	sql := `SELECT password FROM app.users WHERE username = $1;`
	err = ess.PG.QueryRow(sql, dataStruct.Username).Scan(&password)
	if err != nil {
		ess.Log.Error("unable to query for username ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// decrypt the password
	password, err = api.Decrypt(secrets.Key, password)
	if err != nil {
		ess.Log.Error("unable to decrypt password ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// check to see if the password matches
	if dataStruct.Password != password {
		ess.Log.Error("password is incorrect ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// generate a valid jwt token for the user
	token, err := auth.GenerateJwt(dataStruct, secrets.Jwt)
	if err != nil {
		ess.Log.Error("unable to generate jwt token ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// gets a random 14 character authorization token that can be used to designate a session id
	sessionid := auth.GenerateToken(14)

	// insert login data into the database, use variables so sql driver can escape all user entered info.

	// this is used as a subtitute for refresh token, when a token is expired the database will be checked for previous
	// valid token before sending out a new token, if token is invalid it will log the user out user from that session.

	// this will only trigger if the refresh token has been compromised, thus it will send me an email so I am aware.
	sql = `INSERT INTO app.login (username, sessionid, jwt) VALUES ($1, $2, $3);`
	_, err = ess.PG.Exec(sql, dataStruct.Username, sessionid, token)
	if err != nil {
		ess.Log.Error("unable to add data to jwt table ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// creates a jwt cookie for the user
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		MaxAge:   604800, // max age of a week
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	// create a session id cookie for the user
	cookie = http.Cookie{
		Name:   "sessionid",
		Value:  sessionid,
		Path:   "/",
		MaxAge: 604800, // max age of a week
	}
	http.SetCookie(w, &cookie)

	w.Write([]byte(`{"msg":"success"}`))
	return
}
