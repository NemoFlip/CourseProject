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
	InfoLogger  *zerolog.Logger
	ErrorLogger *zerolog.Logger
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

	loggerInfoFile, err := os.OpenFile("auth_service/cmd/log/logs.info", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		panic("unable to open info log file")
	}
	loggerErrorFile, err := os.OpenFile("auth_service/cmd/log/logs.error", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		panic("unable to open error log file")
	}
	infoLogger := zerolog.New(loggerInfoFile).With().
		Timestamp().
		Str("service", "auth_service").
		Caller().
		Logger()
	errorLogger := zerolog.New(loggerErrorFile).With().
		Timestamp().
		Str("service", "auth_service").
		Caller().
		Logger()

	return &Logger{InfoLogger: &infoLogger, ErrorLogger: &errorLogger}
}
