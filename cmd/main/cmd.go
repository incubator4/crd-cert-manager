package main

import (
	"auto-cert/pkg/provider/aggregator"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var rootCmd = &cobra.Command{
	Use:   "auto-cert",
	Short: "A sample CLI application",
	Run: func(cmd *cobra.Command, args []string) {
		var config *rest.Config
		var err error
		if config, err = rest.InClusterConfig(); err != nil {
			// Running outside a Kubernetes cluster.

			// use the current context in kubeconfig
			config, err = clientcmd.BuildConfigFromFlags("", Config.KubeConfig)
			if err != nil {
				fmt.Println("KUBECONFIG ", Config.KubeConfig)
				panic(err.Error())
			}
		}

		ctxPool := context.Background()

		provider := aggregator.NewAggregatorProvider(Config)
		eventCh, err := provider.WatchAll(config, ctxPool.Done())
		if err != nil {
			log.WithError(err).Error("Error watch channel")
		}

		for {
			select {
			case <-ctxPool.Done():
				break
			case event := <-eventCh:
				provider.Provide(ctxPool, event)
			}
		}

	},
}
