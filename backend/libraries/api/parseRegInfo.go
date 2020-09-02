package api

import "encoding/json"

// RegInfo holds the registration info received from the frontend in a struct
type RegInfo struct {
	Username         string `json:"userName"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Captcha          string `json:"captcha"`
	VerificationCode string `json:"verificationCode"`
}

// ParseRegInfo unmarshalls the byte registration data into a golang struct
func ParseRegInfo(data []byte) (RegInfo, error) {
	var info RegInfo
	err := json.Unmarshal(data, &info)
	return info, err
}

// SetUserInfo repacks username and email into RegInfo struct and marshalls them into bytes
func SetUserInfo(username string, email string) ([]byte, error) {
	response := RegInfo{username, email, "", "", ""}
	res, err := json.Marshal(response)
	return res, err
}
