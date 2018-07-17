package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CLIOptions struct {
	configFile  string
	verboseMode bool
}

var opts CLIOptions

var rootCmd = &cobra.Command{
	Use:   "evacuate",
	Short: "Have a misbehaving or compromised instance? Evacuate it.",
	Long:  "Have a misbehaving or compromised instance? Evacuate it.",
	Run: func(cmd *cobra.Command, args []string) {
		Evacuate(opts)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&opts.configFile, "config", "", "config file (default is ./evacuate.json)")
	rootCmd.PersistentFlags().BoolVarP(&opts.verboseMode, "verbose", "v", false, "Enable verbose/debug logging")
}

// GetConfig returns a configuration struct
func GetConfig() Config {
	var c Config

	defaultProviderConfig := ProviderConfig{}
	defaultPluginsConfig := make(map[string]PluginConfig)

	if opts.configFile == "" {
		opts.configFile = "evacuate.json"
	}

	viper.SetConfigFile(opts.configFile)

	viper.SetDefault("provider", defaultProviderConfig)
	viper.SetDefault("plugins", defaultPluginsConfig)

	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&c)

	if err != nil {
		panic(err)
	}

	return c
}
