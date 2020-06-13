package setup

import (
	"database/sql"
	"time"

	//import postgresql library as a blank import
	_ "github.com/lib/pq"
	logger "github.com/sirupsen/logrus"
)

// Essentials holds all essential components of the backend in a struct that is passed by reference between packages
type Essentials struct {
	Log *logger.Logger
	PG  *sql.DB
}

// GetEssentials creates the essential struct that is passed by reference between packages
func GetEssentials() (*Essentials, *Secrets) {
	var ess Essentials

	// create the logger
	ess.Log = logger.New()
	ess.Log.SetLevel(logger.DebugLevel)

	ess.Log.Info("Starting api container")

	// read secrets/credentials from json file and parse it into a golang struct
	secrets, err := GetSecrets()
	if err != nil {
		ess.Log.Error("unable to correctly parse secrets.json ", err)
	}

	// wait a few seconds for postgresql container to startup
	time.Sleep(10 * time.Second)

	// connect to the postgres database
	db, err := sql.Open("postgres", "host=postgresql port=5432 user="+secrets.DbUser+" password="+secrets.DbPassword+" dbname=postgresql sslmode=disable")
	if err != nil {
		ess.Log.Error("unable to connect to postgres database ", err)
	}
	ess.PG = db

	// check to make sure the postgres database is alive
	err = ess.PG.Ping()
	if err != nil {
		ess.Log.Error("unable to ping postgres database ", err)
	}

	// return both the essentials and the secrets struct back to the main function
	return &ess, secrets
}
