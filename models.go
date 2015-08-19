package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Visitor model
type Visitor struct {
	Pass                 string `json:"pass"`
	Email                string `json:"email"`
	Password_reset_token string `json:"t"`
}

// Facebook Response model with access_token
type FacebookAccessTokenResponse struct {
	Access_token string `json:"access_token"`
	Token_type   string `json:"token_type"`
	Expires_in   int    `json:"expires_in"`
}

// Facebook Response model with user info
type FacebookUserInfoResponse struct {
	Id    string `json:"id"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// Member model
type Member struct {
	Id                   bson.ObjectId `bson:"_id,omitempty"`
	Pass                 string
	Email                string
	Facebook_id          string
	Auth_key             string
	Password_hash        string
	Password_reset_token string
	Phone                string
}

// MemberSecure model
type MemberSecure struct {
	A string
	S int
}

func (v Visitor) findMemberByEmail() (Member, int) {
	status := 0

	session, err := mgo.Dial(getMgoConnect())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(getMgoDb()).C("member")

	member := Member{}

	err = c.Find(bson.M{"email": v.Email}).One(&member)
	if err != nil {
		fmt.Println("we must be heere")
		fmt.Println(err)
		status = 0
	} else {
		status = 1
	}

	return member, status
}

func (v Visitor) findMemberByPasswordResetToken() (Member, int) {
	status := 0

	session, err := mgo.Dial(getMgoConnect())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(getMgoDb()).C("member")

	member := Member{}

	err = c.Find(bson.M{"password_reset_token": v.Password_reset_token}).One(&member)
	if err != nil {
		fmt.Println("we must be in find by password reset token error")
		fmt.Println(err)
		status = 0
	} else {
		status = 1
	}

	return member, status
}

func (m Member) Insert() {
	session, err := mgo.Dial(getMgoConnect())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(getMgoDb()).C("member")
	err = c.Insert(&m)
	if err != nil {
		fmt.Println(err)
	}
}

func (m Member) Update() {
	session, err := mgo.Dial(getMgoConnect())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(getMgoDb()).C("member")
	err = c.Update(bson.M{"_id": m.Id}, m)
	if err != nil {
		fmt.Println(err)
	}
}

func (f FacebookUserInfoResponse) findMemberByFacebookId() (Member, int) {
	status := 0

	session, err := mgo.Dial(getMgoConnect())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(getMgoDb()).C("member")

	member := Member{}

	err = c.Find(bson.M{"facebook_id": f.Id}).One(&member)
	if err != nil {
		fmt.Println("we must be find by facebook id")
		fmt.Println(err)
		status = 0
	} else {
		status = 1
	}

	return member, status
}
