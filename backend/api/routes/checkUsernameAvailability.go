package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/setup"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
)

// CheckUsernameAvailability is responsible for ensuring that the username has not been previously registered and is still available
func CheckUsernameAvailability(w http.ResponseWriter, r *http.Request, ess *setup.Essentials) {

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

	// queries for the username to check if it exists
	var exists bool
	sql := `SELECT exists (SELECT 1 FROM app.users WHERE username = $1 and verified = 'yes');`
	err = ess.PG.QueryRow(sql, dataStruct.Username).Scan(&exists)
	if err != nil {
		ess.Log.Error("unable to query for username ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	// if username exists return success, else retrun failure
	if exists {
		w.Write([]byte(`{"msg":"success"}`))
	} else {
		w.Write([]byte(`{"msg":"failure"}`))
	}

	return
}
