package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
	io.Closer
}

func New(logFilePath string) (Logger, error) {
	info, err := os.Stat(logFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return Logger{}, err
		}
	} else if info.IsDir() {
		return Logger{}, getErrLogPathIsDir()
	}

	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return Logger{}, err
	}

	if _, err := f.WriteString("\n\n"); err != nil {
		return Logger{}, err
	}

	logger := zerolog.New(f)

	return Logger{
		Logger: &logger,
		Closer: f,
	}, nil
}
