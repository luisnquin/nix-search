package log

import "errors"

func getErrLogPathIsDir() error {
	return errors.New("log path is a directory")
}
