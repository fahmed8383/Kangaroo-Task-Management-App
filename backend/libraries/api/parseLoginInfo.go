package api

import "encoding/json"

// LoginInfo holds the login info received from the frontend in a struct
type LoginInfo struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

// ParseLoginInfo unmarshalls the byte login data into a golang struct
func ParseLoginInfo(data []byte) (LoginInfo, error) {
	var info LoginInfo
	err := json.Unmarshal(data, &info)
	return info, err
}
