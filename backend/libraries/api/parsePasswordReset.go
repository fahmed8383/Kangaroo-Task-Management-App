package api

import "encoding/json"

// ResetInfo holds the password reset info received from the frontend in a struct
type ResetInfo struct {
	Username string `json:"userName"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ParseResetInfo unmarshalls the byte password reset data into a golang struct
func ParseResetInfo(data []byte) (ResetInfo, error) {
	var info ResetInfo
	err := json.Unmarshal(data, &info)
	return info, err
}
