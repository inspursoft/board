package service

import (
	"time"
)

type ServerTime struct {
	TimeNow int64 `json:"time_now"`
}

func GetServerTime() ServerTime {
	return ServerTime{
		TimeNow: getServerNow(),
	}
}
func getServerNow() int64 {
	return time.Now().Unix()
}
