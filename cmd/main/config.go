package main

import (
	"auto-cert/pkg/config"
	"auto-cert/pkg/provider/traefik"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

var (
	Config config.Configurations
)

func init() {

	viper.BindPFlags(rootCmd.PersistentFlags())

	if home := homedir.HomeDir(); home != "" {
		viper.SetDefault("kubeconfig", filepath.Join(home, ".kube", "config"))
	}
	viper.AutomaticEnv()
	viper.BindEnv("kubeconfig", "KUBECONFIG")

	rootCmd.PersistentFlags().IntVarP(&Config.Count, "count", "c", 1, "the number of times to repeat the greeting")
	rootCmd.PersistentFlags().StringVarP(&Config.KubeConfig, "kubeconfig", "", viper.GetString("kubeconfig"), "Path to kubeconfig file, can be set by env KUBECONFIG")

	rootCmd.PersistentFlags().StringVarP(&Config.Log.Format, "log.format", "", "text", "log format")

	Config.Traefik = &traefik.Provider{
		Resync: 10 * time.Minute,
	}

	var formatter log.Formatter
	if Config.Log.Format == "json" {
		formatter = &log.JSONFormatter{}
	} else {
		formatter = &log.TextFormatter{
			FullTimestamp: true,
		}
	}
	log.SetFormatter(formatter)
}
