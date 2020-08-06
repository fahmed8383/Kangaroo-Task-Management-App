package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/setup"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
)

// ResetUserPassword is responsible for resetting a users password if the username and token provided are correct
func ResetUserPassword(w http.ResponseWriter, r *http.Request, ess *setup.Essentials, secrets *setup.Secrets) {

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
	dataStruct, err := api.ParseResetInfo(body)
	if err != nil {
		ess.Log.Error("cannot parse request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// check to make sure none of the required fields are missing
	if dataStruct.Username == "" || dataStruct.Password == "" || dataStruct.Token == "" {
		ess.Log.Error("missing fields ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// queries for the reset token associated with the username
	var token string
	sql := `SELECT token FROM app.users WHERE username = $1;`
	err = ess.PG.QueryRow(sql, dataStruct.Username).Scan(&token)
	if err != nil {
		ess.Log.Error("unable to query for token ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	if token != dataStruct.Token {
		ess.Log.Error("reset token is invalid ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// encrypt the password before storing in database
	encrypted, err := api.Encrypt(secrets.Key, dataStruct.Password)
	if err != nil {
		ess.Log.Error("unable to encrypt password ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// update the password and token for the user in the database
	sql = `UPDATE app.users SET password = $1, token = NULL WHERE username = $2;`
	_, err = ess.PG.Exec(sql, encrypted, dataStruct.Username)
	if err != nil {
		ess.Log.Error("unable to update password and token for user ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	w.Write([]byte(`{"msg":"success"}`))
	return
}
