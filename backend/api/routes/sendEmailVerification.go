package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/auth"
	"github.com/fahmed8383/SchedulingApp/libraries/email"

	"github.com/fahmed8383/SchedulingApp/libraries/setup"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
)

// SendEmailVerification is responsible for validating the registration recaptcha and sending a confirmation email to the user if everything is correct
func SendEmailVerification(w http.ResponseWriter, r *http.Request, ess *setup.Essentials, secrets *setup.Secrets) {

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
	if dataStruct.Username == "" || dataStruct.Email == "" || dataStruct.Password == "" {
		ess.Log.Error("missing fields ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// recaptcha verification --------------------------------------------------------------------

	// generate url for the request
	secret := secrets.Recaptcha
	url := "https://www.google.com/recaptcha/api/siteverify?secret=" + secret + "&response=" + dataStruct.Captcha
	method := "POST"

	// init http client
	client := &http.Client{}

	// create the http request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		ess.Log.Error("unable to generate http request to verify recaptcha ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// send the http request
	res, err := client.Do(req)
	if err != nil {
		ess.Log.Error("unable to send http request to verify recaptcha ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}
	defer res.Body.Close()

	// read response body
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		ess.Log.Error("unable to read body from recaptcha response ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// returns the parsed request
	recaptchaResp, err := api.ParseRecaptchaResp(body)
	if err != nil {
		ess.Log.Error("unable to parse body from recaptcha response ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// ------------------------------------------------------------------------------------------

	// check if the user passed the recaptcha
	if recaptchaResp.Success != true {
		ess.Log.Error("recaptcha failed ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	// gets a random 6 character authorization token that can be used to verify the email
	token := auth.GenerateToken(6)

	// send the verification code email to the user
	err = email.SendVerificationEmail(secrets, dataStruct, token)
	if err != nil {
		ess.Log.Error("unable to send verification code email ", err)
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

	// insert user data into the database, use variables so sql driver can escape all user entered info
	sql := `INSERT INTO app.users (username, email, password, verified, schedule, sorting) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = ess.PG.Exec(sql, dataStruct.Username, dataStruct.Email, encrypted, token, `[]`, "desc-due")
	if err != nil {
		ess.Log.Error("unable to add data to table ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"msg":"failure"}`))
		return
	}

	w.Write([]byte(`{"msg":"success"}`))
	return
}
