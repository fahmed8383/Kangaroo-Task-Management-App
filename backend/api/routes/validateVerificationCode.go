package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
	"github.com/fahmed8383/SchedulingApp/libraries/setup"
)

// ValidateVerificationCode is responsible for ensuring that the validation code submitted is correct and setting the user status to verified
func ValidateVerificationCode(w http.ResponseWriter, r *http.Request, ess *setup.Essentials) {

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

	// queries for the token created for the specific user
	var token string
	sql := `SELECT verified FROM app.users WHERE username = $1;`
	err = ess.PG.QueryRow(sql, dataStruct.Username).Scan(&token)
	if err != nil {
		ess.Log.Error("unable to query for token ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// check to see if the verfication code matches
	if token != dataStruct.VerificationCode {
		ess.Log.Error("verification code is incorrect ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// update verified status to "yes" for the user in the database
	sql = `UPDATE app.users SET verified = $1 WHERE username = $2;`
	_, err = ess.PG.Exec(sql, "yes", dataStruct.Username)
	if err != nil {
		ess.Log.Error("unable to update verification status for user ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	w.Write([]byte(`{"msg":"success"}`))
	return
}
