package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/luisnquin/nix-search/internal"
	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
	io.Closer
}

func New() (Logger, error) {
	logFilePath := filepath.Join(
		os.TempDir(), fmt.Sprintf("%s.log", internal.PROGRAM_NAME),
	)

	info, err := os.Stat(logFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return Logger{}, err
		}
	} else if info.IsDir() {
		return Logger{}, getErrLogPathIsDir()
	}

	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return Logger{}, err
	}

	if _, err := f.WriteString("\n\n"); err != nil {
		return Logger{}, err
	}

	zerolog.LevelFieldName = "lvl"
	logger := zerolog.New(f).With().Int64("i", time.Now().Unix()).Logger()

	return Logger{
		Logger: &logger,
		Closer: f,
	}, nil
}
