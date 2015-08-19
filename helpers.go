package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
)

const dictionary = "abc~!@$^defghijkmnpq(}*{>)<rstuvwxyzABCDEFGHJKL~!@$^*><MNPQRSTUVWXYZ123456789"

// generates a random string of fixed size
func srand(size int) string {
	bytes := make([]byte, size)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func getHash(word string) string {
	h := sha1.New()
	h.Write([]byte(word))
	sha1Hash := hex.EncodeToString(h.Sum(nil))

	return sha1Hash
}
