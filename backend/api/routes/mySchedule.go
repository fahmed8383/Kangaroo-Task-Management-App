package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fahmed8383/SchedulingApp/libraries/api"
	"github.com/fahmed8383/SchedulingApp/libraries/auth"
	"github.com/fahmed8383/SchedulingApp/libraries/setup"
)

// MySchedule is responsible for saving and retreiving the user schedule and sorting method from the database
func MySchedule(w http.ResponseWriter, r *http.Request, ess *setup.Essentials, secrets *setup.Secrets) {

	// check to make sure jwt is valid
	if r.Method == "POST" || r.Method == "GET" {

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
		// this does not check to see if jwt has expired or not
		claims, token, err := auth.GetToken(jwtCookie.Value, secrets.Jwt)
		if err != nil {
			ess.Log.Error("invalid jwt ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		// check to make sure the jwt is for the correct user
		if claims.Username != usernameCookie.Value {
			ess.Log.Error("jwt user does not match cookie ")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		// if token is expired
		if !token.Valid {

			// queries to see if token matches the last token for the session in the database
			var jwt string
			sql := `SELECT jwt FROM app.login WHERE username = $1 and sessionid = $2;`
			err = ess.PG.QueryRow(sql, usernameCookie.Value, sessIDCookie.Value).Scan(&jwt)
			if err != nil {
				ess.Log.Error("unable to query for jwt ", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"msg":"failure"}`))
				return
			}

			// if the token saved in the database does not match the received token,
			// it is likely the token has been hijacked. Log out user from session and invalidate token
			if token.Raw != jwt {

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

				w.Write([]byte(`{"msg":"logged-out"}`))
				return
			}

			// if token matches the token in database, generate and save new token for user and continue with rest of request

			// generate a valid jwt token for the user
			token, err := auth.GenerateJwt(usernameCookie.Value, secrets.Jwt)
			if err != nil {
				ess.Log.Error("unable to generate jwt token ", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"msg":"failure"}`))
				return
			}

			// update the jwt for the session in the database
			sql = `UPDATE app.login SET jwt = $1 WHERE username = $2 and sessionid = $3;`
			_, err = ess.PG.Exec(sql, token, usernameCookie.Value, sessIDCookie.Value)
			if err != nil {
				ess.Log.Error("unable to update password and token for user ", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"msg":"failure"}`))
				return
			}

			// creates a jwt cookie for the user
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    token,
				Path:     "/",
				MaxAge:   2147483647, // virtually infinite cookie
				HttpOnly: true,
			})

		}

	}

	// if method is post, save the schedule in the database
	if r.Method == "POST" {

		// read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ess.Log.Error("cannot read request body ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		// returns the parsed request
		dataStruct, err := api.ParseScheduleInfo(body)
		if err != nil {
			ess.Log.Error("cannot parse request body ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		// gets username cookie from frontend
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			ess.Log.Error("username cookie missing ")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		// update schedule and sorting method for the user in the database
		sql := `UPDATE app.users SET schedule = $1, sorting = $2 WHERE username = $3;`
		_, err = ess.PG.Exec(sql, dataStruct.Schedule, dataStruct.SortingMethod, usernameCookie.Value)
		if err != nil {
			ess.Log.Error("unable to update verification status for user ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		w.Write([]byte(`{"msg":"success"}`))
		return
	}

	// if method is get, get the schedule for the frontend
	if r.Method == "GET" {

		// gets username cookie from frontend
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			ess.Log.Error("username cookie missing ")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		// queries for the schedule and sorting method of the specific user
		var schedule json.RawMessage
		var sorting string
		sql := `SELECT schedule, sorting FROM app.users WHERE username = $1;`
		err = ess.PG.QueryRow(sql, usernameCookie.Value).Scan(&schedule, &sorting)
		if err != nil {
			ess.Log.Error("unable to query for schedule and sorting method ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		// create response for frontend
		res, err := api.SetScheduleResponse(schedule, sorting)
		if err != nil {
			ess.Log.Error("unable to create response to send schedule to frontend ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"msg":"failure"}`))
			return
		}

		w.Write(res)
		return
	}

	ess.Log.Error("method not POST or GET")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"msg":"failure"}`))
	return
}
