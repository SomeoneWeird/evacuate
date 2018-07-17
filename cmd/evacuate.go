package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/SomeoneWeird/evacuate/plugins"
	"github.com/SomeoneWeird/evacuate/providers"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var loggers = map[string]*logrus.Entry{}

// Evacuate is the main function of this program
func Evacuate(cliOpts CLIOptions) {
	config := GetConfig()

	if cliOpts.verboseMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	availablePlugins := plugins.GetPlugins()

	for name := range availablePlugins {
		log := logrus.New()
		log.Formatter = new(prefixed.TextFormatter)

		if cliOpts.verboseMode {
			log.SetLevel(logrus.DebugLevel)
		}

		loggers[name] = log.WithFields(logrus.Fields{
			"prefix": fmt.Sprintf("plugin:%s", name),
		})
	}

	usablePlugins := getUsablePlugins(config, availablePlugins)

	directory := runPlugins(config, usablePlugins)

	file := bundleResults(config, directory)

	runProvider(config, cliOpts, file)
}

func getUsablePlugins(config Config, input plugins.List) plugins.List {
	output := make(plugins.List)

	for name, plugin := range input {
		// Check if plugins is explicitly enabled or disabled
		if pluginConfig, ok := config.Plugins[name]; ok {
			if pluginConfig.Force {
				output[name] = plugin
				loggers[name].Info("Explicitly enabled by config")
				continue
			} else if pluginConfig.Disable == true {
				loggers[name].Info("Explicitly disabled by config")
				continue
			}
		}

		c := make(chan bool)

		ctx := plugins.PluginContext{
			Config: config.Plugins[name].Options,
			Logger: loggers[name],
			Finish: c,
		}

		go plugin.ShouldRun(ctx)

		run := <-c

		loggers[name].WithFields(logrus.Fields{
			"willExecute": run,
		}).Info("")

		if run {
			output[name] = plugin
		}
	}

	return output
}

func runPlugins(config Config, list plugins.List) string {
	// TODO: err handling
	tmpDir, _ := ioutil.TempDir("", "evacuate")
	evacDir := fmt.Sprintf("%s/evacuate", tmpDir)
	_ = os.Mkdir(evacDir, 0700)

	for name, plugin := range list {
		loggers[name].WithFields(logrus.Fields{
			"running": true,
		}).Info()

		pluginOutputPath := fmt.Sprintf("%s/%s", evacDir, strings.ToLower(name))

		os.Mkdir(pluginOutputPath, 0777)

		c := make(chan bool)

		ctx := plugins.PluginContext{
			OutputPath: pluginOutputPath,
			Config:     config.Plugins[name].Options,
			Logger:     loggers[name],
			Finish:     c,
		}

		go plugin.Run(ctx)

		<-c

		loggers[name].WithFields(logrus.Fields{
			"finished": true,
		}).Info()
	}

	return tmpDir
}

func bundleResults(config Config, directory string) string {
	cmd := exec.Command("tar", "cvfz", "evacuate.tar.gz", "evacuate")
	cmd.Dir = directory

	_, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/evacuate.tar.gz", directory)
}

func runProvider(config Config, cliOpts CLIOptions, file string) {
	if config.Provider.Type == "" {
		logrus.Info("No provider specified. Evacuate finished. Results available: ", file)
		os.Exit(0)
	}

	log := logrus.New()
	log.Formatter = new(prefixed.TextFormatter)

	if cliOpts.verboseMode {
		log.SetLevel(logrus.DebugLevel)
	}

	availableProviders := providers.GetProviders()

	if _, ok := availableProviders[config.Provider.Type]; !ok {
		log.WithFields(logrus.Fields{
			"provider": config.Provider.Type,
		}).Error("Unknown provider")
		os.Exit(1)
	}

	logger := log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf("provider:%s", config.Provider.Type),
	})

	logger.WithFields(logrus.Fields{
		"running": true,
	}).Info()

	c := make(chan string)

	ctx := providers.ProviderContext{
		Config: config.Provider.Options,
		Logger: logger,
		Finish: c,
	}

	go availableProviders[config.Provider.Type].Run(ctx, file)

	output := <-c

	logger.WithFields(logrus.Fields{
		"finished": true,
	}).Info()

	logrus.Info("Evacuate finished: ", output)
}
