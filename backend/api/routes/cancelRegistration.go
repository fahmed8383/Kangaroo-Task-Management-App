package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
	"github.com/fahmed8383/SchedulingApp/libraries/setup"
)

// CancelRegistration is responsible for removing an unverified user from the database
func CancelRegistration(w http.ResponseWriter, r *http.Request, ess *setup.Essentials) {

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
	if dataStruct.Username == "" {
		ess.Log.Error("missing fields ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// delete the user from the database
	sql := `DELETE FROM app.users WHERE username = $1 AND verified != 'yes';`
	_, err = ess.PG.Exec(sql, dataStruct.Username)
	if err != nil {
		ess.Log.Error("unable to delete unverified user ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	w.Write([]byte(`{"msg":"success"}`))
	return
}
