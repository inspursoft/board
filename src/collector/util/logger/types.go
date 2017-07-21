package logger

type NewLog struct {
	loggerInit logInit
	SetDebug   func(parameter ...interface{})
	SetInfo    func(parameter ...interface{})
	SetWarn    func(parameter ...interface{})
	SetError   func(parameter ...interface{})
	SetFatal   func(parameter ...interface{})
}
type logInit struct {
	LogLevel
}
type LogLevel int

const (
	_          LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)
