package config

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	DBConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
		Debug    bool   `mapstructure:"debug"`
	}

	ConsoleConfig struct {
		OrganizationId  string `mapstructure:"organization_id"`
		OrganizationKey string `mapstructure:"organization_key"`
		JWTSecret       string `mapstructure:"jwt_secret"`
	}

	Config struct {
		Port     int      `mapstructure:"port"`
		Env      string   `mapstructure:"env"`
		Postgres DBConfig `mapstructure:"postgres"`
		Sentry   struct {
			DSN string `mapstructure:"dsn"`
		} `mapstructure:"sentry"`
		Console ConsoleConfig `mapstructure:"console"`
	}
)

func ConfigFlags(flags *pflag.FlagSet) {
	flags.String("env", "dev", "Environment (dev, prod)")
	flags.String("postgres.host", "localhost", "DB host")
	flags.Int("postgres.port", 5432, "DB port")
	flags.String("postgres.user", "", "DB user")
	flags.String("postgres.password", "", "DB password")
	flags.String("postgres.database", "", "DB name")
	flags.String("sentry.dsn", "", "Sentry DSN")
	flags.String("console.organization_id", "", "Organization ID(uuid)")
	flags.String("console.organization_key", "", "Organization Key")
	flags.String("console.jwt_secret", "", "JWT Secret for console authentication")
}

func ParseConfig[T any](cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return errors.Wrap(err, "failed to bind flags")
	}

	var config T
	if err := viper.Unmarshal(&config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	cmd.SetContext(context.WithValue(cmd.Context(), cmd, config))

	return nil
}
