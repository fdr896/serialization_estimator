package cmd

import (
	"os"
	estimator "serialization_estimator/libs/services/estimator_service"
	"serialization_estimator/libs/support"

	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
    port string
    method string
)

var service = &cobra.Command{
	Use:   "./estimator --port <port> --method <method>",
	Long: "UPD service that estimates certain serialization methods efficiency.",
	Run: func(cmd *cobra.Command, args []string) {
        if len(port) == 0 {
            port = os.Getenv("PORT")
        }
        if len(method) == 0 {
            method = os.Getenv("METHOD")
        }

        service := estimator.New(port, method)
        if err := service.Start(); err != nil {
            zlog.Err(err).Msg("Failed to start estimator service")
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
    service.Flags().StringVarP(&port, "port", "p", "", "UPD port that estimator service listens")
    service.Flags().StringVarP(&method, "method", "m", "", "Serialization method that service estimates")
}


