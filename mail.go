package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"strings"
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func sendResetPasswordEmail(member Member) {
	templateFile := "./templates/resetpasswordemail.html"
	templateData, err := ioutil.ReadFile(templateFile)

	resetPasswordUri := fmt.Sprintf("%s?t=%s", getResetPasswordUri(), member.Password_reset_token)
	body := strings.Replace(string(templateData), "<<memberEmail>>", member.Email, -1)
	body = strings.Replace(body, "<<resetPasswordUri>>", resetPasswordUri, -1)

	// Set up authentication information.
	// https://gist.github.com/andelf/5004821
	auth := smtp.PlainAuth("", getSmtpUser(), getSmtpPass(), getSmtpHost())

	from := mail.Address{"Support", getSmtpUser()}
	to := mail.Address{member.Email, member.Email}
	title := "Reset password"
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	err = smtp.SendMail(getSmtpHost()+getSmtpPort(), auth, from.Address, []string{to.Address}, []byte(message))
	if err != nil {
		fmt.Println(err)
	}
}
