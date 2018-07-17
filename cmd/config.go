package cmd

type ProviderConfig struct {
	Type    string            `json:"type"`
	Options map[string]string `json:"options"`
}

type PluginConfig struct {
	Force   bool              `json:"force"`
	Disable bool              `json:"disable"`
	Options map[string]string `json:"options"`
}

// Config describes the configuration files for evacuate
type Config struct {
	Provider ProviderConfig          `json:"provider"`
	Plugins  map[string]PluginConfig `json:"plugins"`
}
