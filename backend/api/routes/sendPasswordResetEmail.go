package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/auth"
	"github.com/fahmed8383/SchedulingApp/libraries/email"

	"github.com/fahmed8383/SchedulingApp/libraries/setup"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
)

// SendPasswordResetEmail is responsible for setting a temporary password for the user and sending a password reset link to the user
func SendPasswordResetEmail(w http.ResponseWriter, r *http.Request, ess *setup.Essentials, secrets *setup.Secrets) {

	// make sure the request is a post request
	if r.Method != "POST" {
		ess.Log.Error("method not POST request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ess.Log.Error("cannot read request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// returns the parsed request
	dataStruct, err := api.ParseRegInfo(body)
	if err != nil {
		ess.Log.Error("cannot parse request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// check to make sure none of the required fields are missing
	if dataStruct.Email == "" {
		ess.Log.Error("missing fields ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// queries for the username associated with the email
	sql := `SELECT username FROM app.users WHERE email = $1;`
	err = ess.PG.QueryRow(sql, dataStruct.Email).Scan(&dataStruct.Username)
	if err != nil {
		ess.Log.Error("unable to query for token ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// gets a random 14 character authorization token that can be used as a reset token
	token := auth.GenerateURLFriendlyToken(14)

	// update the token for the user in the database
	sql = `UPDATE app.users SET token = $1 WHERE email = $2;`
	_, err = ess.PG.Exec(sql, token, dataStruct.Email)
	if err != nil {
		ess.Log.Error("unable to update token for user ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	url := "web.127.0.0.1.xip.io/password-reset/" + dataStruct.Username + "/" + token

	// send the password reset email to the user
	err = email.SendPasswordResetEmail(secrets, dataStruct, url)
	if err != nil {
		ess.Log.Error("unable to send password reset email ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	w.Write([]byte(`{"msg":"success"}`))
	return
}
