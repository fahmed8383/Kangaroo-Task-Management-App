package setup

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Secrets holds all info that needs to be kept secure and cant be availible on the public github. Loaded from an external file
type Secrets struct {
	Recaptcha  string `json:"recaptcha"`
	DbUser     string `json:"dbuser"`
	DbPassword string `json:"dbpassword"`
	Key        string `json:"key"`
	Jwt        string `json:"jwt"`
	GmailPass  string `json:"gmailPass"`
}

// GetSecrets unmarshalls and returns the secrets file by reference in a golang struct
func GetSecrets() (*Secrets, error) {

	// open secrets.json file
	jsonFile, err := os.Open("/secrets.json")

	// read the byte value of the secrets.json file
	byteValue, err := ioutil.ReadAll(jsonFile)

	// create a secrets struct variable and unmarshall the byte data into the struct
	var secrets Secrets
	err = json.Unmarshal(byteValue, &secrets)

	return &secrets, err
}
