package common

import (
	"log"
	"os"
	"path"
)

func LogPrgName() {
	ex, err := ExePath()
	if err != nil {
		return
	}
	log.Printf("Executable : %s\n", ex)
	dir := path.Dir(ex)
	log.Printf("Installed in %s\n", dir)
}

func ExePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", FatalError(err)
	}
	return ex, nil
}
