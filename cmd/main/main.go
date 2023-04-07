package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Errorf("error run command")
		os.Exit(1)
	}
}
