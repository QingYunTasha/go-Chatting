package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"
)

type SecretKey struct {
	key        string
	expireTime time.Time
}

func NewSecretKey() *SecretKey {
	sk := SecretKey{}
	sk.key = sk.generateSecretKey()
	sk.expireTime = time.Now().Add(time.Hour * 24)
	return &sk
}

func (k *SecretKey) Get() string {
	if k.expireTime.Before(time.Now()) {
		k.key = k.generateSecretKey()
	}

	return k.key
}

func (k *SecretKey) generateSecretKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}
