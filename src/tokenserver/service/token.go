package service

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
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

func generateRandomString() string {
	length := 32
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	l := len(chars)
	result := make([]byte, length)
	_, err := rand.Read(result)
	if err != nil {
		logs.Error("Error reading random bytes: %v", err)
	}
	for i := 0; i < length; i++ {
		result[i] = chars[int(result[i])%l]
	}
	return string(result)
}

func InitService() error {
	var err error
	iniConf, err = config.NewConfig("ini", "app.conf")
	if err != nil {
		logs.Error("Failed to load config file: %+v", err)
		return err
	}
	expireSeconds, err = iniConf.Int("tokenExpireSeconds")
	if err != nil {
		logs.Error("Failed to get expireSeconds from config file: %+v", err)
		return err
	}
	secretKey = []byte(generateRandomString())
	logs.Info("Token expiration time now is %d second(s).", expireSeconds)
	return nil
}
