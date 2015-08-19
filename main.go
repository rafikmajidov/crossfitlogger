package main

import (
	"log"
	"net/http"
)

// global array to keep all logged members
var loggedMembers = make(map[string]Member)

// global array to keep all facebook members
var fbCodes = make(map[string]string)

func main() {
	http.HandleFunc("/templates/home", handleTemplatesHome)
	http.HandleFunc("/templates/login", handleTemplatesLogin)
	http.HandleFunc("/templates/logout", handleTemplatesLogout)
	http.HandleFunc("/member/authenticate", handleMemberAuthenticate)
	http.HandleFunc("/member/register", handleMemberRegister)
	http.HandleFunc("/member/fb", handleMemberFb)
	http.HandleFunc("/member/logout", handleMemberLogout)
	http.HandleFunc("/member/forgetpassword", handleMemberForgetPassword)
	http.HandleFunc("/templates/forgetpassword", handleTemplatesForgetPassword)
	http.HandleFunc("/templates/resetpassword", handleTemplatesResetPassword)
	http.HandleFunc("/member/resetpassword", handleMemberResetPassword)
	http.Handle("/", http.FileServer(http.Dir("./public")))

	log.Printf("Server started: http://localhost%s", getWebServerPort())
	// log.Println(srand(10))
	log.Fatal(http.ListenAndServe(getWebServerPort(), nil))
}
