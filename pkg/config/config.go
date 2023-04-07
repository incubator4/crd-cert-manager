package config

import (
	"auto-cert/pkg/provider/traefik"
	log "github.com/sirupsen/logrus"
)

type Configurations struct {
	Count      int
	KubeConfig string
	Log        Log
	Traefik    *traefik.Provider
}

type Log struct {
	Format string
	Level  log.Level
}
