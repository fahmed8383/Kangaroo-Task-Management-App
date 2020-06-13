package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
	"github.com/fahmed8383/SchedulingApp/libraries/auth"
	"github.com/fahmed8383/SchedulingApp/libraries/email"
	"github.com/fahmed8383/SchedulingApp/libraries/setup"
)

// ResendEmailVerification is responsible for re-sending a verification code email to the user
func ResendEmailVerification(w http.ResponseWriter, r *http.Request, ess *setup.Essentials, secrets *setup.Secrets) {

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
	dataStruct, err := api.ParseRegInfo(body)
	if err != nil {
		ess.Log.Error("cannot parse request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// check to make sure none of the required fields are missing
	if dataStruct.Username == "" || dataStruct.Email == "" || dataStruct.Password == "" {
		ess.Log.Error("missing fields ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// queries for the username to check if it exists
	var exists bool
	sql := `SELECT exists (SELECT 1 FROM app.users WHERE username = $1);`
	err = ess.PG.QueryRow(sql, dataStruct.Username).Scan(&exists)
	if err != nil {
		ess.Log.Error("unable to query for username ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// if username does not exist, cancel the request
	if !exists {
		ess.Log.Error("attempt to re-send verification for unregistered user ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// gets a random 6 character authorization token that can be used to verify the email
	token := auth.GenerateToken(6)

	// send the verification code email to the user
	err = email.SendVerificationEmail(secrets, dataStruct, token)
	if err != nil {
		ess.Log.Error("unable to send verification code email ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// update the token and email for the user in the database
	sql = `UPDATE app.users SET verified = $1, email = $2 WHERE username = $3;`
	_, err = ess.PG.Exec(sql, token, dataStruct.Email, dataStruct.Username)
	if err != nil {
		ess.Log.Error("unable to update token value for user ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	w.Write([]byte(`{"msg":"success"}`))
	return
}
