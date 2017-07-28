package util

import (
	loggerUtil "git/inspursoft/board/src/collector/util/logger"
)
var Logger loggerUtil.NewLog
func init() {
	Logger.LogRegister(loggerUtil.LevelDebug)
}
