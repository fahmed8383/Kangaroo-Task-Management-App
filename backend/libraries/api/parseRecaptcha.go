package api

import "encoding/json"

// RecaptchaResp holds the response info received from the google recaptcha in a struct
type RecaptchaResp struct {
	Success bool `json:"success"`
}

// ParseRecaptchaResp unmarshalls the byte recaptcha response into a golang struct
func ParseRecaptchaResp(data []byte) (RecaptchaResp, error) {
	var resp RecaptchaResp
	err := json.Unmarshal(data, &resp)
	return resp, err
}
