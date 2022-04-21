package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Option func(l *viperLoader) error

func WithEnv(prefix ...string) Option {
	return func(l *viperLoader) error {
		l.useEnv = true
		if len(prefix) > 0 {
			l.envPrefix = strings.TrimSpace(prefix[0])
		}
		return nil
	}
}

func WithFile(file string) Option {
	return func(l *viperLoader) error {
		l.confFile = file
		return nil
	}
}

func WithCobra(cmd *cobra.Command) Option {
	return func(l *viperLoader) error {
		configFile, _ := cmd.Flags().GetString("config")
		if configFile != "" {
			l.confFile = configFile
		}
		// TODO: bind flags.
		return nil
	}
}

func bindFlags(envPrefix string, cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			_ = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
