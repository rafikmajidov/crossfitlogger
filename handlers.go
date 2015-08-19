package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func handleTemplatesHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}

func handleTemplatesLogin(w http.ResponseWriter, r *http.Request) {
	templateFile := "./templates/login.html"
	templateData, err := ioutil.ReadFile(templateFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", templateFile, err), http.StatusInternalServerError)
		return
	}

	//io.Copy(w, bytes.NewReader(templateData))
	randomString := srand(20)
	fbCodes[randomString] = randomString

	facebookOauthUri := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s&response_type=code&scope=email", getFacebookAuthenticationUri(), getFacebookClientId(), getFacebookCallbackUri(), randomString)

	output := strings.Replace(string(templateData), "<<facebookOauthUri>>", facebookOauthUri, -1)
	output = strings.Replace(output, "<<forgetPasswordUri>>", getForgetpasswordUri(), -1)
	fmt.Fprint(w, output)
}

func handleTemplatesLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}

func handleTemplatesForgetPassword(w http.ResponseWriter, r *http.Request) {
	templateFile := "./templates/forgetpassword.html"
	templateData, err := ioutil.ReadFile(templateFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", templateFile, err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(templateData))
}

func handleTemplatesResetPassword(w http.ResponseWriter, r *http.Request) {
	templateFile := "./templates/resetpassword.html"
	templateData, err := ioutil.ReadFile(templateFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", templateFile, err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(templateData))
}

func handleMemberAuthenticate(w http.ResponseWriter, r *http.Request) {
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Println(e)
	}
	// turn the request body (JSON) into a Visitor object
	var payload Visitor
	e = json.Unmarshal(data, &payload)
	if e != nil {
		fmt.Println(e)
	}
	memberSecure := MemberSecure{}
	member, status := payload.findMemberByEmail()
	// such member found in db, lets validate password
	if status == 1 {
		// hash of password from form matches hash from db
		if member.Password_hash == getHash(payload.Pass) {
			memberSecure.A = member.Auth_key
			memberSecure.S = 1
			loggedMembers[member.Auth_key] = member
			// password mismatch
		} else {
			memberSecure.A = "You have entered wrong password"
			memberSecure.S = 0
		}
		// such person not found in db
	} else {
		memberSecure.A = "Wrong login details"
		memberSecure.S = 0
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(memberSecure); err != nil {
		panic(err)
	}
}

func handleMemberRegister(w http.ResponseWriter, r *http.Request) {
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Println(e)
	}
	// turn the request body (JSON) into a Visitor object
	var payload Visitor
	e = json.Unmarshal(data, &payload)
	if e != nil {
		fmt.Println(e)
	}
	memberSecure := MemberSecure{}
	member, status := payload.findMemberByEmail()
	// such member not found, create new member
	if status == 0 {
		member.Pass = payload.Pass
		member.Email = payload.Email
		member.Password_hash = getHash(payload.Pass)
		member.Auth_key = srand(30)
		fmt.Println(member)
		member.Insert()
		memberSecure.A = member.Auth_key
		memberSecure.S = 1
		loggedMembers[member.Auth_key] = member
		// such member already exists
	} else {
		memberSecure.A = "Such member already exists"
		memberSecure.S = 0
	}
	fmt.Println(loggedMembers)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(memberSecure); err != nil {
		panic(err)
	}
}

func handleMemberFb(w http.ResponseWriter, r *http.Request) {
	memberSecure := MemberSecure{}
	values, err := url.ParseQuery(r.RequestURI)
	if err == nil {
		code := values.Get("/member/fb?code")
		state := values.Get("state")
		fbCode, ok := fbCodes[state]
		if ok == true && state == fbCode {

			url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s", getFacebookAccessTokenUri(), getFacebookClientId(), getFacebookCallbackUri(), getFacebookSecret(), code)

			fbResponse, err := http.Get(url)
			if err == nil {
				defer fbResponse.Body.Close()
				body, err := ioutil.ReadAll(fbResponse.Body)
				if err == nil {
					var payload FacebookAccessTokenResponse
					err = json.Unmarshal(body, &payload)
					if err == nil && payload.Access_token != "" {
						url2 := fmt.Sprintf("%s?access_token=%s", getFacebookUserInfoUri(), payload.Access_token)
						fbResponse2, err := http.Get(url2)
						if err == nil {
							defer fbResponse2.Body.Close()
							body2, err := ioutil.ReadAll(fbResponse2.Body)
							if err == nil {
								var payload2 FacebookUserInfoResponse
								err = json.Unmarshal(body2, &payload2)
								if err == nil && payload2.Id != "" {
									fmt.Println(payload2)
									member, status := payload2.findMemberByFacebookId()
									// such member not found, create new member
									if status == 0 {
										member.Facebook_id = payload2.Id
										member.Pass = srand(20)
										if payload2.Email != "" {
											member.Email = payload2.Email
										}
										if payload2.Phone != "" {
											member.Phone = payload2.Phone
										}
										member.Password_hash = getHash(member.Pass)
										member.Auth_key = srand(30)
										fmt.Println(member)
										member.Insert()
										memberSecure.A = member.Auth_key
										memberSecure.S = 1
										loggedMembers[member.Auth_key] = member
										// such member already exists
									} else {
										memberSecure.A = member.Auth_key
										memberSecure.S = 1
										loggedMembers[member.Auth_key] = member
									}
								}
							}
						}
					}
				}
			}

			//fmt.Fprint(w, "we will send data to fb again ", code)

		} else {
			//fmt.Fprint(w, "Some data from Fb do not match")
			memberSecure.A = "Some data from Fb do not match"
			memberSecure.S = 0
		}
	}

	fmt.Println(loggedMembers)
	res := ""
	if memberSecure.S == 1 {
		res = fmt.Sprintf("<script>window.localStorage.member='%s'; </script>", memberSecure.A)
	}

	fmt.Fprint(w, res, "<script>window.location.href='/';</script>")

}

func handleMemberLogout(w http.ResponseWriter, r *http.Request) {
	_, file, line, _ := runtime.Caller(1)
	// fmt.Println(loggedMembers)
	log.Printf(" debug %s:%d %v", file, line, loggedMembers)
	if isAuthorized(w, r) == false {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
	} else {
		delete(loggedMembers, r.Header.Get("Authorization"))
		memberSecure := MemberSecure{}
		memberSecure.A = ""
		memberSecure.S = 1
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(memberSecure); err != nil {
			panic(err)
		}
	}
}

func handleMemberForgetPassword(w http.ResponseWriter, r *http.Request) {
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Println(e)
	}
	// turn the request body (JSON) into a Visitor object
	var payload Visitor
	e = json.Unmarshal(data, &payload)
	if e != nil {
		fmt.Println(e)
	}
	memberSecure := MemberSecure{}
	member, status := payload.findMemberByEmail()
	// such member found in db, lets send email
	if status == 1 {
		// if reset_password_token is invalid, reset it
		// if reset_password_token is valid, do not do anything
		if member.Password_reset_token != "" {
			parts := strings.Split(member.Password_reset_token, "_")
			timestamp, _ := strconv.Atoi(parts[1])
			// token expiry is in the past, lets reset it
			if int64(timestamp+getPasswordResetTokenExpire()) < time.Now().Unix() {
				member.Password_reset_token = fmt.Sprintf("%s_%d", srand(20), time.Now().Unix())
				member.Update()
				// token is valid
			} else {
				fmt.Println("token is valid")
			}
			// empty token,lets set it
		} else {
			member.Password_reset_token = fmt.Sprintf("%s_%d", srand(20), time.Now().Unix())
			member.Update()
		}
		sendResetPasswordEmail(member)
		memberSecure.A = "Email has been sent"
		// such person not found in db
	} else {
		memberSecure.A = "No such member"

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(memberSecure); err != nil {
		panic(err)
	}
}

func handleMemberResetPassword(w http.ResponseWriter, r *http.Request) {
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Println(e)
	}
	// turn the request body (JSON) into a Visitor object
	var payload Visitor
	e = json.Unmarshal(data, &payload)
	if e != nil {
		fmt.Println(e)
	}
	memberSecure := MemberSecure{}
	fmt.Println("in handleMemberResetPassword")
	fmt.Println(payload)
	member, status := payload.findMemberByPasswordResetToken()
	// such member found in db, lets update
	// pass, password_hash, auth_key, password_reset_token
	if status == 1 {
		member.Pass = payload.Pass
		member.Password_hash = getHash(payload.Pass)
		member.Auth_key = srand(30)
		member.Password_reset_token = ""
		fmt.Println(member)
		member.Update()
		memberSecure.A = member.Auth_key
		memberSecure.S = 1
		loggedMembers[member.Auth_key] = member
	} else {
		memberSecure.A = "No such member"
		memberSecure.S = 0
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(memberSecure); err != nil {
		panic(err)
	}
}

func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Authorization") == "" {
		return false
	}

	_, ok := loggedMembers[r.Header.Get("Authorization")]
	return ok

}
