package main

import (
	"flag"
	"fmt"
)

var webServerPort int
var webServerHost string
var facebookClientId string
var facebookSecret string
var smtpUser string
var smtpPass string
var smtpHost string
var smtpPort int
var mongoUser string
var mongoPass string
var mongoHost string
var mongoPort int
var mongoDb string

func init() {
	flag.StringVar(&webServerHost, "webServerHost", "localhost", "Specify the web server host.")
	flag.IntVar(&webServerPort, "webServerPort", 3000, "Specify the web server port to listen to.")
	flag.StringVar(&facebookClientId, "facebookClientId", "...", "Facebook Client Id.")
	flag.StringVar(&facebookSecret, "facebookSecret", "...", "Facebook Secret.")
	flag.StringVar(&smtpUser, "smtpUser", "test@gmail.com", "Specify the smtp user.")
	flag.StringVar(&smtpPass, "smtpPass", "pass", "Specify the smtp user password.")
	flag.StringVar(&smtpHost, "smtpHost", "smtp.gmail.com", "Specify the smtp host.")
	flag.IntVar(&smtpPort, "smtpPort", 587, "Specify the smtp port.")
	flag.StringVar(&mongoUser, "mongoUser", "test", "Specify the mongodb user.")
	flag.StringVar(&mongoPass, "mongoPass", "pass", "Specify the the mongodb user password.")
	flag.StringVar(&mongoHost, "mongoHost", "localhost", "Specify the mongodb host.")
	flag.IntVar(&mongoPort, "mongoPort", 27017, "Specify the mongodb port.")
	flag.StringVar(&mongoDb, "mongoDb", "test", "Specify the mongodb database.")
	flag.Parse()
}

// get mongo connection string
func getMgoConnect() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", mongoUser, mongoPass, mongoHost, mongoPort, mongoDb)
}

// get mongo database name
func getMgoDb() string {
	return mongoDb
}

func getFacebookCallbackUri() string {
	return fmt.Sprintf("http://%s:%d/member/fb", webServerHost, webServerPort)
}

func getResetPasswordUri() string {
	return fmt.Sprintf("http://%s:%d/#/member/resetpassword", webServerHost, webServerPort)
}

func getForgetpasswordUri() string {
	return fmt.Sprintf("http://%s:%d/#/member/forgetpassword", webServerHost, webServerPort)
}

func getWebServerPort() string {
	return fmt.Sprintf(":%d", webServerPort)
}

func getWebServerHost() string {
	return webServerHost
}

func getFacebookClientId() string {
	return facebookClientId
}

func getFacebookSecret() string {
	return facebookSecret
}

func getFacebookAuthenticationUri() string {
	return "https://www.facebook.com/dialog/oauth"
}

func getFacebookAccessTokenUri() string {
	return "https://graph.facebook.com/v2.3/oauth/access_token"
}

func getFacebookUserInfoUri() string {
	return "https://graph.facebook.com/v2.3/me"
}

func getPasswordResetTokenExpire() int {
	return 3600
}

func getSmtpUser() string {
	return smtpUser
}

func getSmtpPass() string {
	return smtpPass
}

func getSmtpHost() string {
	return smtpHost
}

func getSmtpPort() string {
	return fmt.Sprintf(":%d", smtpPort)
}
