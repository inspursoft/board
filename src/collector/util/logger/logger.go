package logger

import (
	"log"
	"net/http"
	"time"
)

func (a *NewLog) newLog() {
	switch a.loggerInit.LogLevel {
	case LevelFatal:
		*a = NewLog{logInit{LevelDebug}, noneOut,
			noneOut, noneOut, noneOut,printOut}
	case LevelError:
		*a = NewLog{logInit{LevelInfo}, noneOut,
			noneOut, noneOut, printOut, printOut}
	case LevelWarn:
		*a = NewLog{logInit{LevelWarn}, noneOut,
			noneOut, printOut, printOut, printOut}
	case LevelInfo:
		*a = NewLog{logInit{LevelError}, noneOut,
			printOut, printOut, printOut, printOut}
	case LevelDebug:
		*a = NewLog{logInit{LevelFatal}, printOut,
			printOut, printOut, printOut, printOut}
	}
}

func (a *NewLog) LogRegister(c LogLevel) {

	a.loggerInit.LogLevel = c
	a.newLog()

}

func noneOut(none ...interface{}) {
}

func printOut(parameter ...interface{}) {
	log.Println(parameter...)
}

func (NewLog) HttpLog(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
