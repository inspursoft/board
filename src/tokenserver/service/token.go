package service

import (
	"fmt"

	"time"

	"encoding/base64"

	"log"

	"github.com/astaxie/beego/config"
	jwt "github.com/dgrijalva/jwt-go"
)

var iniConf config.Configer

var secretKey []byte
var expireSeconds int

type tokenClaims struct {
	Payload map[string]interface{}
	jwt.StandardClaims
}

func Sign(payload map[string]interface{}) (string, error) {
	claims := tokenClaims{
		payload,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expireSeconds)).Unix(),
		},
	}
	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(secretKey)
}

func Verify(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims.Payload, nil
	}
	return nil, nil
}

func init() {
	var err error
	iniConf, err = config.NewConfig("ini", "app.conf")
	if err != nil {
		log.Fatalf("Failed to load config file: %+v\n", err)
	}
	expireSeconds, err = iniConf.Int("tokenExpireSeconds")
	if err != nil {
		log.Fatalf("Failed to get expireSeconds from config file: %+v\n", err)
	}
	encodedKey := iniConf.String("tokenSecretKey")
	secretKey, err = base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		log.Fatalf("Failed to decode secret key from config file: %+v\n", err)
	}
	fmt.Printf("Token expiration time now is %d second(s).\n", expireSeconds)
}
