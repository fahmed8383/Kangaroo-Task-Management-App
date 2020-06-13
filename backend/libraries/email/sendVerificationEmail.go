package email

import (
	"github.com/fahmed8383/SchedulingApp/libraries/api"
	"github.com/fahmed8383/SchedulingApp/libraries/setup"

	"gopkg.in/gomail.v2"
)

//SendVerificationEmail is responsible for sending the user a verification code to confirm their email
func SendVerificationEmail(secrets *setup.Secrets, info api.RegInfo, token string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "kangaroo.devplayground@gmail.com")
	mail.SetHeader("To", info.Email)
	mail.SetHeader("Subject", "Kangaroo: Email Verification")
	mail.SetBody("text/html", "Your email verification code for Kangaroo is: <br><br>"+token)
	d := gomail.NewDialer("smtp.gmail.com", 587, "kangaroo.devplayground@gmail.com", secrets.GmailPass)
	err := d.DialAndSend(mail)
	return err
}
