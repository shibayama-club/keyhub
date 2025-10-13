package main

import (
	"context"
	"log"
	"strings"

	"github.com/shibayama-club/keyhub/cmd/serve"
	"github.com/shibayama-club/keyhub/internal/domain/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	debug := false
	configFile := ""

	cmd := &cobra.Command{
		Short: "keyhub",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.SetupLogger(debug)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cobra.OnInitialize(func() {
		if configFile != "" {
			viper.SetConfigFile(configFile)
			if err := viper.ReadInConfig(); err != nil {
				log.Fatal(err)
			}
		}
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
	})

	cmd.AddCommand(serve.ServeApp())
	cmd.AddCommand(serve.ServeConsole())
	flags := cmd.PersistentFlags()
	flags.BoolVar(&debug, "debug", false, "Enable debug mode")
	flags.StringVar(&configFile, "config", "", "Path to config file")

	if err := cmd.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
}
