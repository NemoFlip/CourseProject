package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type Logger struct {
	Logger *zerolog.Logger
}

func UnitFormatter() {
	zerolog.TimestampFunc = func() time.Time {
		format := "2006-01-02 15:04:05"
		timeString := time.Now().Format(format)
		timeFormatted, _ := time.Parse(format, timeString)
		return timeFormatted
	}

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		// get function name
		funcPath := runtime.FuncForPC(pc).Name()
		funcName := filepath.Base(funcPath)

		// get directory name
		dirPath := filepath.Dir(file)
		dirName := filepath.Base(dirPath)

		return fmt.Sprintf("%s/%s:%s", dirName, funcName, strconv.Itoa(line))
	}
}

func InitLogger() *Logger {
	UnitFormatter()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	loggerFile, err := os.OpenFile("auth_service/cmd/log/logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic("error opening log file")
	}
	newLogger := zerolog.New(loggerFile).With().
		Timestamp().
		Str("service", "auth_service").
		Caller().
		Logger()

	return &Logger{Logger: &newLogger}
}
