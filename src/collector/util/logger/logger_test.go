package logger

import "testing"

func TestGenLog_OutL1(t *testing.T)  {
	var s NewLog
	s.LogRegister(LevelInfo)
	s.SetInfo("asd","sdfsdf")
}