package main

import (
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/setup"

	"github.com/fahmed8383/SchedulingApp/api/routes"
)

func main() {

	// get essentials struct and secrets struct by reference. the essentials struct contains the logger and the database, while the secrets struct
	// contains all credentials that can not be seen on the public repo
	ess, secrets := setup.GetEssentials()

	ess.Log.Info("api container built, starting http server")

	// create the neccesary schemas and tables for the postgres database
	SetUpDBStructure(ess)

	// create the route handlers for the http server -------------------------------------------------------------------------------
	// each call is logged in debug

	// registers /api/send-email-verification route
	http.HandleFunc("/api/send-email-verification", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.SendEmailVerification(w, r, ess, secrets)
	})

	// registers /api/resend-email-verification route
	http.HandleFunc("/api/resend-email-verification", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.ResendEmailVerification(w, r, ess, secrets)
	})

	// registers /api/check-username-availability route
	http.HandleFunc("/api/check-username-availability", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.CheckUsernameAvailability(w, r, ess)
	})

	// registers /api/check-email-availability route
	http.HandleFunc("/api/check-email-availability", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.CheckEmailAvailability(w, r, ess)
	})

	// registers /api/cancel-registration route
	http.HandleFunc("/api/cancel-registration", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.CancelRegistration(w, r, ess)
	})

	// registers /api/validate-verification-code
	http.HandleFunc("/api/validate-verification-code", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.ValidateVerificationCode(w, r, ess)
	})

	// registers /api/login
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.Login(w, r, ess, secrets)
	})

	// registers /api/logout
	http.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.LogOut(w, r, ess, secrets)
	})

	// registers /api/send-password-reset-email
	http.HandleFunc("/api/send-password-reset-email", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.SendPasswordResetEmail(w, r, ess, secrets)
	})

	// registers /api/reset-user-password
	http.HandleFunc("/api/reset-user-password", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.ResetUserPassword(w, r, ess, secrets)
	})

	// registers /api/my-schedule
	http.HandleFunc("/api/my-schedule", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.MySchedule(w, r, ess, secrets)
	})

	// registers /api/get-user-info
	http.HandleFunc("/api/get-user-info", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.GetUserInfo(w, r, ess, secrets)
	})

	// registers /api/reset-email
	http.HandleFunc("/api/reset-email", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.ResetEmail(w, r, ess, secrets)
	})

	// registers /api/reset-loggedin-password
	http.HandleFunc("/api/reset-loggedin-password", func(w http.ResponseWriter, r *http.Request) {
		ess.Log.Debugf("%s request on %s", r.Method, r.RequestURI)
		routes.ResetLoggedinPassword(w, r, ess, secrets)
	})

	//-----------------------------------------------------------------------------------------------------------------------------

	// start the  http server
	http.ListenAndServe(":6060", nil)
	ess.Log.Info("http server built and listening")
}

// SetUpDBStructure sets up the schemas and the tables in the postgres database if they have not already been created
func SetUpDBStructure(ess *setup.Essentials) {

	// create the base schema for the application if it does not already exist.
	// since this is a relatively small app, more than one schema would not be needed
	_, err := ess.PG.Exec(`CREATE SCHEMA IF NOT EXISTS app`)
	if err != nil {
		ess.Log.Error("unable to create app schema for database ", err)
	}

	// creates a users table if it does not exist, this table stores all basic info about the user including an encrypted password credentials
	_, err = ess.PG.Exec(`CREATE TABLE IF NOT EXISTS app.users (
		username TEXT UNIQUE NOT NULL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password BYTEA NOT NULL,
		verified TEXT NOT NULL,
		token TEXT,
		schedule json NOT NULL,
		sorting TEXT NOT NULL
	)`)
	if err != nil {
		ess.Log.Error("unable to create users table for database ", err)
	}

	// creates a login table if it does not exist, this table stores the jwt for each log in to be used to refresh tokens
	_, err = ess.PG.Exec(`CREATE TABLE IF NOT EXISTS app.login (
		username TEXT NOT NULL REFERENCES app.users,
		sessionid TEXT UNIQUE NOT NULL,
		jwt TEXT NOT NULL
	)`)
	if err != nil {
		ess.Log.Error("unable to create login table for database ", err)
	}

}
