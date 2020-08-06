package email

import (
	"github.com/fahmed8383/SchedulingApp/libraries/api"
	"github.com/fahmed8383/SchedulingApp/libraries/setup"

	"gopkg.in/gomail.v2"
)

//SendPasswordResetEmail is responsible for sending the user a link to reset their password
func SendPasswordResetEmail(secrets *setup.Secrets, info api.RegInfo, link string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "kangaroo.devplayground@gmail.com")
	mail.SetHeader("To", info.Email)
	mail.SetHeader("Subject", "Kangaroo: Password Reset")
	mail.SetBody("text/html", "A request has been made to reset your password for <b>"+info.Username+"</b>. <br><br>Click the following link to complete the request <br><br><a href='"+link+"'>Reset My Password</a><br><br>If you did not request this change, please ignore this request.")
	d := gomail.NewDialer("smtp.gmail.com", 587, "kangaroo.devplayground@gmail.com", secrets.GmailPass)
	err := d.DialAndSend(mail)
	return err
}
