package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func New() (*log.Logger, error) {
	date := time.Now().Format("2006-01-02")

	logDir := filepath.Join(
		"storage",
		"logs",
	)

	err := os.MkdirAll(
		logDir,
		0755,
	)

	if err != nil {
		return nil, err
	}

	path := filepath.Join(
		logDir,
		fmt.Sprintf("%s.log", date),
	)

	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return nil, err
	}

	return log.New(
		file,
		"",
		log.LstdFlags,
	), nil
}
