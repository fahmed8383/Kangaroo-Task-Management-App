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
