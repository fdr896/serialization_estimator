package cmd

import (
	"os"

	proxy "serialization_estimator/libs/services/proxy_service"
	"serialization_estimator/libs/support"

	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var port string

var service = &cobra.Command{
	Use:   "./proxy --port <port>",
	Long: "Service that proxies request to estimator services",
    Run: func(cmd *cobra.Command, args []string) {
        if len(port) == 0 {
            port = os.Getenv("PORT")
        }

        service, err := proxy.New(port)
        if err != nil {
            zlog.Err(err).Msg("Failed to create proxy service")
            return
        }
        if err := service.Start(); err != nil {
            zlog.Err(err).Msg("Failed to start proxy service")
        }
    },
}

func Execute() {
    support.InitLogger()

    if err := service.Execute(); err != nil {
        zlog.Err(err).Msg("Failed to run estimator service")
        os.Exit(1)
    }
}

func init() {
    service.Flags().StringVarP(&port, "port", "p", "", "UPD port that proxy service listens")
}


